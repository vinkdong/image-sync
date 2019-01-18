package slice

type Runtime struct {
}

func Difference(sliceA []string, sliceB []string) []string {
	diff := make([]string, 0)
	diffMap := make(map[string]int)

	for _, v := range sliceA {
		diffMap[v] = 1
	}
	for _, v := range sliceB {
		diffMap[v] = diffMap[v] - 1
	}

	for k, v := range diffMap {
		if v > 0 {
			diff = append(diff, k)
		}
	}
	return diff
}
