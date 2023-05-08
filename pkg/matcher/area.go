package matcher

import "github.com/samber/lo"

func AreaOverlapse(areaIDs1 []string, areaIDs2 []string) (overlaps []string, hasOverlapse bool) {
	if len(areaIDs1) == 0 && len(areaIDs2) != 0 {
		return areaIDs2, true
	}
	if len(areaIDs2) == 0 && len(areaIDs1) != 0 {
		return areaIDs1, true
	}
	if len(areaIDs1) == 0 && len(areaIDs2) == 0 {
		return []string{}, true
	}
	overlaps = lo.Intersect(areaIDs1, areaIDs2)
	return overlaps, len(overlaps) != 0
}
