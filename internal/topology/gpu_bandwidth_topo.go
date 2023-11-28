package topology

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"k8s.io/klog/v2"
)

const (
	DefaultExecName = "p2pBandwidthLatency"
	DefaultExecPath = "./"
)

func init() {
	// Make sure it is in the same order as the GPUs seen by nvidia-smi
	_ = os.Setenv("CUDA_DEVICE_ORDER", "PCI_BUS_ID")
}

var _ Collector = &gpuTopologyBandwidthCollector{}

func NewBandwidthGPUTopology() Collector {
	currentDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		klog.Errorf("get current process exec dir error: %+v", err)
		currentDir = DefaultExecPath
	}
	currentDir = "/root/go/src/github.com/DaoCloud-OpenSource/k8s-device-plugin/internal/topology/testdata"
	bandwidthExecPath := filepath.Join(currentDir, DefaultExecName)
	_, err = os.Stat(bandwidthExecPath)
	if err != nil {
		panic(err)
	}
	return &gpuTopologyBandwidthCollector{
		ExecName: DefaultExecName,
		ExecPath: currentDir,
	}
}

type gpuTopologyBandwidthCollector struct {
	ExecPath string
	ExecName string
}

func (g *gpuTopologyBandwidthCollector) Collect(ctx context.Context) (Matrix, error) {
	cmdPath := path.Join(g.ExecPath, g.ExecName)
	cmd := exec.Command(cmdPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}
	out := string(output)
	klog.V(5).Infof("exec gpu bandwidth out:\n%s", out)
	// 使用正则表达式匹配带宽矩阵的部分
	re := regexp.MustCompile(`Unidirectional P2P=Enabled Bandwidth \(P2P Writes\) Matrix \(GB/s\)\n([\s\S]+?)Bidirectional`)
	match := re.FindStringSubmatch(out)
	if len(match) != 2 {
		log.Println(out)
		return nil, fmt.Errorf("bandwidth matrix not found in output")
	}

	// 将匹配的部分按行拆分
	lines := strings.Split(match[1], "\n")

	// 解析每行的数字
	var matrix [][]float64
	for _, line := range lines {
		if strings.Contains(line, "D\\D") {
			continue // 跳过包含"D\D"的行
		}
		fields := strings.Fields(line)
		if len(fields) > 1 {
			var row []float64
			for _, field := range fields[1:] {
				value, err := strconv.ParseFloat(field, 64)
				if err != nil {
					return nil, err
				}
				row = append(row, value)
			}
			matrix = append(matrix, row)
		}
	}

	return matrix, nil
}

func getCurrentExecPath() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		klog.Error(err)
		return DefaultExecPath
	}
	return dir
}
