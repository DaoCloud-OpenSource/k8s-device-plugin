package lm

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	spec "github.com/DaoCloud-OpenSource/k8s-device-plugin/api/config/v1"
	k8s "github.com/DaoCloud-OpenSource/k8s-device-plugin/internal/kubernetes"
	"github.com/DaoCloud-OpenSource/k8s-device-plugin/internal/resource"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"
)

func NewAnnotation(manager resource.Manager, config *spec.Config, refreshChan chan<- struct{}) (Labeler, error) {
	if config.Flags.TopologyEnabled != nil && *config.Flags.TopologyEnabled {
		topologyLabeler, err := NewTopologyLabel(manager, config, refreshChan)
		if err != nil {
			return nil, fmt.Errorf("error creating Topology labeler: %v", err)
		}
		l := Merge(topologyLabeler)
		return l, nil
	}
	return nil, errors.New("topology don't enable")
}

// TODO when upgrade nfd version, remove this method, current nfd can't update node annotation, but mina code having this function
func (labels Labels) UpdateNodeAnnotation() error {
	clientSet, err := k8s.GetKubernetesClientSet()
	if err != nil {
		return fmt.Errorf("failed to get Kubernetes client: %v", err)
	}
	patchAnnotations := map[string]interface{}{"metadata": map[string]map[string]string{"annotations": labels}}
	b, err := json.Marshal(patchAnnotations)
	if err != nil {
		return err
	}
	nodeName := k8s.NodeName()
	_, err = clientSet.CoreV1().Nodes().Patch(context.TODO(), nodeName, types.StrategicMergePatchType, b, metav1.PatchOptions{})
	if err != nil {
		klog.Errorf("Failed to update node annotation, lables is %+v, error %v.", labels, err)
	} else {
		klog.Info("Updated node annotation successfully.")
	}
	return err
}
