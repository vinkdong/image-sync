package random

import (
	"testing"
	"time"
	"fmt"
)

func TestRangeInt(t *testing.T)  {
	Seed(time.Now().UnixNano())
	i := RangeInt(1,80)
	checkInt(i,1,80,t)
}

func checkInt(i, min, max int, t *testing.T) bool {
	fmt.Println(i, min, max)
	if i >= min && i <= max {
		return true
	}
	t.Fatalf("check int range faile %d not in  min: %d, max:%d", i, min, max)
	return false
}

func TestRangeIntWithExclude(t *testing.T) {

	Seed(time.Now().UnixNano())
	i := RangeIntInclude(Slice{2,79})

	checkInt(i,2,79,t)

	i = RangeIntInclude(Slice{2,4})
	checkInt(i,2,4,t)

	i = RangeIntInclude(Slice{2,4},Slice{3,8})
	checkInt(i,2,8,t)
}
