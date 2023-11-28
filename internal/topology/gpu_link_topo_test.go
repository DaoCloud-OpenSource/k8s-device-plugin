package topology

import (
	"context"
	"fmt"
	"testing"
)

func Test_printTopology(t *testing.T) {
	inputMatrix := [][]LinkType{
		{P2PLinkSameBoard, P2PLinkSingleSwitch, P2PLinkHostBridge, P2PLinkHostBridge, P2PLinkCrossCPU},
		{P2PLinkSingleSwitch, P2PLinkSameBoard, P2PLinkHostBridge, P2PLinkHostBridge, P2PLinkCrossCPU},
		{P2PLinkHostBridge, P2PLinkHostBridge, P2PLinkSameBoard, P2PLinkSingleSwitch, P2PLinkCrossCPU},
		{P2PLinkHostBridge, P2PLinkHostBridge, P2PLinkSingleSwitch, P2PLinkSameBoard, P2PLinkCrossCPU},
		{P2PLinkCrossCPU, P2PLinkCrossCPU, P2PLinkCrossCPU, P2PLinkCrossCPU, P2PLinkSameBoard},
	}
	g := gpuTopologyLinkCollector{}
	g.printTopology(inputMatrix)
}

func Test_convertToMatrix(t *testing.T) {
	inputMatrix := [][]LinkType{
		{P2PLinkSameBoard, P2PLinkSingleSwitch, P2PLinkHostBridge, P2PLinkHostBridge, P2PLinkCrossCPU},
		{P2PLinkSingleSwitch, P2PLinkSameBoard, P2PLinkHostBridge, P2PLinkHostBridge, P2PLinkCrossCPU},
		{P2PLinkHostBridge, P2PLinkHostBridge, P2PLinkSameBoard, P2PLinkSingleSwitch, P2PLinkCrossCPU},
		{P2PLinkHostBridge, P2PLinkHostBridge, P2PLinkSingleSwitch, P2PLinkSameBoard, P2PLinkCrossCPU},
		{P2PLinkCrossCPU, P2PLinkCrossCPU, P2PLinkCrossCPU, P2PLinkCrossCPU, P2PLinkSameBoard},
	}
	g := gpuTopologyLinkCollector{}
	matrix := g.convertToMatrix(inputMatrix)
	t.Logf("%+v", matrix)
}

func Test_LinkCollect(t *testing.T) {
	g := gpuTopologyLinkCollector{}
	matrix, err := g.Collect(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Extracted GPU Connections Matrix:")
	for _, row := range matrix {
		fmt.Println(row)
	}
}
