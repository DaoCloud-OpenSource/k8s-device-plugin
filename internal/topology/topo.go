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
	"encoding/json"
	"fmt"
	"strings"
)

type Matrix [][]float64

type Type string

// Collector define collect single node all gpu device topology matrix info.
type Collector interface {
	// Collect function return matrix info, if collect error, return error.
	Collect(ctx context.Context) (Matrix, error)
}

func (m Matrix) String() string {
	// 打印表头
	builder := strings.Builder{}
	builder.WriteString("\t")
	for i := 0; i < len(m); i++ {
		builder.WriteString(fmt.Sprintf("GPU%d\t", i))
	}
	builder.WriteString("\n")

	// 打印内容
	for i := 0; i < len(m); i++ {
		builder.WriteString(fmt.Sprintf("GPU%d\t", i))
		for j := 0; j < len(m[i]); j++ {
			builder.WriteString(fmt.Sprintf("%.2f\t", m[i][j]))
		}
		builder.WriteString("\n")
	}
	return builder.String()
}

func (m Matrix) Marshal() string {
	bytes, _ := json.Marshal(m)
	return string(bytes)
}
