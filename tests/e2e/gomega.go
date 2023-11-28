/*
 * Copyright (c) 2023, NVIDIA CORPORATION.  All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package e2e

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/DaoCloud-OpenSource/k8s-device-plugin/tests/e2e/common"
	. "github.com/onsi/gomega"
	gomegatypes "github.com/onsi/gomega/types"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"
	e2elog "k8s.io/kubernetes/test/e2e/framework"
)

type k8sLabels map[string]string

// eventuallyNonControlPlaneNodes is a helper for asserting node properties
func eventuallyNonControlPlaneNodes(ctx context.Context, cli clientset.Interface) AsyncAssertion {
	return Eventually(func(g Gomega, ctx context.Context) ([]corev1.Node, error) {
		return common.GetNonControlPlaneNodes(ctx, cli)
	}).WithPolling(1 * time.Second).WithTimeout(10 * time.Second).WithContext(ctx)
}

// MatchLabels returns a specialized Gomega matcher for checking if a list of
// nodes are labeled as expected.
func MatchLabels(expectedNew map[string]k8sLabels, oldNodes []corev1.Node) gomegatypes.GomegaMatcher {
	return &nodeListPropertyRegexpMatcher[k8sLabels]{
		propertyName: "labels",
		expected:     expectedNew,
		oldNodes:     oldNodes,
	}
}

// nodeListPropertyRegexpMatcher is a generic Gomega matcher for asserting one property a group of nodes.
type nodeListPropertyRegexpMatcher[T any] struct {
	expected map[string]k8sLabels
	oldNodes []corev1.Node

	propertyName string
	node         *corev1.Node //nolint:unused
	missing      []string     //nolint:unused
	invalidValue []string     //nolint:unused
}

// Match method of the GomegaMatcher interface.
func (m *nodeListPropertyRegexpMatcher[T]) Match(actual interface{}) (bool, error) {
	nodes, ok := actual.([]corev1.Node)
	if !ok {
		return false, fmt.Errorf("expected []corev1.Node, got: %T", actual)
	}

	targetNode := corev1.Node{}
	for _, node := range nodes {
		_, ok := m.expected[node.Name]
		if !ok {
			e2elog.Logf("Skipping node %q as no expected was specified", node.Name)
			continue
		}
		targetNode = node
		break
	}

	for labelKey, labelValue := range m.expected[targetNode.Name] {
		// missing key
		if _, ok := targetNode.Labels[labelKey]; !ok {
			m.missing = append(m.missing, labelKey)
			continue
		}
		// invalid value
		regexMatcher := regexp.MustCompile(labelValue)
		if !regexMatcher.MatchString(targetNode.Labels[labelKey]) {
			m.invalidValue = append(m.invalidValue, fmt.Sprintf("%s: %s", labelKey, targetNode.Labels[labelKey]))
			return false, fmt.Errorf("node %q label %q value %q does not match %q", targetNode.Name, labelKey, targetNode.Labels[labelKey], labelValue)
		}
	}

	return true, nil
}

// FailureMessage method of the GomegaMatcher interface.
func (m *nodeListPropertyRegexpMatcher[T]) FailureMessage(actual interface{}) string {
	return m.message()
}

// NegatedFailureMessage method of the GomegaMatcher interface.
func (m *nodeListPropertyRegexpMatcher[T]) NegatedFailureMessage(actual interface{}) string {
	return fmt.Sprintf("Node %q matched unexpectedly", m.node.Name)
}

// TODO remove nolint when golangci-lint is able to cope with generics
//
//nolint:unused
func (m *nodeListPropertyRegexpMatcher[T]) message() string {
	msg := fmt.Sprintf("Node %q %s did not match:", m.node.Name, m.propertyName)
	if len(m.missing) > 0 {
		msg += fmt.Sprintf("\n  missing:\n    %s", strings.Join(m.missing, "\n    "))
	}
	if len(m.invalidValue) > 0 {
		msg += fmt.Sprintf("\n  invalid value:\n    %s", strings.Join(m.invalidValue, "\n    "))
	}
	return msg
}

// JobIsCompleted checks if a job is completed
func JobIsCompleted(ctx context.Context, cli clientset.Interface, namespace, podName string) bool {
	pod, err := cli.CoreV1().Pods(namespace).Get(ctx, podName, metav1.GetOptions{})
	if err != nil {
		return false
	}
	// Check if the pod's phase is Succeeded.
	if pod.Status.Phase == "Succeeded" {
		return true
	}
	return false
}
