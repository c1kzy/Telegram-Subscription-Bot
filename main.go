package main

import (
	"context"
	"fmt"
	"net/http"
	"subscriptionbot/db"
	"subscriptionbot/service"
	"subscriptionbot/utilities"
	weatherAPI "subscriptionbot/weather"
	"time"

	tgapi "github.com/c1kzy/Telegram-API"
	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
	"github.com/phuslu/log"
)

func main() {
	ctx := context.Background()

	if envErr := godotenv.Load(); envErr != nil {
		log.Fatal().Err(fmt.Errorf("error loading .env file: %w", envErr))
	}
	cfg := &tgapi.Config{}

	log.DefaultLogger = log.Logger{
		Level:      log.InfoLevel,
		Caller:     cfg.Caller,
		TimeField:  cfg.TimeField,
		TimeFormat: time.RFC850,
		Writer:     &log.ConsoleWriter{},
	}

	if err := env.Parse(cfg); err != nil {
		log.Error().Err(err)
	}

	api := tgapi.GetAPI(cfg)
	database := db.GetDB()
	weather := weatherAPI.GetWeatherAPI()

	tgService := service.NewService(database, weather, api)

	go func() {
		tgService.Notify(ctx)
	}()

	api.RegisterCommand("/start", utilities.StartResponse)
	api.RegisterInput(tgService.AddSubscription)
	http.HandleFunc("/telegram", api.TelegramHandler)

	go func() {
		err := http.ListenAndServe(fmt.Sprintf(":%v", cfg.Port), nil)
		if err != nil {
			log.Fatal().Err(err)
		}

	}()

	log.Info().Msg("Server started")

}
