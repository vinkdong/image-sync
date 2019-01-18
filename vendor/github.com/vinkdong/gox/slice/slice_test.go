package slice

import (
	"testing"
	"fmt"
)

func TestDifference(t *testing.T) {
	sla := []string{"a", "b", "c", "f"}
	slb := []string{"b", "c", "d", "e"}
	diff := Difference(sla, slb)
	fmt.Println(diff)
	if len(diff) != 2 {
		t.Error("slice difference set error")
	}
	diff2 := Difference(diff, sla)
	if len(diff2) != 0 {
		t.Error("slice difference set error")
	}
}
