package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"perfume-bot/handler"
	"perfume-bot/repository"

	"github.com/go-telegram/bot"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

// Send any text message to the bot after the bot has been started

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	pool, err := pgxpool.New(ctx, os.Getenv("PG_DSN"))
	if err != nil {
		log.Fatalf("Error pgxpool new: %v", err)
	}
	err = pool.Ping(ctx)
	if err != nil {
		log.Fatalf("Error postgres ping: %v", err)
	}
	repo := repository.New(pool)
	updatesHandler := handler.New(repo)

	opts := []bot.Option{
		bot.WithDefaultHandler(updatesHandler.All),
	}

	b, err := bot.New(os.Getenv("TG_BOT_TOKEN"), opts...)
	if err != nil {
		panic(err)
	}
	botUser, err := b.GetMe(ctx)
	if err != nil {
		log.Fatalf("Error bot get me: %v", err)
	}
	log.Printf("Bot is started- %s", botUser.Username)
	b.Start(ctx)
}
