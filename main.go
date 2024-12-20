package main

import (
	"log"
	"os"

	"gin-auth-mongo/databases"
	"gin-auth-mongo/graph"
	"gin-auth-mongo/graph/resolvers"
	"gin-auth-mongo/middlewares"

	"gin-auth-mongo/routes"
	"gin-auth-mongo/utils"
	"gin-auth-mongo/utils/consts"
	"gin-auth-mongo/utils/cron"
	"gin-auth-mongo/utils/jwkmanager"
	"gin-auth-mongo/utils/mail"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// load env
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	utils.InitLoggerDir()

	// init databases
	databases.InitRedis()
	databases.InitMinio()
	databases.InitMongoDB()
	err = databases.RunMigrations()
	if err != nil {
		log.Fatalf("Error running migrations: %v", err)
	}
	log.Println("Migrations completed")

	// init mail
	mail.InitMail()

	// init jwk manager
	err = jwkmanager.LoadSigningKeys(consts.PUBLIC_KEYS_FILE, consts.PRIVATE_KEYS_FILE)
	if err != nil {
		log.Printf("Error loading jwk keys: %v", err)
		log.Println("Trying to create new jwk keys...")

		// create new keys
		err := jwkmanager.UpdateKeys()
		if err != nil {
			log.Fatalf("Error creating new jwk keys: %v", err)
			panic(err)
		}
	}

	// init gin
	r := gin.Default()
	r.MaxMultipartMemory = 500 << 20 // 500MB

	cron.StartCron()

	// init routes
	routes.SetupRoutes(r)

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &resolvers.Resolver{}}))

	// support multipartform
	srv.AddTransport(transport.MultipartForm{
		MaxUploadSize: 500 << 20, // 500MB
	})

	// GraphQL routes
	r.GET("/graphql", gin.WrapH(playground.Handler("GraphQL playground", "/graphql/query")))
	graphqlGroup := r.Group("/graphql")
	graphqlGroup.Use(middlewares.CORSMiddleware(), middlewares.FlowLimitMiddleware(), middlewares.GraphQLMiddleware())
	{
		graphqlGroup.POST("/query", gin.WrapH(srv))
	}

	// run server
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	r.Run(host + ":" + port)
}
