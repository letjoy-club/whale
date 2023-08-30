package modelutil

import (
	"context"
	"whale/pkg/gqlient/hoopoe"
	"whale/pkg/loader"
	"whale/pkg/models"

	"github.com/letjoy-club/mida-tool/midacontext"
)

type MotionSummary struct {
	FemaleNum int `json:"femaleNum"`
	MaleNum   int `json:"maleNum"`

	Cities []CityMotions `json:"cities"`
}

type CityMotions struct {
	CityName  string `json:"cityName"`
	Code      string `json:"code"`
	FemaleNum int    `json:"femaleNum"`
	MaleNum   int    `json:"maleNum"`

	Categories []CityCategoryMotions `json:"categories"`
}

type CityCategoryMotions struct {
	ID           string `json:"id"`
	CategoryName string `json:"name"`
	FemaleNum    int    `json:"femaleNum"`
	MaleNum      int    `json:"maleNum"`

	Topics []CityTopicSummary `json:"topics"`
}

type CityTopicSummary struct {
	TopicName string `json:"name"`
	ID        string `json:"id"`

	FemaleNum int `json:"femaleNum"`
	MaleNum   int `json:"maleNum"`
}

func GenMotionSummary(ctx context.Context) (*MotionSummary, error) {
	loader := midacontext.GetLoader[loader.Loader](ctx).AllMotion
	if err := loader.Load(ctx); err != nil {
		return nil, err
	}
	citiesToMotion := loader.GetCityToMotions()

	topicIDToName := map[string]string{}
	categoryToName := map[string]string{}
	cityCodeToName := map[string]string{}

	{
		resp, err := hoopoe.GetAllTopicAndCategoryAndTopicName(ctx, midacontext.GetServices(ctx).Hoopoe, &hoopoe.GraphQLPaginator{
			Page: 1,
			Size: 9999,
		})
		if err != nil {
			return nil, err
		}

		for _, topic := range resp.Topics {
			topicIDToName[topic.Id] = topic.Name
		}
		for _, category := range resp.TopicCategories {
			categoryToName[category.Name] = category.Desc
		}
		for _, city := range resp.Areas {
			cityCodeToName[city.Code] = city.Name
		}
	}

	ret := &MotionSummary{}

	for cityId, cityCategoriesMotion := range citiesToMotion {
		cityMotion := CityMotions{
			Code:       cityId,
			CityName:   cityCodeToName[cityId],
			Categories: []CityCategoryMotions{},
		}

		cityFemaleNum := 0
		cityMaleNum := 0

		for categoryID, topicMotions := range cityCategoriesMotion {
			categorySummary := CityCategoryMotions{
				ID:           categoryID,
				CategoryName: categoryToName[categoryID],
				Topics:       []CityTopicSummary{},
			}

			femaleNum := 0
			maleNum := 0

			type TopicFemaleMale struct {
				Female int
				Male   int
			}

			topicFemaleMale := map[string]*TopicFemaleMale{}

			for _, topicMotion := range topicMotions {
				if _, ok := topicFemaleMale[topicMotion.TopicID]; !ok {
					topicFemaleMale[topicMotion.TopicID] = &TopicFemaleMale{}
				}
				if topicMotion.MyGender == models.GenderF.String() {
					femaleNum++
					topicFemaleMale[topicMotion.TopicID].Female++
				} else {
					maleNum++
					topicFemaleMale[topicMotion.TopicID].Male++
				}
			}

			for topicId, topicFM := range topicFemaleMale {
				categorySummary.Topics = append(categorySummary.Topics, CityTopicSummary{
					ID:        topicId,
					TopicName: topicIDToName[topicId],
					MaleNum:   topicFM.Male,
					FemaleNum: topicFM.Female,
				})
			}

			categorySummary.FemaleNum = femaleNum
			categorySummary.MaleNum = maleNum

			cityFemaleNum += femaleNum
			cityMaleNum += maleNum

			cityMotion.Categories = append(cityMotion.Categories, categorySummary)
		}

		cityMotion.FemaleNum = cityFemaleNum
		cityMotion.MaleNum = cityMaleNum

		ret.FemaleNum += cityFemaleNum
		ret.MaleNum += cityMaleNum
		ret.Cities = append(ret.Cities, cityMotion)
	}

	return ret, nil
}
