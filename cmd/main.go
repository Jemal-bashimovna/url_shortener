package main

import (
	"log"
	"os"
	urlshortener "shotenedurl"
	"shotenedurl/pkg/handler"
	"shotenedurl/pkg/repository"
	"shotenedurl/pkg/service"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {

	if err := initConfig(); err != nil {
		logrus.Fatalf("Error initializing configs :%s", err)
	}

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading env variables: %s", err)
	}

	db, err := repository.NewPostgres(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		User:     viper.GetString("db.user"),
		DBName:   viper.GetString("db.dbname"),
		Password: os.Getenv("DB_PASSWORD"),
		SSLMode:  viper.GetString("db.sslmode"),
	})

	if err != nil {
		logrus.Fatalf("Error connecting to database: %s", err)
	}

	repo := repository.NewRepository(db)
	services := service.NewService(repo)
	handlers := handler.NewHandler(services)

	srv := new(urlshortener.Server)
	if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
		logrus.Fatalf("Error running server :%s", err)
	}

}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
