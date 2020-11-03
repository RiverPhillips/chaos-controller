// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2020 Datadog, Inc.

package injector

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/DataDog/chaos-controller/network"
	chaostypes "github.com/DataDog/chaos-controller/types"
	"go.uber.org/zap"
)

const (
	flowEgress  = "egress"
	flowIngress = "ingress"
)

// linkOperation represents a tc operation on a single network interface combined with the parent to bind to and the handle identifier to use
type linkOperation func(network.NetlinkLink, string, uint32) error

// NetworkDisruptionConfig provides an interface for using the network traffic controller for new disruptions
type NetworkDisruptionConfig interface {
	AddNetem(delay time.Duration, drop int, corrupt int)
	AddOutputLimit(bytesPerSec uint)
	ApplyOperations() error
	ClearOperations() error
}

// NetworkDisruptionConfigStruct contains all needed drivers to create a network disruption using `tc`
type NetworkDisruptionConfigStruct struct {
	Log               *zap.SugaredLogger
	TrafficController network.TrafficController
	NetlinkAdapter    network.NetlinkAdapter
	DNSClient         network.DNSClient
	hosts             []string
	port              int
	protocol          string
	flow              string
	operations        []linkOperation
}

// NewNetworkDisruptionConfig creates a new network disruption object using the given netlink, dns, etc.
func NewNetworkDisruptionConfig(logger *zap.SugaredLogger, tc network.TrafficController, netlink network.NetlinkAdapter, dns network.DNSClient, hosts []string, port int, protocol string, flow string) NetworkDisruptionConfig {
	return &NetworkDisruptionConfigStruct{
		Log:               logger,
		TrafficController: tc,
		NetlinkAdapter:    netlink,
		DNSClient:         dns,
		hosts:             hosts,
		port:              port,
		protocol:          protocol,
		flow:              flow,
		operations:        []linkOperation{},
	}
}

// NewNetworkDisruptionConfigWithDefaults creates a new network disruption object using default netlink, dns, etc.
func NewNetworkDisruptionConfigWithDefaults(logger *zap.SugaredLogger, hosts []string, port int, protocol string, flow string) NetworkDisruptionConfig {
	return NewNetworkDisruptionConfig(logger, network.NewTrafficController(logger), network.NewNetlinkAdapter(), network.NewDNSClient(), hosts, port, protocol, flow)
}

// getInterfacesByIP returns the interfaces used to reach the given hosts
// if hosts is empty, all interfaces are returned
func (c *NetworkDisruptionConfigStruct) getInterfacesByIP(hosts []string) (map[string][]*net.IPNet, error) {
	linkByIP := map[string][]*net.IPNet{}

	if len(hosts) > 0 {
		c.Log.Info("auto-detecting interfaces to apply disruption to...")

		// resolve hosts
		ips, err := resolveHosts(c.DNSClient, hosts)
		if err != nil {
			return nil, fmt.Errorf("can't resolve given hosts: %w", err)
		}

		// get the association between IP and interfaces to know
		// which interfaces we have to inject disruption to
		for _, ip := range ips {
			// get routes for resolved destination IP
			routes, err := c.NetlinkAdapter.RoutesForIP(ip)
			if err != nil {
				return nil, fmt.Errorf("can't get route for IP %s: %w", ip.String(), err)
			}

			// for each route, get the related interface and add it to the association
			// between interfaces and IPs
			for _, route := range routes {
				c.Log.Infof("IP %s belongs to interface %s", ip.String(), route.Link().Name())

				// store association, initialize the map entry if not present yet
				if _, ok := linkByIP[route.Link().Name()]; !ok {
					linkByIP[route.Link().Name()] = []*net.IPNet{}
				}

				linkByIP[route.Link().Name()] = append(linkByIP[route.Link().Name()], ip)
			}
		}
	} else {
		c.Log.Info("no hosts specified, all interfaces will be impacted")

		// prepare links/IP association by pre-creating links
		// exclude lo interface
		links, err := c.NetlinkAdapter.LinkList()
		if err != nil {
			c.Log.Fatalf("can't list links: %w", err)
		}

		for _, link := range links {
			if link.Name() != "lo" {
				c.Log.Infof("adding interface %s", link.Name())
				linkByIP[link.Name()] = []*net.IPNet{}
			}
		}
	}

	return linkByIP, nil
}

