package handler

import (
	"context"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (h *Handler) CategoriesCallbackHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.CallbackQuery != nil {
		b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: update.CallbackQuery.ID,
		})
		categories, err := h.repo.GetAllCategories(ctx)
		if err != nil {
			log.Printf("[CategoriesCallbackHandler] failed to get categories: %v", err)
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.CallbackQuery.From.ID,
				Text:   "Произошла ошибка при загрузке категорий",
			})
			return
		}
		if len(categories) == 0 {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.CallbackQuery.From.ID,
				Text:   "Категорий пока нет.",
			})
			return
		}
		var kb models.InlineKeyboardMarkup
		kb.InlineKeyboard = make([][]models.InlineKeyboardButton, 0)
		for _, cat := range categories {
			var row []models.InlineKeyboardButton
			row = append(row, models.InlineKeyboardButton{
				Text:         cat.Title,
				CallbackData: fmt.Sprintf("cp:%d:0", cat.ID),
			})

			kb.InlineKeyboard = append(kb.InlineKeyboard, row)
		}
		kb.InlineKeyboard = append(kb.InlineKeyboard, []models.InlineKeyboardButton{
			{
				Text:         "🏠 Главное меню",
				CallbackData: "main_menu",
			},
		})
		b.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
			MessageID:   update.CallbackQuery.Message.Message.ID,
			Text:        "Выберите категорию:",
			ReplyMarkup: kb,
		})
	}
}

func (h *Handler) CategoryProductsCallbackHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.CallbackQuery != nil {
		b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: update.CallbackQuery.ID,
		})
		data := update.CallbackQuery.Data
		parts := strings.Split(data, ":")
		catID, _ := strconv.Atoi(parts[1])
		page, _ := strconv.Atoi(parts[2])

		limit := 5
		offset := page * limit

		totalCount, err := h.repo.GetProductsCountByCategoryID(ctx, catID)
		if err != nil {
			log.Printf("[CategoryProductsCallbackHandler] failed to get count: %v", err)
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.CallbackQuery.From.ID,
				Text:   "Произошла ошибка базы данных",
			})
			return
		}

		if totalCount == 0 {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.CallbackQuery.From.ID,
				Text:   "В этой категории пока нет парфюма.",
			})
			return
		}

		products, err := h.repo.GetProductsPageByCategoryID(ctx, catID, limit, offset)
		if err != nil {
			log.Printf("[CategoryProductsCallbackHandler] failed to get products: %v", err)
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.CallbackQuery.From.ID,
				Text:   "Ошибка при загрузке списка товаров.",
			})
			return
		}

		var kb models.InlineKeyboardMarkup
		for _, p := range products {
			kb.InlineKeyboard = append(kb.InlineKeyboard,
				[]models.InlineKeyboardButton{
					{
						Text:         fmt.Sprintf("%s %s - %d руб.", p.Brand.Title, p.Title, p.Price),
						CallbackData: fmt.Sprintf("p:%d|cp:%d:%d", p.ID, catID, page),
					},
				})
		}
		totalPages := int(math.Ceil(float64(totalCount) / float64(limit)))

		var navRow []models.InlineKeyboardButton
		if page > 0 {
			navRow = append(navRow, models.InlineKeyboardButton{
				Text:         "⬅️ Назад",
				CallbackData: fmt.Sprintf("cp:%d:%d", catID, page-1),
			})
		}
		navRow = append(navRow, models.InlineKeyboardButton{
			Text:         fmt.Sprintf("Стр. %d из %d", page+1, totalPages),
			CallbackData: "ignore",
		})

		if page < totalPages-1 {
			navRow = append(navRow, models.InlineKeyboardButton{
				Text:         "Далее ➡️",
				CallbackData: fmt.Sprintf("cp:%d:%d", catID, page+1),
			})
		}

		kb.InlineKeyboard = append(kb.InlineKeyboard, navRow)
		kb.InlineKeyboard = append(kb.InlineKeyboard, []models.InlineKeyboardButton{
			{
				Text:         "Вернутья к категориям",
				CallbackData: "categories",
			},
		})

		if update.CallbackQuery.Message.Message.Text == "" {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
				Text:        "Выберите парфюм из списка:",
				ReplyMarkup: kb,
			})
		} else {
			b.EditMessageText(ctx, &bot.EditMessageTextParams{
				ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
				MessageID:   update.CallbackQuery.Message.Message.ID,
				Text:        "Выберите парфюм из списка:",
				ReplyMarkup: kb,
			})
		}

	}
}
