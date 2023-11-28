/*
Copyright 2022 The Kangaroo Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package topology

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/NVIDIA/go-nvml/pkg/nvml"
	"k8s.io/klog/v2"
)

const (
	// TopologyNv is gpy link type ls nvlink
	TopologyNv nvml.GpuTopologyLevel = 60
)

func NewLinkCollector() Collector {
	return &gpuTopologyLinkCollector{}
}

var _ Collector = &gpuTopologyLinkCollector{}

type gpuTopologyLinkCollector struct {
}

func (g *gpuTopologyLinkCollector) Collect(ctx context.Context) (Matrix, error) {
	result := nvml.Init()
	if result != nvml.SUCCESS {
		return nil, fmt.Errorf("unable to initialize NVML: %v", nvml.ErrorString(result))
	}
	defer func() {
		result = nvml.Shutdown()
		if result != nvml.SUCCESS {
			log.Printf("unable to shutdown NVML: %v", nvml.ErrorString(result))
		}
	}()
	matrix, err := g.linkTypeTopology()
	if err != nil {
		return nil, err
	}
	return g.convertToMatrix(matrix), nil
}

func (g *gpuTopologyLinkCollector) printTopology(matrix [][]LinkType) {
	log.Println("gpu topology link type:")
	// 打印表头
	fmt.Print("\t")
	for i := 0; i < len(matrix); i++ {
		fmt.Printf("GPU%d\t", i)
	}
	fmt.Println()

	// 打印内容
	for i := 0; i < len(matrix); i++ {
		fmt.Printf("GPU%d\t", i)
		for j := 0; j < len(matrix[i]); j++ {
			fmt.Printf("%s\t", LinkTypeToStr(matrix[i][j]))
		}
		fmt.Println()
	}
}

func (g *gpuTopologyLinkCollector) convertToMatrix(matrix [][]LinkType) Matrix {
	g.printNvidiaSmiTopologyMatrix()
	g.printTopology(matrix)
	resultMatrix := make(Matrix, len(matrix))
	for i := range matrix {
		resultMatrix[i] = make([]float64, len(matrix))
		for j := range matrix[i] {
			resultMatrix[i][j] = float64(matrix[i][j])
		}
	}
	return resultMatrix
}

func (g *gpuTopologyLinkCollector) printNvidiaSmiTopologyMatrix() {
	cmd := exec.Command("nvidia-smi", "topo", "-m")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return
	}
	klog.V(5).Infof("nvidia-smi topo -m:\n%s", string(output))
}

func (g *gpuTopologyLinkCollector) linkTypeTopology() ([][]LinkType, error) {
	cmd := exec.Command("nvidia-smi", "topo", "-m")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}
	output := string(out)
	count, ret := nvml.DeviceGetCount()
	if ret != nvml.SUCCESS {
		return nil, errors.New(nvml.ErrorString(ret))
	}

	linkTypes := make([][]LinkType, count)
	lines := strings.Split(output, "\n")
	// 提取连接矩阵
	for i := 0; i < count; i++ {
		line := lines[i+1]
		fileds := strings.Split(line, "\t")
		linkTypes[i] = make([]LinkType, count)
		for j := 0; j < count; j++ {
			linkTypes[i][j] = StrToLinkType(fileds[j+1])
		}
	}
	return linkTypes, nil
}

// collectDeviceToDeviceLinkType is get gpu to gpu link type
//func (g *gpuTopologyLinkCollector) collectDeviceToDeviceLinkType(src, dest nvml.Device) (LinkType, nvml.Return) {
//	srcUUID, ret := src.GetUUID()
//	if ret != nvml.SUCCESS {
//		return P2PLinkUnknown, ret
//	}
//	destUUID, ret := dest.GetUUID()
//	if ret != nvml.SUCCESS {
//		return P2PLinkUnknown, ret
//	}
//	if srcUUID == destUUID {
//		return P2PLinkSameBoard, nvml.SUCCESS
//	}
//	linkType, err := getP2PLink(src, dest)
//	if err != nvml.SUCCESS {
//		return linkType, err
//	}
//	nvLink, err := getNVLink(src, dest)
//	if err != nvml.SUCCESS {
//		klog.Errorf("get nvlink type result is not success, %s", nvml.ErrorString(err))
//		return linkType, nvml.SUCCESS
//	}
//	return nvLink, nvml.SUCCESS
//}
//
//func getP2PLink(dev1, dev2 nvml.Device) (link LinkType, err nvml.Return) {
//	level, ret := nvml.DeviceGetTopologyCommonAncestor(dev1, dev2)
//	if ret != nvml.SUCCESS {
//		return P2PLinkUnknown, ret
//	}
//	switch level {
//	case nvml.TOPOLOGY_INTERNAL:
//		link = P2PLinkSameBoard
//	case nvml.TOPOLOGY_SINGLE:
//		link = P2PLinkSingleSwitch
//	case nvml.TOPOLOGY_MULTIPLE:
//		link = P2PLinkMultiSwitch
//	case nvml.TOPOLOGY_HOSTBRIDGE:
//		link = P2PLinkHostBridge
//	case nvml.TOPOLOGY_NODE:
//		link = P2PLinkSameCPU
//	case nvml.TOPOLOGY_SYSTEM:
//		link = P2PLinkCrossCPU
//	default:
//		err = nvml.ERROR_NOT_SUPPORTED
//	}
//	return
//}
//
//func deviceGetAllNvLinkRemotePciInfo(dev nvml.Device) ([]string, nvml.Return) {
//	var busIds []string
//	//busIdMap := make(map[string]string)
//	// TODO int maxNvLinks = (sm < 60) ? 0 : (sm < 70) ? 4 : (sm < 80) ? 6 : (sm < 90) ? 12 : 18;
//	for i := 0; i < nvml.NVLINK_MAX_LINKS; i++ {
//		nvLinkState, ret := nvml.DeviceGetNvLinkState(dev, i)
//		if ret != nvml.SUCCESS {
//			if len(busIds) > 0 {
//				return busIds, nvml.SUCCESS
//			}
//			log.Printf("nvml call DeviceGetNvLinkState ret is: %d", ret)
//			return nil, ret
//		}
//		if nvLinkState == nvml.FEATURE_ENABLED {
//			remotePciInfo, ret := nvml.DeviceGetNvLinkRemotePciInfo(dev, i)
//			if ret != nvml.SUCCESS {
//				log.Printf("nvml call DeviceGetNvLinkRemotePciInfo ret is: %d", ret)
//				return nil, ret
//			}
//
//			busId := fmt.Sprintf("%c", remotePciInfo.BusId[0])
//			busIds = append(busIds, busId)
//			//if _, ok := busIdMap[busId]; !ok {
//			//	busIdMap[busId] = busId
//			//}
//		}
//	}
//	return busIds, nvml.ERROR_NOT_SUPPORTED
//}
//
//func getNVLink(dev1, dev2 nvml.Device) (link LinkType, err nvml.Return) {
//	nvBusIds1, err := deviceGetAllNvLinkRemotePciInfo(dev1)
//	if err != nvml.SUCCESS || nvBusIds1 == nil {
//		return P2PLinkUnknown, err
//	}
//	pciInfo, _ := dev2.GetPciInfo()
//	busId := fmt.Sprintf("%c", pciInfo.BusId[0])
//	nvLink := P2PLinkUnknown
//	for _, nvBusID := range nvBusIds1 {
//		if nvBusID == busId {
//			switch nvLink {
//			case P2PLinkUnknown:
//				nvLink = SingleNVLINKLink
//			case SingleNVLINKLink:
//				nvLink = TwoNVLINKLinks
//			case TwoNVLINKLinks:
//				nvLink = ThreeNVLINKLinks
//			case ThreeNVLINKLinks:
//				nvLink = FourNVLINKLinks
//			case FourNVLINKLinks:
//				nvLink = FiveNVLINKLinks
//			case FiveNVLINKLinks:
//				nvLink = SixNVLINKLinks
//			case SixNVLINKLinks:
//				nvLink = SevenNVLINKLinks
//			case SevenNVLINKLinks:
//				nvLink = EightNVLINKLinks
//			case EightNVLINKLinks:
//				nvLink = NineNVLINKLinks
//			case NineNVLINKLinks:
//				nvLink = TenNVLINKLinks
//			case TenNVLINKLinks:
//				nvLink = ElevenNVLINKLinks
//			case ElevenNVLINKLinks:
//				nvLink = TwelveNVLINKLinks
//			}
//		}
//	}
//	// TODO: Handle NVSwitch semantics
//	return nvLink, nvml.SUCCESS
//}
