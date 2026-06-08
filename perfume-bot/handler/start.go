package handler

import (
	"context"
	"log"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (h *Handler) StartHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	kb := models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "Каталог", CallbackData: "catalog:0"},
				{Text: "Категории", CallbackData: "categories"},
			},
			{
				{Text: "Бренды", CallbackData: "brands:0"},
			},
			{
				{Text: "Корзина", CallbackData: "cart"},
			},
		},
	}
	textMsg := "<b>Добро пожаловать в парфюмерный бутик!</b> ✨\n\nЗдесь собраны лучшие ароматы: от признанной мировой классики до редкого селективного парфюма.\n\n<i>Выберите раздел в меню ниже, чтобы найти свой идеальный аромат:</i>"

	if update.Message != nil && update.Message.Text == "/start" {
		err := h.repo.CreateUserIfNotExists(ctx, update.Message.From.ID)
		if err != nil {
			log.Printf("[StartHandler] failed to save user (ID: %d): %v", update.Message.From.ID, err)
		}

		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      update.Message.Chat.ID,
			Text:        textMsg,
			ParseMode:   models.ParseModeHTML,
			ReplyMarkup: kb,
		})
		return
	}

	if update.CallbackQuery != nil && update.CallbackQuery.Data == "main_menu" {
		b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: update.CallbackQuery.ID,
		})
		if update.CallbackQuery.Message.Message == nil {
			return
		}
		if update.CallbackQuery.Message.Message.Text == "" {
			b.DeleteMessage(ctx, &bot.DeleteMessageParams{
				ChatID:    update.CallbackQuery.Message.Message.Chat.ID,
				MessageID: update.CallbackQuery.Message.Message.ID,
			})
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
				Text:        textMsg,
				ParseMode:   models.ParseModeHTML,
				ReplyMarkup: kb,
			})
		} else {
			b.EditMessageText(ctx, &bot.EditMessageTextParams{
				ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
				MessageID:   update.CallbackQuery.Message.Message.ID,
				Text:        textMsg,
				ParseMode:   models.ParseModeHTML,
				ReplyMarkup: kb,
			})
		}
		return
	}
}

func (h *Handler) DefaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.CallbackQuery != nil {
		log.Printf("[DefaultHandler] unhandled callback data: %s", update.CallbackQuery.Data)
		b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: update.CallbackQuery.ID,
		})
		if update.CallbackQuery.Data == "ignore" {
			return
		}

		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.CallbackQuery.From.ID,
			Text:   "Эта кнопка больше недействительна. Чтобы вызвать актуальное меню, нажмите /start",
		})
		return
	}

	if update.Message != nil {
		log.Printf("[DefaultHandler] unhandled message: '%s'", update.Message.Text)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Чтобы вызвать меню, нажмите /start",
		})
		return
	}
}
