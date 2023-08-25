package main

import (
	"flag"
	"fmt"
	"whale/pkg/models"
	"whale/pkg/whaleconf"

	"gorm.io/gen"
)

func main() {
	conf := flag.String("conf", "conf.yaml", "Configuration file")
	initDB := flag.Bool("init", false, "Initialize the database")

	flag.Parse()

	g := gen.NewGenerator(gen.Config{
		OutPath: "./pkg/dbquery",
		Mode:    gen.WithQueryInterface,
	})

	fmt.Println("Generating schema code")
	m := []interface{}{
		// &models.HotTopic{},
		&models.HotTopicsInArea{},
		&models.CityTopics{},
		&models.MatchingReview{},
		&models.UserJoinTopic{},
		&models.RecentMatching{},
		&models.WhaleConfig{},

		&models.MatchingDurationConstraint{},
		&models.MatchingInvitation{},
		&models.MatchingQuota{},
		&models.MatchingResultConfirmAction{},
		&models.MatchingResult{},
		&models.Matching{},

		&models.Motion{},
		&models.MotionOfferRecord{},
		&models.MotionViewHistory{},
		&models.MotionReview{},
		&models.UserLikeMotion{},

		&models.DurationConstraint{},

		// 用户查看匹配
		&models.UserViewMatching{},
		&models.MatchingViewHistory{},
		&models.MatchingView{},

		// 用户匹配意向
		&models.MatchingOfferRecord{},
		&models.MatchingOfferSummary{},
		&models.UserThumbsUpMotion{},
	}

	g.ApplyBasic(m...)
	g.Execute()

	if *initDB {
		c := whaleconf.ReadConf(*conf)
		db := c.DB()
		if err := db.AutoMigrate(m...); err != nil {
			panic(err)
		}
	}
}
