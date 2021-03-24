package main

import (
	"github.com/akwanmaroso/ddd-drugs/infrastructure/auth"
	"github.com/akwanmaroso/ddd-drugs/infrastructure/persistence"
	"github.com/akwanmaroso/ddd-drugs/interfaces"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("not env gotten")
	}
}

func main() {
	dbDriver := os.Getenv("DB_DRIVER")
	dbHost := os.Getenv("DB_HOST")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbUser := os.Getenv("DB_USER")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	// Redis details
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisPassword := os.Getenv("REDIS_PASSWORD")

	services, err := persistence.NewRepositories(dbDriver, dbUser, dbPassword, dbPort, dbHost, dbName)
	if err != nil {
		panic(err)
	}

	defer services.Close()
	services.Automigrate()

	redisService, err := auth.NewRedisDB(redisHost, redisPort, redisPassword)
	if err != nil {
		log.Fatal(err)
	}

	token := auth.NewToken()

	users := interfaces.NewUsers(services.User, token, redisService.Auth)
	drugs := interfaces.NewDrug(services.Drug, services.User, token, redisService.Auth)
	//	authenticate := interfaces.NewAuthenticate(services.User, redisService.Auth, token)

	r := gin.Default()

	// users route
	r.GET("/users", users.GetUsers)
	r.GET("/users/user_id", users.GetUser)

	// drugs route
	r.GET("/drugs", drugs.GetAllDrug)

	PORT := os.Getenv("PORT") // implement in heroku
	if PORT == "" {
		PORT = "5000"
	}
	log.Fatal(r.Run(":" + PORT))
}
