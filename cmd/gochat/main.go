package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/sazonovItas/gochat-tcp/cmd/gochat/app"
	"github.com/sazonovItas/gochat-tcp/internal/utils"
)

func main() {
	// get value of env variable
	configEnv := utils.GetEnv()

	// load env variables from file
	err := godotenv.Load("./config/.env." + configEnv)
	if err != nil {
		log.Fatalf("%s: %s", "error to load env variables from file", err.Error())
	}

	// init application configuration
	cfg, err := app.InitAppConfig(&app.Options{
		Env: configEnv,

		LogWriter: os.Stdout,
	})
	if err != nil || cfg == nil {
		log.Fatalf("%s: %s", "error to init app config", err.Error())
	}

	// init application
	app, err := app.InitApp(cfg)
	if err != nil {
		log.Fatalf("%s: %s", "error to init app", err.Error())
	}

	// Run application
	log.Fatal(app.Run())
}
