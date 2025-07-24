package user

import (
	"github.com/nghianx1211/golang/graph"
	"github.com/nghianx1211/golang/graph/generated"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
)

func RegisterGraphQLRoutes(r *gin.Engine) {
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(
		generated.Config{Resolvers: &graph.Resolver{}},
	))

	r.POST("/graphql", gin.WrapH(srv))
	r.GET("/playground", gin.WrapH(playground.Handler("GraphQL", "/graphql")))
}
