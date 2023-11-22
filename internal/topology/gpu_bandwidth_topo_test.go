package topology

import (
	"context"
	"fmt"
	"testing"
)

func Test_Collect(t *testing.T) {
	b := gpuTopologyBandwidthCollector{
		ExecPath: "./testdata",
		ExecName: DefaultExecName,
	}
	matrix, err := b.Collect(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("\n%s", matrix.String())
	fmt.Println("Extracted GPU Connections Matrix:")
	for _, row := range matrix {
		fmt.Println(row)
	}
}
