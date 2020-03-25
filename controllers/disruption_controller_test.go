// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2020 Datadog, Inc.

package controllers

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"sigs.k8s.io/controller-runtime/pkg/client"

	chaosv1beta1 "github.com/DataDog/chaos-controller/api/v1beta1"
	chaostypes "github.com/DataDog/chaos-controller/types"
	"golang.org/x/net/context"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
)

// listChaosPods returns all the chaos pods for the given instance and mode
func listChaosPods(instance *chaosv1beta1.Disruption, mode chaostypes.PodMode) (corev1.PodList, error) {
	l := corev1.PodList{}
	ls := labels.NewSelector()
	instancePods := corev1.PodList{}

	// create requirements
	targetPodRequirement, _ := labels.NewRequirement(chaostypes.TargetPodLabel, selection.In, []string{"foo", "bar"})
	podModeRequirement, _ := labels.NewRequirement(chaostypes.PodModeLabel, selection.Equals, []string{string(mode)})

	// add requirements to label selector
	ls = ls.Add(*targetPodRequirement, *podModeRequirement)

	// get matching pods
	if err := k8sClient.List(context.Background(), &l, &client.ListOptions{
		Namespace:     "default",
		LabelSelector: ls,
	}); err != nil {
		return corev1.PodList{}, fmt.Errorf("can't list chaos pods: %w", err)
	}

	// filter to get only pods owned by the given instance
	for _, pod := range l.Items {
		if metav1.IsControlledBy(&pod, instance) {
			instancePods.Items = append(instancePods.Items, pod)
		}
	}

	return instancePods, nil
}

// expectChaosPod retrieves the list of created chaos pods related to the given and to the
// given mode (inject or clean) and returns an error if it doesn't
// equal the given count
func expectChaosPod(instance *chaosv1beta1.Disruption, mode chaostypes.PodMode, count int) error {
	l, err := listChaosPods(instance, mode)
	if err != nil {
		return err
	}

	// ensure count is correct
	if len(l.Items) != count {
		return fmt.Errorf("unexpected injection pods count: %d", len(l.Items))
	}

	// ensure generated pods have the needed fields
	for _, p := range l.Items {
		if p.GenerateName == "" {
			return fmt.Errorf("GenerateName field can't be empty")
		}
		if p.Namespace != instance.Namespace {
			return fmt.Errorf("pod namesapce must match instance namespace")
		}
		if p.Labels[chaostypes.PodModeLabel] != string(mode) {
			return fmt.Errorf("pod mode label should be set and match the actual mode")
		}
		if len(p.Spec.Containers[0].Args) == 0 {
			return fmt.Errorf("pod container args must be set")
		}
		if p.Spec.Containers[0].Image == "" {
			return fmt.Errorf("pod container image must be set")
		}
	}

	return nil
}

var _ = Describe("Disruption Controller", func() {
	var disruption *chaosv1beta1.Disruption

	BeforeEach(func() {
		disruption = &chaosv1beta1.Disruption{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: "default",
			},
			Spec: chaosv1beta1.DisruptionSpec{
				Selector: map[string]string{"foo": "bar"},
				NetworkFailure: &chaosv1beta1.NetworkFailureSpec{
					Hosts:       []string{"127.0.0.1"},
					Port:        80,
					Probability: 0,
					Protocol:    "tcp",
				},
				NetworkLatency: &chaosv1beta1.NetworkLatencySpec{
					Delay: 1000,
					Hosts: []string{"10.0.0.0/8"},
				},
			},
		}
	})

	AfterEach(func() {
		// delete disruption resource
		_ = k8sClient.Delete(context.Background(), disruption)
	})

	JustBeforeEach(func() {
		By("Creating disruption resource")
		Expect(k8sClient.Create(context.Background(), disruption)).To(BeNil())
	})

	Context("target all pods", func() {
		BeforeEach(func() {
			disruption.Spec.Count = -1
		})

		It("should target all the selected pods", func() {
			By("Ensuring that the inject pod has been created")
			Eventually(func() error { return expectChaosPod(disruption, chaostypes.PodModeInject, 4) }, timeout).Should(Succeed())

			By("Deleting the disruption resource")
			Expect(k8sClient.Delete(context.Background(), disruption)).To(BeNil())
			Eventually(func() error { return k8sClient.Get(context.Background(), instanceKey, disruption) }, timeout).Should(Succeed())

			By("Ensuring that the cleanup pod has been created")
			Eventually(func() error { return expectChaosPod(disruption, chaostypes.PodModeClean, 4) }, timeout).Should(Succeed())

			By("Simulating the completion of the cleanup pod by removing the finalizer")
			Eventually(func() error {
				if err := k8sClient.Get(context.Background(), instanceKey, disruption); err != nil {
					return err
				}
				disruption.ObjectMeta.Finalizers = []string{}
				return k8sClient.Update(context.Background(), disruption)
			}, timeout).Should(Succeed())

			By("Waiting for disruption resource to be deleted")
			Eventually(func() error { return k8sClient.Get(context.Background(), instanceKey, disruption) }, timeout).Should(MatchError("Disruption.chaos.datadoghq.com \"foo\" not found"))
		})
	})

	Context("target one pod only", func() {
		BeforeEach(func() {
			disruption.Spec.Count = 1
		})

		It("should target all the selected pods", func() {
			By("Ensuring that the inject pod has been created")
			Eventually(func() error { return expectChaosPod(disruption, chaostypes.PodModeInject, 2) }, timeout).Should(Succeed())

			By("Deleting the disruption resource")
			Expect(k8sClient.Delete(context.Background(), disruption)).To(BeNil())
			Eventually(func() error { return k8sClient.Get(context.Background(), instanceKey, disruption) }, timeout).Should(Succeed())

			By("Ensuring that the cleanup pod has been created")
			Eventually(func() error { return expectChaosPod(disruption, chaostypes.PodModeClean, 2) }, timeout).Should(Succeed())

			By("Simulating the completion of the cleanup pod by removing the finalizer")
			Eventually(func() error {
				if err := k8sClient.Get(context.Background(), instanceKey, disruption); err != nil {
					return err
				}
				disruption.ObjectMeta.Finalizers = []string{}
				return k8sClient.Update(context.Background(), disruption)
			}, timeout).Should(Succeed())

			By("Waiting for disruption resource to be deleted")
			Eventually(func() error { return k8sClient.Get(context.Background(), instanceKey, disruption) }, timeout).Should(MatchError("Disruption.chaos.datadoghq.com \"foo\" not found"))
		})
	})
})