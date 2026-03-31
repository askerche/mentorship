package handler

import (
	"context"
	"fmt"
	"log"
	"perfume-bot/repository"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type Handler struct {
	repo *repository.Repository
}

func New(repo *repository.Repository) *Handler {
	return &Handler{
		repo: repo,
	}
}

func (h *Handler) All(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.CallbackQuery != nil {
		b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: update.CallbackQuery.ID,
		})
		callback := update.CallbackQuery.Data
		var responseText string
		switch callback {
		case "catalog":
			products, err := h.repo.GetAllProducts(ctx)
			if err != nil {
				log.Printf("error repo.GetAllProducts: %v", err)
				responseText = "Произошла ошибка при загрузке каталога."
			} else if len(products) == 0 {
				responseText = "Каталог пока пуст."
			} else {
				responseText = "Наш каталог:\n\n"
				for _, p := range products {
					responseText += fmt.Sprintf("%s от %s \nЦена: %d руб.\nОписание: %s\n\n", p.Title, p.Brand.Title, p.Price, p.Description)
				}
			}
		case "categories":
			categories, err := h.repo.GetAllCategories(ctx)
			if err != nil {
				log.Printf("error repo.GetAllCategories: %v", err)
				responseText = "Произошла ошибка при загрузке категорий"
			} else if len(categories) == 0 {
				responseText = "Категорий пока нет."
			} else {

			}
		case "brands":
			responseText = "BVLGARY, Roja, Tom Ford, HFC."
		case "cart":
			responseText = "Ваша корзина пуста.Данный раздел в разработке"
		default:
			return
		}
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.CallbackQuery.From.ID,
			Text:   responseText,
		})
		return
	}
	if update.Message.Text == "/start" && update.Message != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Привет\n\nТут ты можешь купить духи, парфюм, туалетную воду и тд. Пользуйся кнопками ниже",
			ReplyMarkup: models.InlineKeyboardMarkup{
				InlineKeyboard: [][]models.InlineKeyboardButton{
					{
						{Text: "Каталог", CallbackData: "catalog"},
						{Text: "Категории", CallbackData: "categories"},
					},
					{
						{Text: "Бренды", CallbackData: "brands"},
					},
					{
						{Text: "Корзина(недоступно)", CallbackData: "cart"},
					},
				},
			},
		})
		return
	}
}
