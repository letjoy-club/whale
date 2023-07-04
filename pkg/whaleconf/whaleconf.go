package whaleconf

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	sql "github.com/go-sql-driver/mysql"
	"github.com/goccy/go-yaml"
	"github.com/letjoy-club/mida-tool/pulsarutil"
	"github.com/redis/go-redis/v9"
	cls "github.com/tencentcloud/tencentcloud-cls-sdk-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	mysql "gorm.io/driver/mysql"
	"gorm.io/gorm/schema"

	"gorm.io/gorm"
)

type Conf struct {
	DSN      string `yaml:"dsn"`
	DBPrefix string `yaml:"db-prefix"`

	RedisConf RedisConf `yaml:"redis"`
	Secret    string    `yaml:"secret"`

	MQ pulsarutil.PulsarClient `yaml:"mq"`

	ServiceConf ServiceConf `yaml:"services"`
	TraceConf   TraceConf   `yaml:"trace"`
	QCloud      QCloudConf  `yaml:"qcloud"`
}

type QCloudConf struct {
	CLSConf CLSConf `yaml:"cls"`
}

type CLSConf struct {
	Endpoint  string `yaml:"endpoint"`
	SecretID  string `yaml:"secret-id"`
	SecretKey string `yaml:"secret-key"`
	TopicID   string `yaml:"topic-id"`
}

type TraceConf struct {
	Jaeger string `yaml:"jaeger"`
	Token  string `yaml:"token"`
}

type ServiceConf struct {
	// 基础服务
	Hoopoe string `yaml:"hoopoe"`
	// IM 服务
	Smew string `yaml:"smew"`
	// 通知服务
	Scream string `yaml:"scream"`
}

func (l Conf) CLS() *cls.AsyncProducerClient {
	c := cls.GetDefaultAsyncProducerClientConfig()
	c.Endpoint = l.QCloud.CLSConf.Endpoint
	c.AccessKeyID = l.QCloud.CLSConf.SecretID
	c.AccessKeySecret = l.QCloud.CLSConf.SecretKey
	c.Retries = 1
	client, err := cls.NewAsyncProducerClient(c)
	if err != nil {
		panic(err)
	}
	return client
}

func (l Conf) Trace() *sdktrace.TracerProvider {
	ctx := context.Background()
	fmt.Println(" - trace:", l.TraceConf.Jaeger)

	opts := []otlptracegrpc.Option{
		otlptracegrpc.WithEndpoint(l.TraceConf.Jaeger),
		otlptracegrpc.WithInsecure(),
	}

	res, err := resource.New(ctx,
		//设置 Token 值
		resource.WithAttributes(attribute.KeyValue{
			Key: "token", Value: attribute.StringValue(l.TraceConf.Token),
		}),
		//设置服务名
		resource.WithAttributes(attribute.KeyValue{
			Key: "service.name", Value: attribute.StringValue("whale-matching"),
		}),
	)
	if err != nil {
		panic(err)
	}

	exporter, err := otlptracegrpc.New(ctx, opts...)
	if err != nil {
		panic(err)
	}

	//创建新的TracerProvider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter, sdktrace.WithBatchTimeout(10*time.Second)),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp
}

func (l Conf) Redis() *redis.Client {
	fmt.Println("using redis", l.RedisConf.Address)
	client := redis.NewClient(&redis.Options{
		Addr:     l.RedisConf.Address,
		Password: l.RedisConf.Password,
		DB:       l.RedisConf.DB,
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
	{
		sql, _ := db.DB()
		sql.SetMaxOpenConns(10)
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
		RedisConf: RedisConf{
			Address:  "127.0.0.1:6379",
			Password: "letjoy2023",
			DB:       1,
		},
		ServiceConf: ServiceConf{
			Hoopoe: "http://localhost:11011/hoopoe/v2/query",
			Smew:   "http://localhost:11012/smew/v2/query",
			Scream: "http://localhost:11014/scream/v2/query",
		},
		Secret: "youyue",
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
