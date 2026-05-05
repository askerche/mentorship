package handler

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (h *Handler) CartCallbackHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.CallbackQuery != nil {
		b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: update.CallbackQuery.ID,
		})

		telegramID := update.CallbackQuery.From.ID
		items, err := h.repo.GetCartByTelegramID(ctx, telegramID)
		if err != nil {
			log.Printf("[CartCallbackHandler] failed to get cart for user %d: %v", telegramID, err)
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: telegramID,
				Text:   "Произошла ошибка при загрузке корзины.",
			})
			return
		}
		var kb models.InlineKeyboardMarkup
		if len(items) == 0 {
			kb.InlineKeyboard = [][]models.InlineKeyboardButton{
				{
					{
						Text:         "🏠 Главное меню",
						CallbackData: "main_menu",
					},
				},
			}
			msgText := "<b>Ваша корзина пуста</b> \n\nПерейдите в каталог или категории, чтобы выбрать свой идеальный аромат."
			if update.CallbackQuery.Message.Message != nil {
				b.EditMessageText(ctx, &bot.EditMessageTextParams{
					ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
					MessageID:   update.CallbackQuery.Message.Message.ID,
					Text:        msgText,
					ParseMode:   models.ParseModeHTML,
					ReplyMarkup: kb,
				})
			}
			return
		}
		var receipt strings.Builder
		receipt.WriteString("🛒 <b>Ваш заказ:</b>\n\n")

		totalSum := 0

		for i, item := range items {
			itemTotal := item.Price * item.Quantity
			totalSum += itemTotal
			receipt.WriteString(fmt.Sprintf("%d. <b>%s %s</b>\n", i+1, item.BrandName, item.Title))
			receipt.WriteString(fmt.Sprintf("▪️ %d шт. х %d руб. = <b>%d руб.</b>\n\n", item.Quantity, item.Price, itemTotal))
		}
		receipt.WriteString(fmt.Sprintf("<b>ИТОГО: %d руб.</b>", totalSum))

		kb.InlineKeyboard = [][]models.InlineKeyboardButton{
			{
				{Text: "💳 Оформить заказ", CallbackData: "checkout"},
			},
			{
				{Text: "🗑 Очистить корзину", CallbackData: "cart_clear"},
			},
			{
				{Text: "🏠 Главное меню", CallbackData: "main_menu"},
			},
		}
		if update.CallbackQuery.Message.Message != nil {
			b.EditMessageText(ctx, &bot.EditMessageTextParams{
				ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
				MessageID:   update.CallbackQuery.Message.Message.ID,
				Text:        receipt.String(),
				ParseMode:   models.ParseModeHTML,
				ReplyMarkup: kb,
			})
		}
	}
}

func (h *Handler) AddToCartCallbackHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.CallbackQuery != nil {
		data := update.CallbackQuery.Data
		parts := strings.Split(data, ":")
		if len(parts) != 2 {
			log.Printf("[AddToCart] invalid data format: %s", data)
			return
		}
		productId, err := strconv.Atoi(parts[1])
		if err != nil {
			log.Printf("[AddToCart] invalid product ID: %v", err)
			return
		}
		telegramId := update.CallbackQuery.From.ID

		err = h.repo.AddToCart(ctx, telegramId, productId)
		if err != nil {
			log.Printf("[AddToCart] failed to add product %d for user %d: %v", productId, telegramId, err)
			b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
				CallbackQueryID: update.CallbackQuery.ID,
				Text:            "❌ Ошибка базы данных при добавлении товара.",
				ShowAlert:       true,
			})
			return
		}
		b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: update.CallbackQuery.ID,
			Text:            "✅ Товар успешно добавлен в корзину!",
			ShowAlert:       false,
		})
	}
}

func (h *Handler) ClearCartCallbackHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.CallbackQuery != nil {
		telegramID := update.CallbackQuery.From.ID

		err := h.repo.ClearCart(ctx, telegramID)
		if err != nil {
			log.Printf("[ClearCart] failed to clear cart for user %d: %v", telegramID, err)
			b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
				CallbackQueryID: update.CallbackQuery.ID,
				Text:            "❌ Ошибка, не удалость очистить корзину.",
				ShowAlert:       true,
			})
			return
		}
		b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: update.CallbackQuery.ID,
			Text:            "✅ Корзина очищена!",
			ShowAlert:       false,
		})

		var kb models.InlineKeyboardMarkup
		kb.InlineKeyboard = [][]models.InlineKeyboardButton{
			{
				{
					Text:         "🏠 Главное меню",
					CallbackData: "main_menu",
				},
			},
		}

		msgText := "<b>Ваша корзина пуста</b> \n\nПерейдите в каталог или категории, чтобы выбрать свой идеальный аромат."
		if update.CallbackQuery.Message.Message != nil {
			b.EditMessageText(ctx, &bot.EditMessageTextParams{
				ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
				MessageID:   update.CallbackQuery.Message.Message.ID,
				Text:        msgText,
				ParseMode:   models.ParseModeHTML,
				ReplyMarkup: kb,
			})
		}
	}
}
