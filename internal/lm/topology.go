package lm

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"time"

	spec "github.com/DaoCloud-OpenSource/k8s-device-plugin/api/config/v1"
	"github.com/DaoCloud-OpenSource/k8s-device-plugin/internal/resource"
	"github.com/DaoCloud-OpenSource/k8s-device-plugin/internal/rm"
	"github.com/DaoCloud-OpenSource/k8s-device-plugin/internal/topology"

	"github.com/NVIDIA/go-nvlib/pkg/nvlib/device"
	"github.com/NVIDIA/go-nvlib/pkg/nvml"
	"k8s.io/klog/v2"
	pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
)

var once sync.Once

func NewTopologyLabel(manager resource.Manager, config *spec.Config, refreshChan chan<- struct{}) (Labeler, error) {
	linkTypeCollect := topology.NewLinkCollector()
	bandwidthCollect := topology.NewBandwidthGPUTopology()
	nvmllib := nvml.New()
	devicelib := device.New(device.WithNvml(nvmllib))

	once.Do(func() {
		config.Resources.AddGPUResource("*", "gpu")
		rms, err := rm.NewNVMLResourceManagers(nvmllib, config)
		if err != nil {
			klog.Errorf("new nvml resource manager error: %+v", err)
			return
		}
		health := make(chan *rm.Device)
		stop := make(chan interface{})
		go func() {
			err := rms[0].CheckHealth(stop, health)
			if err != nil {
				klog.Infof("Failed to start health check: %v; continuing with health checks disabled", err)
			}
		}()
		go func() {
			for {
				select {
				case <-stop:
					return
				case d := <-health:
					// FIXME: there is no way to recover from the Unhealthy state.
					d.Health = pluginapi.Unhealthy
					klog.Infof("'%s' device marked unhealthy: %s", rms[0].Resource(), d.ID)
					refreshChan <- struct{}{}
				}
			}
		}()
	})

	return &topologyLabel{
		linkTypeCollect:  linkTypeCollect,
		bandwidthCollect: bandwidthCollect,
		manager:          manager,
		config:           config,
		devicelib:        devicelib,
	}, nil
}

type topologyLabel struct {
	linkTypeCollect  topology.Collector
	bandwidthCollect topology.Collector
	manager          resource.Manager
	config           *spec.Config
	devicelib        device.Interface
}

type DeviceInfo struct {
	ID      string // device UUID
	Healthy bool   // Device healthy or not
}

// Labels
// current don't patch node annotation, only patch node labels, waite this pr release. in the before, use client to patch node
// https://github.com/kubernetes-sigs/node-feature-discovery/pull/1417/files
func (t *topologyLabel) Labels() (Labels, error) {
	interval := time.Now().Add(time.Duration(*t.config.Flags.GFD.SleepInterval))
	ctx, _ := context.WithDeadline(context.Background(), interval)
	linkTypeMatrix, err := t.linkTypeCollect.Collect(ctx)
	if err != nil {
		klog.Errorf("collect gpu link type topology matrix error: %+v", err)
		return nil, err
	}
	klog.Infof("collect gpu link type matrix:\n%s", linkTypeMatrix.String())

	bandwidthMatrix, err := t.bandwidthCollect.Collect(ctx)
	if err != nil {
		klog.Errorf("collect gpu bandwidth type topology matrix error: %+v", err)
		return nil, err
	}
	klog.Infof("collect gpu bandwidth type matrix:\n%s", bandwidthMatrix.String())

	devices, err := t.collectDevices()
	if err != nil {
		klog.Errorf("collect device error: %+v", err)
		return nil, err
	}
	klog.Infof("collect gpu device info:\n%v", devices)

	l := Labels{
		"node.gpu.info/topology":  linkTypeMatrix.Marshal(),
		"node.gpu.info/bandwidth": bandwidthMatrix.Marshal(),
		"node.gpu.info/devices":   devices,
	}
	return l, nil
}

func (t *topologyLabel) collectDevices() (string, error) {
	devices, err := t.devicelib.GetDevices()
	if err != nil {
		klog.Errorf("get all devices error: %+v", err)
		return "", err
	}
	deviceInfos := make([]DeviceInfo, 0)
	for _, device := range devices {
		uuid, ret := device.GetUUID()
		if ret != nvml.SUCCESS {
			return "", errors.New(ret.Error())
		}
		deviceInfos = append(deviceInfos, DeviceInfo{
			ID:      uuid,
			Healthy: true,
		})
	}
	values, err := json.Marshal(deviceInfos)
	if err != nil {
		return "", err
	}
	return string(values), nil
}