// ApplyOperations applies the added operations by building a tc tree
// Here's what happen on tc side:
//  - a first prio qdisc will be created and attached to root
//    - it'll be used to apply the fisrt filter, filtering on packet IP destination, source/destination ports and protocol
//  - a second prio qdisc will be created and attached to the first one
//    - it'll be used to apply the second filter, filtering on packet classid to identify packets coming from the targeted process
//  - operations will be chained to the second band of the second prio qdisc
//  - a cgroup filter will be created to classify packets according to their classid (if any)
//  - a filter will be created to redirect traffic related to the specified host(s) through the last prio band
//    - if no host, port or protocol is specified, a filter redirecting all the traffic (0.0.0.0/0) to the disrupted band will be created
//  - a last filter will be created to redirect traffic related to the local node through a not disrupted band
//
// Here's the tc tree representation:
// root (1:) <-- prio qdisc with 4 bands with a filter classifying packets matching the given dst ip, src/dst ports and protocol with class 1:4
//   |- (1:1) <-- first band
//   |- (1:2) <-- second band
//   |- (1:3) <-- third band
//   |- (1:4) <-- fourth band
//     |- (2:) <-- prio qdisc with 2 bands with a cgroup filter to classify packets according to their classid (packets with classid 2:2 will be affected by operations)
//       |- (2:1) <-- first band
//       |- (2:2) <-- second band
//         |- (3:) <-- first operation
//           |- (4:) <-- second operation
//             ...
func (c *NetworkDisruptionConfigStruct) ApplyOperations() error {
	// get the default route to exclude it later
	defaultRoute, err := c.NetlinkAdapter.DefaultRoute()
	if err != nil {
		return fmt.Errorf("error getting the default route: %w", err)
	}

	c.Log.Infof("detected default gateway IP %s on interface %s", defaultRoute.Gateway().String(), defaultRoute.Link().Name())

	// get the targeted pod node IP from the environment variable
	hostIP, ok := os.LookupEnv(chaostypes.TargetPodHostIPEnv)
	if !ok {
		return fmt.Errorf("%s environment variable must be set with the target pod node IP", chaostypes.TargetPodHostIPEnv)
	}

	c.Log.Infof("target pod node IP is %s", hostIP)

	hostIPNet := &net.IPNet{
		IP:   net.ParseIP(hostIP),
		Mask: net.CIDRMask(32, 32),
	}

	// get routes going to this node IP to add a filter excluding this IP from the disruptions later
	// it is used to allow the node to reach the pod even with disruptions applied
	hostIPRoutes, err := c.NetlinkAdapter.RoutesForIP(hostIPNet)
	if err != nil {
		return fmt.Errorf("error getting target pod node IP routes: %w", err)
	}

	// get the interfaces per IP map
	linkByIP, err := c.getInterfacesByIP(c.hosts)
	if err != nil {
		return fmt.Errorf("can't get interfaces per IP listing: %w", err)
	}

	// for each link/ip association, add disruption
	for linkName, ips := range linkByIP {
		clearTxQlen := false

		// retrieve link from name
		link, err := c.NetlinkAdapter.LinkByName(linkName)
		if err != nil {
			return fmt.Errorf("can't retrieve link %s: %w", linkName, err)
		}

		// set the tx qlen if not already set as it is required to create a prio qdisc without dropping
		// all the outgoing traffic
		// this qlen will be removed once the injection is done if it was not present before
		if link.TxQLen() == 0 {
			c.Log.Infof("setting tx qlen for interface %s", link.Name())

			// set clear flag to true so we can clean up this once qdiscs are created
			clearTxQlen = true

			// set qlen
			if err := link.SetTxQLen(1000); err != nil {
				return fmt.Errorf("can't set tx queue length on interface %s: %w", link.Name(), err)
			}
		}

		// create a new qdisc for the given interface of type prio with 4 bands instead of 3
		// we keep the default priomap, the extra band will be used to filter traffic going to the specified IP
		// we only create this qdisc if we want to target traffic going to some hosts only, it avoids to apply disruptions to all the traffic for a bit of time
		priomap := [16]uint32{1, 2, 2, 2, 1, 2, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1}

		if err := c.TrafficController.AddPrio(link.Name(), "root", 1, 4, priomap); err != nil {
			return fmt.Errorf("can't create a new qdisc for interface %s: %w", link.Name(), err)
		}

		// create second prio with only 2 bands to filter traffic with a specific classid
		if err := c.TrafficController.AddPrio(link.Name(), "1:4", 2, 2, [16]uint32{}); err != nil {
			return fmt.Errorf("can't create a new qdisc for interface %s: %w", link.Name(), err)
		}

		// create cgroup filter
		if err := c.TrafficController.AddCgroupFilter(link.Name(), "2:0", 2); err != nil {
			return fmt.Errorf("can't create the cgroup filter for interface %s: %w", link.Name(), err)
		}

		// parent 2:2 refers to the 2nd band of the 2nd prio qdisc
		// handle starts from 3 because 1 and 2 are used by the 2 prio qdiscs
		parent := "2:2"
		handle := uint32(3)

		// add operations
		for _, operation := range c.operations {
			if err := operation(link, parent, handle); err != nil {
				return fmt.Errorf("could not perform operation on newly created qdisc for interface %s: %w", link.Name(), err)
			}

			// update parent reference and handle identifier for the next operation
			// the next operation parent will be the current handle identifier
			// the next handle identifier is just an increment of the actual one
			parent = fmt.Sprintf("%d:", handle)
			handle++
		}

		// handle flow direction
		// if flow is egress, filter will be on destination port
		// if flow is ingress, filter will be on source port
		var srcPort, dstPort int

		switch c.flow {
		case flowEgress:
			dstPort = c.port
		case flowIngress:
			srcPort = c.port
		default:
			return fmt.Errorf("unsupported flow: %s", c.flow)
		}

		// if some hosts are targeted, create one filter per host to redirect the traffic to the disrupted band
		// otherwise, create a filter redirecting all the traffic (0.0.0.0/0) using the given port and protocol to the disrupted band
		if len(ips) > 0 {
			for _, ip := range ips {
				if err := c.TrafficController.AddFilter(link.Name(), "1:0", 0, nil, ip, srcPort, dstPort, c.protocol, "1:4"); err != nil {
					return fmt.Errorf("can't add a filter to interface %s: %w", link.Name(), err)
				}
			}
		} else {
			_, nullIP, _ := net.ParseCIDR("0.0.0.0/0")
			if err := c.TrafficController.AddFilter(link.Name(), "1:0", 0, nil, nullIP, srcPort, dstPort, c.protocol, "1:4"); err != nil {
				return fmt.Errorf("can't add a filter to interface %s: %w", link.Name(), err)
			}
		}

		// the following lines are used to allow the node and the pod to communicate even with disruptions applied
		// depending on the network configuration, only one of those filters can be useful but we must add all of them
		// NOTE: those filters must be added after every other filters applied to the interface so they are used first

		// this filter allows the pod to communicate with the default route gateway IP
		if defaultRoute.Link().Name() == link.Name() {
			gatewayIP := &net.IPNet{
				IP:   defaultRoute.Gateway(),
				Mask: net.CIDRMask(32, 32),
			}

			if err := c.TrafficController.AddFilter(link.Name(), "1:0", 0, nil, gatewayIP, 0, 0, "", "1:1"); err != nil {
				return fmt.Errorf("can't add the default route gateway IP filter: %w", err)
			}
		}

		// this filter allows the pod to communicate with the node IP
		for _, hostIPRoute := range hostIPRoutes {
			if hostIPRoute.Link().Name() == link.Name() {
				if err := c.TrafficController.AddFilter(link.Name(), "1:0", 0, nil, hostIPNet, 0, 0, "", "1:1"); err != nil {
					return fmt.Errorf("can't add the target pod node IP filter: %w", err)
				}
			}
		}

		// reset the interface transmission queue length once filters have been created (only if qlen has been set earlier)
		if clearTxQlen {
			c.Log.Infof("clearing tx qlen for interface %s", link.Name())

			if err := link.SetTxQLen(0); err != nil {
				return fmt.Errorf("can't clear %s link transmission queue length: %w", link.Name(), err)
			}
		}
	}

	return nil
}

