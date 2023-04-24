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
	m := []interface{}{&models.Matching{}, &models.MatchingResult{}, &models.MatchingResultConfirmAction{}, &models.MatchingQuota{}, &models.MatchingInvitation{}}

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
