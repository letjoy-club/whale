package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"whale/graph"
	"whale/pkg/loader"
	"whale/pkg/restfulserver"
	"whale/pkg/whaleconf"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/letjoy-club/mida-tool/authenticator"
	"github.com/letjoy-club/mida-tool/dbutil"
	"github.com/letjoy-club/mida-tool/midacode"
	"github.com/letjoy-club/mida-tool/midacontext"
	"github.com/letjoy-club/mida-tool/redisutil"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func main() {

	port := flag.Int("port", 11013, "port to run server on")
	confPath := flag.String("conf", "conf.yaml", "path to config file")
	flag.Parse()

	conf := whaleconf.ReadConf(*confPath)
	db := conf.DB()
	redis := conf.Redis()
	loader := loader.NewLoader(db)

	gqlConf := graph.Config{Resolvers: &graph.Resolver{}}
	gqlConf.Directives.AdminOnly = func(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
		token := midacontext.GetClientToken(ctx)
		if token.IsAdmin() {
			return next(ctx)
		}
		return nil, midacode.ErrNotPermitted
	}

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(gqlConf))
	srv.AroundOperations(func(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
		oc := graphql.GetOperationContext(ctx)
		fmt.Printf("RawQuery: %s\n", oc.RawQuery)
		return next(ctx)
		// return graphql.ResponseHandler(func(ctx context.Context) *graphql.Response {
		// 	return &graphql.Response{
		// 		Errors: gqlerror.List{&gqlerror.Error{Message: "DEEP_QUERY"}},
		// 	}
		// })
	})
	srv.SetErrorPresenter(func(ctx context.Context, e error) *gqlerror.Error {
		var ok bool
		var midaErr midacode.Error2
		tmpErr := e

		for {
			tmpErr = errors.Unwrap(tmpErr)
			if tmpErr == nil {
				break
			}
			midaErr, ok = tmpErr.(midacode.Error2)
			if ok {
				break
			}
		}
		err := graphql.DefaultErrorPresenter(ctx, e)
		err.Extensions = map[string]interface{}{"cn": midaErr.CN()}
		return err
	})

	r := chi.NewRouter()

	secret := []byte(conf.Secret)
	auth := authenticator.Authenticator{Key: secret}

	adminToken, _ := auth.SignID("1000")
	fmt.Println("Admin Token: ", adminToken)

	services := midacontext.Services{
		Hoopoe: midacontext.NewServices(conf.ServicesSetting.Hoopoe, "1000"),
		Smew:   midacontext.NewServices(conf.ServicesSetting.Smew, "1000"),
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

			ctx = dbutil.WithDB(ctx, db)
			ctx = redisutil.WithRedis(ctx, redis)
			ctx = midacontext.WithLoader(ctx, loader)
			ctx = midacontext.WithServices(ctx, services)
			ctx = midacontext.WithClientToken(ctx, token)
			ctx = midacontext.WithAuthenticator(ctx, auth)

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

	log.Printf("connect to http://localhost:%d/whale/v2 for GraphQL playground", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), r))
}

func queryDepth() {
}
