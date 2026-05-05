package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"perfume-bot/api"
	minio_cl "perfume-bot/clients/minio"
	"perfume-bot/handler"
	"perfume-bot/repository"

	"github.com/go-telegram/bot"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
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

	mcl, err := minio.New(os.Getenv("MINIO_ENDPOINT"), &minio.Options{
		Creds:  credentials.NewStaticV4(os.Getenv("MINIO_ACCESS_KEY"), os.Getenv("MINIO_SECRET_KEY"), ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalf("Error minio new: %v", err)
	}

	fileClient := minio_cl.New(mcl, os.Getenv("MINIO_BUCKET"), os.Getenv("MINIO_ENDPOINT"))

	apiServer := api.New(repo, fileClient)
	go func() {
		log.Println("API сервер запущен на порту :8181")
		err := apiServer.Run()
		if err != nil {
			log.Fatal("Ошибка API сервера: %v", err)
		}
	}()

	opts := []bot.Option{
		bot.WithDefaultHandler(updatesHandler.DefaultHandler),
	}

	b, err := bot.New(os.Getenv("TG_BOT_TOKEN"), opts...)
	if err != nil {
		panic(err)
	}
	b.RegisterHandler(bot.HandlerTypeMessageText, "start", bot.MatchTypeCommand, updatesHandler.StartHandler)
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "main_menu", bot.MatchTypeExact, updatesHandler.StartHandler)
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "cart", bot.MatchTypeExact, updatesHandler.CartCallbackHandler)
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "cart_clear", bot.MatchTypeExact, updatesHandler.ClearCartCallbackHandler)
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "catalog:", bot.MatchTypePrefix, updatesHandler.CatalogCallbackHandler)
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "categories", bot.MatchTypeExact, updatesHandler.CategoriesCallbackHandler)
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "brands:", bot.MatchTypePrefix, updatesHandler.BrandsCallbackHandler)
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "cart_add:", bot.MatchTypePrefix, updatesHandler.AddToCartCallbackHandler)
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "cp:", bot.MatchTypePrefix, updatesHandler.CategoryProductsCallbackHandler)
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "bp:", bot.MatchTypePrefix, updatesHandler.BrandProductsCallbackHandler)
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "p:", bot.MatchTypePrefix, updatesHandler.ProductCardCallbackHandler)

	botUser, err := b.GetMe(ctx)
	if err != nil {
		log.Fatalf("Error bot get me: %v", err)
	}
	log.Printf("Bot is started- %s", botUser.Username)
	b.Start(ctx)
}
