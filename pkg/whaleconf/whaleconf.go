package whaleconf

import (
	"fmt"
	"log"
	"os"

	sql "github.com/go-sql-driver/mysql"
	"github.com/goccy/go-yaml"
	"github.com/redis/go-redis/v9"
	mysql "gorm.io/driver/mysql"
	"gorm.io/gorm/schema"

	"gorm.io/gorm"
)

type Conf struct {
	DSN      string `yaml:"dsn"`
	DBPrefix string `yaml:"db-prefix"`

	RedisSetting RedisConf `yaml:"redis"`
	Secret       string    `yaml:"secret"`

	ServicesSetting ServiceConf `yaml:"services"`
}

type ServiceConf struct {
	// 基础服务
	Hoopoe string `yaml:"hoopoe"`
	// IM 服务
	Smew string `yaml:"smew"`
}

func (l Conf) Redis() *redis.Client {
	fmt.Println("using redis", l.RedisSetting.Address)
	client := redis.NewClient(&redis.Options{
		Addr:     l.RedisSetting.Address,
		Password: l.RedisSetting.Password,
		DB:       l.RedisSetting.DB,
	})
	return client
}

func (c *Conf) DB() *gorm.DB {
	dsnObj, err := sql.ParseDSN(c.DSN)
	if err != nil {
		log.Panic(err)
	}
	dsnObj.ParseTime = true
	dsnObj.Collation = "utf8mb4_general_ci"
	dsn := dsnObj.FormatDSN() + "&charset=utf8mb4&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: &schema.NamingStrategy{TablePrefix: c.DBPrefix},
	})
	if err != nil {
		panic(err)
	}
	db = db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci")
	return db
}

type RedisConf struct {
	Address  string `yaml:"address"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

func ReadConf(path string) *Conf {
	var conf Conf
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Failed to read conf, generate a new config")
		WriteDefaultConf(path)
		panic(err)
	}
	defer file.Close()
	if err := yaml.NewDecoder(file).Decode(&conf); err != nil {
		panic(err)
	}
	return &conf
}

func DefaultConf() Conf {
	return Conf{
		DSN:      "staging_user:1Vg4UwRCRFiD@tcp(localhost:3306)/youyue_staging",
		DBPrefix: "whale_",
		RedisSetting: RedisConf{
			Address:  "127.0.0.1:6379",
			Password: "letjoy2023",
			DB:       1,
		},
		ServicesSetting: ServiceConf{
			Hoopoe: "http://localhost:11011",
			Smew:   "http://localhost:11012",
		},
		Secret: "letjoy",
	}
}

func WriteDefaultConf(path string) {
	conf := DefaultConf()
	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	if err := yaml.NewEncoder(file).Encode(&conf); err != nil {
		panic(err)
	}
}
