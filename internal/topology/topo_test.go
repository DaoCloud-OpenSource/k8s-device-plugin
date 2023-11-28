package topology

import (
	"fmt"
	"testing"
)

func Test_String(t *testing.T) {
	m := Matrix{
		{6, 5, 3, 3, 1},
		{5, 6, 3, 3, 1},
		{3, 3, 6, 5, 1},
		{3, 3, 5, 6, 1},
		{1, 1, 1, 1, 6},
	}
	s := m.String()
	fmt.Println(s)

	m = Matrix{
		{735.29, 48.48, 24.25, 48.48},
		{48.48, 778.91, 48.48, 24.25},
		{24.25, 48.48, 779.3, 24.25},
		{48.48, 24.25, 24.25, 779.69},
	}
	s = m.String()
	fmt.Println(s)
}
