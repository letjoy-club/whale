package matcher

import (
	"time"
	"whale/pkg/gqlient/hoopoe"
	"whale/pkg/loader"
	"whale/pkg/models"
)

type EvaluatorResult struct {
	Score int

	TimeScore  int
	Properties []int

	FailedReason MatchingReason
}

func Evaluate(topicOption *hoopoe.TopicOptionConfigFields, matching1, matching2 *models.Matching) (result EvaluatorResult) {
	if topicOption == nil {
		result.Properties = make([]int, 0)
		return result
	}
	allWeight := topicOption.TimeWeight
	result.Properties = make([]int, len(topicOption.Properties))
	// 判断时间相似度
	timeScore := CompareDateSimilarity(matching1, matching2)
	score := 0

	if len(result.Properties) > 0 {

		for i, property := range topicOption.Properties {
			if !property.Enabled {
				result.Properties[i] = -1
				continue
			}
			allWeight += property.Weight
			// 判断数值类型的相似度属性
			if property.Comparable && len(property.Options) > 0 {
				v1 := matching1.GetSingleValueProperty(property.Name)
				v2 := matching2.GetSingleValueProperty(property.Name)
				score = CompareNumeralProperty(property.Options, v1, v2)
				result.Properties[i] = score
				continue
			}
			// 判断 enum 类型的相似度属性
			score = CompareEnumSimilarity(property, matching1, matching2)
			result.Properties[i] = score
		}
	}
	if allWeight > 0 {
		result.TimeScore = topicOption.TimeWeight * timeScore / allWeight
	}
	result.Score += result.TimeScore
	if len(result.Properties) > 0 {
		for i, p := range result.Properties {
			if p == -1 {
				continue
			}
			result.Score += topicOption.Properties[i].Weight * p / allWeight
		}
	}
	return result
}

func WeekOfDate(date string) time.Weekday {
	p, err := time.Parse("20060102", date)
	if err != nil {
		return -1
	}
	return p.Weekday()
}

var (
	suitableDateForWeekend          = []string{string(models.DatePeriodWeekend), string(models.DatePeriodWeekendNight), string(models.DatePeriodWeekendAfternoon)}
	suitableDateForWeekendNight     = []string{string(models.DatePeriodWeekendNight), string(models.DatePeriodWeekend)}
	suitableDateForWeekendAfternoon = []string{string(models.DatePeriodWeekendAfternoon), string(models.DatePeriodWeekend)}

	suitableDateForWorkday          = []string{string(models.DatePeriodWorkday), string(models.DatePeriodWorkdayNight), string(models.DatePeriodWorkdayAfternoon)}
	suitableDateForWorkdayNight     = []string{string(models.DatePeriodWorkdayNight), string(models.DatePeriodWorkday)}
	suitableDateForWorkdayAfternoon = []string{string(models.DatePeriodWorkdayAfternoon), string(models.DatePeriodWorkday)}
)

func SuitablePeriods(m *models.Matching, periods []string) map[string]struct{} {
	suitablePeriods := map[string]struct{}{}
	for _, period := range periods {
		var toBeAdded []string
		if period == models.DatePeriodWeekend.String() {
			toBeAdded = suitableDateForWeekend
		} else if period == models.DatePeriodWeekendNight.String() {
			toBeAdded = suitableDateForWeekendNight
		} else if period == models.DatePeriodWeekendAfternoon.String() {
			toBeAdded = suitableDateForWeekendAfternoon
		} else if period == models.DatePeriodWorkday.String() {
			toBeAdded = suitableDateForWorkday
		} else if period == models.DatePeriodWorkdayNight.String() {
			toBeAdded = suitableDateForWorkdayNight
		} else if period == models.DatePeriodWorkdayAfternoon.String() {
			toBeAdded = suitableDateForWorkdayAfternoon
		}
		if len(toBeAdded) == 0 {
			continue
		}
		for _, p := range toBeAdded {
			suitablePeriods[p] = struct{}{}
		}
	}
	return suitablePeriods
}

func CompareDateSimilarity(matching1, matching2 *models.Matching) (score int) {
	if len(matching1.DayRange) != 0 {
		// matching1 有确定的日期，直接比较
		for _, day := range matching1.DayRange {
			if matching2.HasSpecificDay(day) {
				return 100
			}
		}
	} else {
		// matching1 没有确定的日期，比较时间段
		if len(matching1.PreferredPeriods) != 0 && len(matching2.PreferredPeriods) != 0 {
			// matching1 和 matching2 都有确定的星期
			suitablePeriods := SuitablePeriods(matching1, matching1.PreferredPeriods)
			for _, period := range matching2.PreferredPeriods {
				if _, ok := suitablePeriods[period]; ok {
					return 100
				}
			}
		}
		// 其中一个没有确定的日期，但是被认为全选
		return 100
	}

	if len(matching2.DayRange) != 0 {
		// matching2 有确定的日期，直接比较
		for _, day := range matching2.DayRange {
			if matching1.HasSpecificDay(day) {
				return 100
			}
		}
	} else {
		// matching2 没有确定的日期，比较时间段
		if len(matching1.PreferredPeriods) != 0 && len(matching2.PreferredPeriods) != 0 {
			// matching1 和 matching2 都有确定的星期
			suitablePeriods := SuitablePeriods(matching2, matching2.PreferredPeriods)
			for _, period := range matching1.PreferredPeriods {
				if _, ok := suitablePeriods[period]; ok {
					return 100
				}
			}
		}
		// 其中一个没有确定的日期，但是被认为全选
		return 100
	}

	return 0
}

func CompareEnumSimilarity(property *loader.TopicOptionConfigProperty, matching1, matching2 *models.Matching) (score int) {
	v1 := matching1.GetProperty(property.Id)
	v2 := matching2.GetProperty(property.Id)

	// 如果 matching 未填属性
	if len(v1) == 0 || len(v2) == 0 {
		// 如果默认是全选，则认为该属性匹配度为 100 (matching 一定互相包容)
		if property.DefaultSelectAll {
			return 100
		} else {
			// 如果两个都是未填，则认为该属性匹配度为 100
			if len(v1) == 0 || len(v2) == 0 {
				return 100
			}
			// 这时候选得越多越不相似，每出现一个选项，扣 10 分
			return max(0, 100-10*(abs(len(v1)-len(v2))))
		}
	}

	sameCount := 0
	for _, item1 := range v1 {
		for _, item2 := range v2 {
			if item1 == item2 {
				sameCount++
			}
		}
	}
	allItems := len(v1) + len(v2)
	// 如果 matching 填了属性
	if property.DefaultSelectAll {
		// 如果默认是全选，则是按照相同的选项数量来计算分数
		return sameCount * 2 * 100 / allItems
	} else {
		// 如果默认不是全选，则是按照不同的选项数量来计算分数
		uniqueItems := allItems - sameCount*2
		return max(100-10*uniqueItems, 0)
	}
}

// CompareNumeralProperty 返回数值类型的属性的匹配分数，范围为0-100
func CompareNumeralProperty(options []*loader.TopicOptionConfigOptions, label1, label2 string) (score int) {
	// 如果配置属性的选项数量小于2，则认为该配置为及格分 60
	if len(options) < 2 {
		return 60
	}

	value1 := 0
	value2 := 0
	for _, option := range options {
		if option.Name == label1 {
			value1 = option.Value
		}
		if option.Name == label2 {
			value2 = option.Value
		}
	}

	// |100 - distance|
	return min(100, abs(100-abs(value1-value2)))
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func abs(a int) int {
	if a > 0 {
		return a
	}
	return -a
}
