package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (h *Handler) ProductCardCallbackHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.CallbackQuery != nil {
		b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: update.CallbackQuery.ID,
		})

		data := update.CallbackQuery.Data
		mainParts := strings.Split(data, "|")
		if len(mainParts) != 2 {
			log.Printf("[ProductCard] invalid composite data format: %s", data)
			return
		}
		productCommand := mainParts[0]
		returnCommand := mainParts[1]
		productParts := strings.Split(productCommand, ":")
		if len(productParts) != 2 {
			log.Printf("[ProductCard] invalid product command format: %s", productCommand)
			return
		}
		productID, err := strconv.Atoi(productParts[1])
		if err != nil {
			log.Printf("[ProductCard] invalid product ID: %v", err)
			return
		}

		product, err := h.repo.GetProductByID(ctx, productID)
		if err != nil {
			log.Printf("[ProductCard] failed to get product by id: %v", err)
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.CallbackQuery.From.ID,
				Text:   "Ошибка при загрузке товара.",
			})
			return
		}
		if product == nil {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.CallbackQuery.From.ID,
				Text:   "Товар не найден или был удален.",
			})
			return
		}

		var kb models.InlineKeyboardMarkup
		kb.InlineKeyboard = [][]models.InlineKeyboardButton{
			{
				{
					Text:         "🛒 Добавить в корзину",
					CallbackData: fmt.Sprintf("cart_add:%d", productID),
				},
			},
			{
				{
					Text:         "Вернуться назад",
					CallbackData: returnCommand,
				},
				{
					Text:         "Перейти в корзину",
					CallbackData: "cart",
				},
			},
			{
				{
					Text:         "🏠 Главное меню",
					CallbackData: "main_menu",
				},
			},
		}

		caption := fmt.Sprintf("✨ <b>%s %s</b>\n\n▪️ <b>Цена:</b> %d руб.\n\n▪️ <b>Ноты:</b>\n<i>%s</i>",
			product.Brand.Title, product.Title, product.Price, product.Description,
		)
		if product.ImageFileID != "" {
			resp, err := http.Get(product.ImageFileID)
			if err != nil {
				log.Printf("[product.ImageFileID] Ошибка скачивания из MinIO: %v", err)
				b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
					Text:        "🖼 <i>(Фото скоро появится)</i>\n\n" + caption,
					ParseMode:   models.ParseModeHTML,
					ReplyMarkup: kb,
				})
			} else {
				_, err := b.SendPhoto(ctx, &bot.SendPhotoParams{
					ChatID: update.CallbackQuery.Message.Message.Chat.ID,
					Photo: &models.InputFileUpload{
						Filename: "image.png",
						Data:     resp.Body,
					},
					Caption:     caption,
					ParseMode:   models.ParseModeHTML,
					ReplyMarkup: kb,
				})

				resp.Body.Close()

				if err != nil {
					log.Printf("[product.ImageFileID] Ошибка отправки фото: %v (ImageFileID: %s)", err, product.ImageFileID)
				}
			}
		} else {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
				Text:        "🖼 <i>(Фото скоро появится)</i>\n\n" + caption,
				ParseMode:   models.ParseModeHTML,
				ReplyMarkup: kb,
			})
		}
	}
}
