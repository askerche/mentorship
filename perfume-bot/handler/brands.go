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

func (h *Handler) BrandsCallbackHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.CallbackQuery != nil {
		b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: update.CallbackQuery.ID,
		})

		data := update.CallbackQuery.Data
		parts := strings.Split(data, ":")
		page, err := strconv.Atoi(parts[1])
		if err != nil {
			log.Printf("[BrandsCallbackHandler] invalid page number: %v", err)
			return
		}

		limit := 5
		offset := page * limit

		totalCount, err := h.repo.GetBrandsCount(ctx)
		if err != nil {
			log.Printf("[BrandsCallbackHandler] failed to get count: %v", err)
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.CallbackQuery.From.ID,
				Text:   "Произошла ошибка базы данных",
			})
			return
		}

		if totalCount == 0 {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.CallbackQuery.From.ID,
				Text:   "Бренды пока не добавлены.",
			})
			return
		}

		brands, err := h.repo.GetBrandsPage(ctx, limit, offset)
		if err != nil {
			log.Printf("[BrandsCallbackHandler] failed to get brands: %v", err)
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.CallbackQuery.From.ID,
				Text:   "Произошла ошибка при загрузке брендов",
			})
			return
		}
		var kb models.InlineKeyboardMarkup
		kb.InlineKeyboard = make([][]models.InlineKeyboardButton, 0)
		for _, br := range brands {
			kb.InlineKeyboard = append(kb.InlineKeyboard, []models.InlineKeyboardButton{
				{
					Text:         br.Title,
					CallbackData: fmt.Sprintf("bp:%d:0", br.ID),
				},
			})
		}

		totalPages := int(math.Ceil(float64(totalCount) / float64(limit)))

		var navRow []models.InlineKeyboardButton
		if page > 0 {
			navRow = append(navRow, models.InlineKeyboardButton{
				Text:         "« ⬅️ Назад »",
				CallbackData: fmt.Sprintf("brands:%d", page-1),
			})
		}

		navRow = append(navRow, models.InlineKeyboardButton{
			Text:         fmt.Sprintf("Стр. %d из %d", page+1, totalPages),
			CallbackData: "brands:ignore",
		})

		if page < totalPages-1 {
			navRow = append(navRow, models.InlineKeyboardButton{
				Text:         "« Далее ➡️ »",
				CallbackData: fmt.Sprintf("brands:%d", page+1),
			})
		}

		kb.InlineKeyboard = append(kb.InlineKeyboard, navRow)
		kb.InlineKeyboard = append(kb.InlineKeyboard, []models.InlineKeyboardButton{
			{
				Text:         "« 🏠 Главное меню »",
				CallbackData: "main_menu",
			},
		})

		b.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
			MessageID:   update.CallbackQuery.Message.Message.ID,
			Text:        "Выберите бренд:",
			ReplyMarkup: kb,
		})
	}
}

func (h *Handler) BrandProductsCallbackHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.CallbackQuery != nil {
		b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: update.CallbackQuery.ID,
		})
		data := update.CallbackQuery.Data
		parts := strings.Split(data, ":")
		if len(parts) != 3 {
			log.Printf("[BrandProducts] invalid callback data format: %s", data)
			return
		}
		brandID, _ := strconv.Atoi(parts[1])
		page, _ := strconv.Atoi(parts[2])

		limit := 5
		offset := page * limit

		totalCount, err := h.repo.GetProductsCountByBrandID(ctx, brandID)
		if err != nil {
			log.Printf("[BrandProducts] failed to get count: %v", err)
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.CallbackQuery.From.ID,
				Text:   "Произошла ошибка базы данных",
			})
			return
		}
		if totalCount == 0 {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.CallbackQuery.From.ID,
				Text:   "У этого бренда пока нет парфюма в наличии.",
			})
			return
		}
		products, err := h.repo.GetProductsPageByBrandID(ctx, brandID, limit, offset)
		if err != nil {
			log.Printf("[BrandProductsCallbackHandler] failed to get products for brand: %v", err)
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.CallbackQuery.From.ID,
				Text:   "Произошла ошибка при загрузке парфюма этого бренда.",
			})
			return
		}

		var kb models.InlineKeyboardMarkup

		for _, p := range products {
			kb.InlineKeyboard = append(kb.InlineKeyboard, []models.InlineKeyboardButton{
				{
					Text:         fmt.Sprintf("🔸 %s %s - %d руб.", p.Title, p.Brand.Title, p.Price),
					CallbackData: fmt.Sprintf("p:%d|bp:%d:%d", p.ID, brandID, page),
				},
			})
		}

		totalPages := int(math.Ceil(float64(totalCount) / float64(limit)))

		var navRow []models.InlineKeyboardButton
		if page > 0 {
			navRow = append(navRow, models.InlineKeyboardButton{
				Text:         "« ⬅️ Назад »",
				CallbackData: fmt.Sprintf("bp:%d:%d", brandID, page-1),
			})
		}

		navRow = append(navRow, models.InlineKeyboardButton{
			Text:         fmt.Sprintf("Стр. %d из %d", page+1, totalPages),
			CallbackData: "ignore",
		})

		if page < totalPages-1 {
			navRow = append(navRow, models.InlineKeyboardButton{
				Text:         "« Далее ➡️ »",
				CallbackData: fmt.Sprintf("bp:%d:%d", brandID, page+1),
			})
		}

		kb.InlineKeyboard = append(kb.InlineKeyboard, navRow)
		kb.InlineKeyboard = append(kb.InlineKeyboard, []models.InlineKeyboardButton{
			{
				Text:         "« 🔙 К списку брендов »",
				CallbackData: "brands:0",
			},
		})
		msg := update.CallbackQuery.Message.Message
		textMsg := "Выберите парфюм из списка:"

		if msg.Text == "" {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:      msg.Chat.ID,
				Text:        textMsg,
				ReplyMarkup: kb,
			})
		} else {
			b.EditMessageText(ctx, &bot.EditMessageTextParams{
				ChatID:      msg.Chat.ID,
				MessageID:   msg.ID,
				Text:        textMsg,
				ReplyMarkup: kb,
			})
		}
	}
}
