package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"prayertimes/client/aladhan"
	"prayertimes/handler"

	"github.com/go-telegram/bot"
	"github.com/joho/godotenv"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	aladhanClient := aladhan.New()
	botHandler := handler.New(aladhanClient)

	opts := []bot.Option{
		bot.WithDefaultHandler(botHandler.Handle),
	}
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	b, err := bot.New(os.Getenv("TG_BOT_TOKEN"), opts...)
	if err != nil {
		panic(err)
	}
	b.Start(ctx)
}