// AddNetem adds network disruptions using the drivers in the NetworkDisruptionConfigStruct
func (c *NetworkDisruptionConfigStruct) AddNetem(delay time.Duration, drop int, corrupt int) {
	// closure which adds netem disruptions
	operation := func(link network.NetlinkLink, parent string, handle uint32) error {
		return c.TrafficController.AddNetem(link.Name(), parent, handle, delay, drop, corrupt)
	}

	c.operations = append(c.operations, operation)
}

// AddOutputLimit adds a network bandwidth disruption using the drivers in the NetworkDisruptionConfigStruct
func (c *NetworkDisruptionConfigStruct) AddOutputLimit(bytesPerSec uint) {
	// closure which adds a bandwidth limit
	operation := func(link network.NetlinkLink, parent string, handle uint32) error {
		return c.TrafficController.AddOutputLimit(link.Name(), parent, handle, bytesPerSec)
	}

	c.operations = append(c.operations, operation)
}

// ClearOperations removes all disruptions by clearing all custom qdiscs created for the given config struct (filters will be deleted as well)
func (c *NetworkDisruptionConfigStruct) ClearOperations() error {
	linkByIP, err := c.getInterfacesByIP(c.hosts)
	if err != nil {
		return fmt.Errorf("can't get interfaces per IP map: %w", err)
	}

	for linkName := range linkByIP {
		c.Log.Infof("clearing root qdisc for interface %s", linkName)

		// retrieve link from name
		link, err := c.NetlinkAdapter.LinkByName(linkName)
		if err != nil {
			return fmt.Errorf("can't retrieve link %s: %w", linkName, err)
		}

		// ensure qdisc isn't cleared before clearing it to avoid any tc error
		cleared, err := c.TrafficController.IsQdiscCleared(link.Name())
		if err != nil {
			return fmt.Errorf("can't ensure the %s link qdisc is cleared or not: %w", link.Name(), err)
		}

		// clear link qdisc if needed
		if !cleared {
			if err := c.TrafficController.ClearQdisc(link.Name()); err != nil {
				return fmt.Errorf("can't delete the %s link qdisc: %w", link.Name(), err)
			}
		} else {
			c.Log.Infof("%s link qdisc is already cleared, skipping", link.Name())
		}
	}

	return nil
}