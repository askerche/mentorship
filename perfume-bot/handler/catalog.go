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

func (h *Handler) CatalogCallbackHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.CallbackQuery != nil {
		b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: update.CallbackQuery.ID,
		})

		data := update.CallbackQuery.Data
		parts := strings.Split(data, ":")
		page, _ := strconv.Atoi(parts[1])

		limit := 5
		offset := page * limit

		totalCount, err := h.repo.GetAllProductsCounts(ctx)
		if err != nil {
			log.Printf("[CatalogCallbackHandler] failed to get count: %v", err)
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.CallbackQuery.From.ID,
				Text:   "Произошла ошибка базы данных",
			})
			return
		}

		if totalCount == 0 {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.CallbackQuery.From.ID,
				Text:   "В данный момент каталог пустой.",
			})
			return
		}

		totalPages := int(math.Ceil(float64(totalCount) / float64(limit)))

		products, err := h.repo.GetAllProductsPage(ctx, limit, offset)
		if err != nil {
			log.Printf("[CatalogCallbackHandler] failed to get products for catalog: %v", err)
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.CallbackQuery.From.ID,
				Text:   "Произошла ошибка базы данных",
			})
			return
		}

		var kb models.InlineKeyboardMarkup

		for _, p := range products {
			kb.InlineKeyboard = append(kb.InlineKeyboard, []models.InlineKeyboardButton{
				{
					Text:         fmt.Sprintf("🔸 %s %s - %d руб.", p.Title, p.Brand.Title, p.Price),
					CallbackData: fmt.Sprintf("p:%d|catalog:%d", p.ID, page),
				},
			})
		}

		var navRow []models.InlineKeyboardButton
		if page > 0 {
			navRow = append(navRow, models.InlineKeyboardButton{
				Text:         "« ⬅️ Назад »",
				CallbackData: fmt.Sprintf("catalog:%d", page-1),
			})
		}

		navRow = append(navRow, models.InlineKeyboardButton{
			Text:         fmt.Sprintf("Стр. %d из %d", page+1, totalPages),
			CallbackData: "ignore",
		})

		if page < totalPages-1 {
			navRow = append(navRow, models.InlineKeyboardButton{
				Text:         "« Далее ➡️ »",
				CallbackData: fmt.Sprintf("catalog:%d", page+1),
			})
		}

		kb.InlineKeyboard = append(kb.InlineKeyboard, navRow)
		kb.InlineKeyboard = append(kb.InlineKeyboard, []models.InlineKeyboardButton{
			{
				Text:         "« 🏠 Главное меню »",
				CallbackData: "main_menu",
			},
		})

		msgText := "Наш каталог:\n\n"

		if update.CallbackQuery.Message.Message.Text == "" {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
				Text:        msgText,
				ReplyMarkup: kb,
			})
		} else {
			b.EditMessageText(ctx, &bot.EditMessageTextParams{
				ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
				MessageID:   update.CallbackQuery.Message.Message.ID,
				Text:        msgText,
				ReplyMarkup: kb,
			})
		}
	}
}
