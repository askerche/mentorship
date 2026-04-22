package handler

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

// PhotoInterceptorHandler ловит отправленные боту фотографии и выдает их File ID
func (h *Handler) PhotoInterceptorHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	// Базовая защита: если это не сообщение или в нем нет фото - игнорируем
	if update.Message == nil || len(update.Message.Photo) == 0 {
		return
	}

	// АРХИТЕКТУРНЫЙ НЮАНС:
	// Telegram автоматически сжимает отправленное фото и присылает массив из нескольких размеров.
	// Самое высокое качество ВСЕГДА находится в самом конце этого массива.
	bestResolutionPhoto := update.Message.Photo[len(update.Message.Photo)-1]

	fileID := bestResolutionPhoto.FileID

	// Формируем красивый ответ с тегом <code>, чтобы ты мог скопировать ID в один клик с телефона или ПК
	msgText := fmt.Sprintf(
		"📸 <b>Фотография успешно загружена на сервера Telegram!</b>\n\n"+
			"Уникальный File ID (нажмите, чтобы скопировать):\n<code>%s</code>\n\n"+
			"<i>Вставьте этот ID в базу данных для нужного товара.</i>",
		fileID,
	)

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      msgText,
		ParseMode: models.ParseModeHTML,
	})
}
