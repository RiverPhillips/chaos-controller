// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2020 Datadog, Inc.

package metrics

import (
	"fmt"
	"time"

	"github.com/DataDog/chaos-controller/metrics/datadog"
	"github.com/DataDog/chaos-controller/metrics/noop"
	"github.com/DataDog/chaos-controller/metrics/types"
	chaostypes "github.com/DataDog/chaos-controller/types"
)

// Sink describes a metric sink
type Sink interface {
	Close() error
	EventWithTags(title, text string, tags []string) error
	Flush() error
	GetSinkName() string
	MetricCleaned(containerID, uid string, succeed bool, kind chaostypes.DisruptionKind, tags []string) error
	MetricCleanupDuration(duration time.Duration, tags []string) error
	MetricInjectDuration(duration time.Duration, tags []string) error
	MetricInjected(containerID, uid string, succeed bool, kind chaostypes.DisruptionKind, tags []string) error
	MetricIPTablesRulesInjected(containerID, uid string, kind chaostypes.DisruptionKind, tags []string) error
	MetricPodsCreated(targetPod, instanceName, namespace, phase string, succeed bool) error
	MetricReconcile() error
	MetricReconcileDuration(duration time.Duration, tags []string) error
}

// GetSink returns an initiated sink
func GetSink(driver types.SinkDriver, app types.SinkApp) (Sink, error) {
	switch driver {
	case types.SinkDriverDatadog:
		return datadog.New(app)
	case types.SinkDriverNoop:
		return noop.New(), nil
	default:
		return nil, fmt.Errorf("unsupported metrics sink: %s", driver)
	}
}