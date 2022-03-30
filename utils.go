package lookup

import (
	"math"
	"strconv"
	"strings"
)

//I had planned to just do everything through maps, but because maps are unordered I'm just going to use a struct with corresponding arrays
//I also could have just sorted the map and printed directly from this function, but it felt unprofessional to not have a saved order
func SortBySpecificPrefix(prefixMap map[string]int, maxPrefix int) OrderedData {
	descendingOrderMap := OrderedData{}
	for interval := maxPrefix; interval >= 0; interval-- {
		for key, val := range prefixMap {
			prefixLength, err := strconv.Atoi(strings.Split(key, forwardslash)[1])
			if err != nil {
				continue
			}
			if prefixLength == interval {
				descendingOrderMap.prefixes = append(descendingOrderMap.prefixes, key)
				descendingOrderMap.asns = append(descendingOrderMap.asns, val)
			}
		}
	}
	return descendingOrderMap
}

func PowInt(x, y int) int {
	return int(math.Pow(float64(x), float64(y)))
}
