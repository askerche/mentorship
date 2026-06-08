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
	"github.com/redis/go-redis/v9"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	redisAddr := os.Getenv("REDIS_ADDR")
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	defer rdb.Close()

	aladhanClient := aladhan.New(rdb)
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
