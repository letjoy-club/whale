package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"whale/graph"
	"whale/pkg/loader"
	"whale/pkg/mq"
	"whale/pkg/restfulserver"
	"whale/pkg/whaleconf"

	"github.com/letjoy-club/mida-tool/qcloudutil"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/letjoy-club/mida-tool/authenticator"
	"github.com/letjoy-club/mida-tool/dbutil"
	"github.com/letjoy-club/mida-tool/graphqlutil"
	"github.com/letjoy-club/mida-tool/midacontext"
	"github.com/letjoy-club/mida-tool/proxy"
	"github.com/letjoy-club/mida-tool/pulsarutil"
	"github.com/letjoy-club/mida-tool/qcloudutil/clsutil"
	"github.com/letjoy-club/mida-tool/redisutil"
	"github.com/letjoy-club/mida-tool/tracerutil"
)

func main() {
	port := flag.Int("port", 11013, "port to run server on")
	confPath := flag.String("conf", "conf.yaml", "path to config file")

	upstream := flag.String("upstream", "", "upstream url, example: ws://localhost:11013/whale/v1/ws?userIds=1000,u_abc")

	flag.Parse()

	whaleConf := whaleconf.ReadConf(*confPath)
	db := whaleConf.DB()
	redis := whaleConf.Redis()
	loader := loader.NewLoader(db)
	whaleConf.QCloud.Init()
	cls := whaleConf.QCloud.CLS.Client
	if cls != nil && whaleConf.QCloud.CLS.TopicID != "" {
		cls.Start()
	}

	publisher := whaleConf.MatchingPublisher()
	subscriber := whaleConf.CreateSubscriber()
	pullTopics(subscriber, db, redis, loader, whaleConf.QCloud)

	tp := whaleConf.Trace()
	defer func() { tp.Shutdown(context.Background()) }()
	tr := tp.Tracer("whale-graph")

	gqlConf := graph.Config{Resolvers: &graph.Resolver{}}
	gqlConf.Directives.AdminOnly = graphqlutil.AdminOnly

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(gqlConf))
	srv.AroundOperations(graphqlutil.AroundOperations(tr, "whale"))
	srv.SetErrorPresenter(graphqlutil.ErrorPresenter)

	r := chi.NewRouter()

	secret := []byte(whaleConf.Secret)
	auth := authenticator.Authenticator{Key: secret}

	services := midacontext.Services{
		Hoopoe: midacontext.NewServices(whaleConf.ServiceConf.Hoopoe, "1000"),
		Smew:   midacontext.NewServices(whaleConf.ServiceConf.Smew, "1000"),
		Scream: midacontext.NewServices(whaleConf.ServiceConf.Scream, "1000"),
	}

	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			token := midacontext.ParseToken(r, auth)

			if token.IsInvalid() {
				w.Header().Set("Content-Type", "application/json")
				encoder := json.NewEncoder(w)
				encoder.Encode(midacontext.GraphQLResp{
					Data:   struct{}{},
					Errors: []midacontext.GraphQLErr{{Message: string(token), Path: []string{}}},
				})
				return
			}

			if proxy.RedirectProxy(token.String(), w, r) {
				return
			}

			ctx = dbutil.WithDB(ctx, db)
			ctx = clsutil.WithGraphLogger(ctx, whaleConf.QCloud.CLS.Client, whaleConf.QCloud.CLS.TopicID, "whale")
			ctx = redisutil.WithRedis(ctx, redis)
			ctx = pulsarutil.WithMQ(ctx, publisher)
			ctx = midacontext.WithLoader(ctx, loader)
			ctx = midacontext.WithServices(ctx, services)
			ctx = midacontext.WithClientToken(ctx, token)
			ctx = midacontext.WithAuthenticator(ctx, auth)
			ctx = tracerutil.WithSpanContext(ctx, r.Header)

			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	})

	r.Use(middleware.Logger)

	r.Route("/whale/", func(r chi.Router) {
		r.Route("/v2", func(r chi.Router) {
			r.Handle("/query", srv)
			r.Handle("/", playground.Handler("GraphQL playground", "/whale/v2/query"))
		})
		r.Route("/v1", func(r chi.Router) {
			restfulserver.Mount(r)
		})
	})

	if *upstream != "" {
		go proxy.Gateway(*upstream, r)
	}

	log.Printf("connect to http://localhost:%d/whale/v2 for GraphQL playground", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), r))
}

func pullTopics(subscriber *whaleconf.Subscriber, db *gorm.DB, redis *redis.Client, loader *loader.Loader, qCloud qcloudutil.QCloudConf) {
	ctx := context.Background()
	ctx = qcloudutil.WithQCloud(ctx, qCloud)
	ctx = redisutil.WithRedis(ctx, redis)
	ctx = midacontext.WithLoader(ctx, loader)
	ctx = dbutil.WithDB(ctx, db)
	ctx = pulsarutil.WithMQ(ctx, subscriber)
	go mq.UserLevelChangeListener(ctx)
}
