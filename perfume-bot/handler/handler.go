package handler

import (
	"context"
	"fmt"
	"log"
	"perfume-bot/repository"
	"strconv"
	"strings"

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

func (h *Handler) StartHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
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

func (h *Handler) CatalogCallbackHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.CallbackQuery != nil {
		b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: update.CallbackQuery.ID,
		})
		var responseText string
		products, err := h.repo.GetAllProducts(ctx)
		if err != nil {
			log.Printf("[CatalogCallbackHandler] failed to get products for catalog: %v", err)
			responseText = "Произошла ошибка при загрузке каталога."
			return
		} else if len(products) == 0 {
			responseText = "Каталог пока пуст."
			return
		}
		responseText = "Наш каталог:\n\n"
		for _, p := range products {
			responseText += fmt.Sprintf("%s от %s \nЦена: %d руб.\nОписание: %s\n\n", p.Title, p.Brand.Title, p.Price, p.Description)
		}
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.CallbackQuery.From.ID,
			Text:   responseText,
		})
	}
}

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
				CallbackData: fmt.Sprintf("category_id:%d", cat.ID),
			})

			kb.InlineKeyboard = append(kb.InlineKeyboard, row)
		}
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      update.CallbackQuery.From.ID,
			Text:        "Выберите категорию парфюма:",
			ReplyMarkup: kb,
		})
	}
}

func (h *Handler) CategoryProductsCallbackHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.CallbackQuery != nil {
		b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: update.CallbackQuery.ID,
		})
		categoryIdStr := strings.TrimPrefix(update.CallbackQuery.Data, "category_id:")
		categoryId, err := strconv.Atoi(categoryIdStr)
		if err != nil {
			log.Printf("[CategoryProductsCallbackHandler] failed to parse category id: %v", err)
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.CallbackQuery.From.ID,
				Text:   "Произошла техническая ошибка. Пожалуйста, попробуйте еще раз.",
			})
			return
		}
		products, err := h.repo.GetProductsByCategoryID(ctx, categoryId)

		var responseText string
		if err != nil {
			log.Printf("[CategoryProductsCallbackHandler] failed to get products for category id: %v", err)
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.CallbackQuery.From.ID,
				Text:   "Произошла ошибка при загрузке парфюма из этой категории.",
			})
			return
		}
		if len(products) == 0 {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.CallbackQuery.From.ID,
				Text:   "В данной категории нет парфюма.",
			})
			return
		}
		responseText = "В этой категории представлены:\n\n"
		for _, p := range products {
			responseText += fmt.Sprintf("%s от %s \nЦена: %d руб.\nОписание: %s\n\n", p.Title, p.Brand.Title, p.Price, p.Description)
		}
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.CallbackQuery.From.ID,
			Text:   responseText,
		})
	}
}

func (h *Handler) BrandsCallbackHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.CallbackQuery != nil {
		b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: update.CallbackQuery.ID,
		})
		brands, err := h.repo.GetAllBrands(ctx)
		if err != nil {
			log.Printf("[BrandsCallbackHandler]failed to get brands: %v", err)
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.CallbackQuery.From.ID,
				Text:   "Произошла ошибка при загрузке брендов",
			})
			return
		}
		if len(brands) == 0 {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.CallbackQuery.From.ID,
				Text:   "Брендов пока нет.",
			})
			return
		}
		var kb models.InlineKeyboardMarkup
		kb.InlineKeyboard = make([][]models.InlineKeyboardButton, 0)
		var row []models.InlineKeyboardButton
		for i, br := range brands {
			row = append(row, models.InlineKeyboardButton{
				Text:         br.Title,
				CallbackData: fmt.Sprintf("brand_id:%d", br.ID),
			})
			if (i+1)%2 == 0 || i == len(brands)-1 {
				kb.InlineKeyboard = append(kb.InlineKeyboard, row)
				row = nil
			}
		}
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      update.CallbackQuery.From.ID,
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
		brandStrId := strings.TrimPrefix(update.CallbackQuery.Data, "brand_id:")
		brandId, err := strconv.Atoi(brandStrId)
		if err != nil {
			log.Printf("[BrandProductsCallbackHandler] failed to parse brand id: %v", err)
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.CallbackQuery.From.ID,
				Text:   "Произошла ошибка при обработке кнопки. Попробуйте еще раз.",
			})
			return
		}
		products, err := h.repo.GetProductsByBrandID(ctx, brandId)
		var responseText string
		if err != nil {
			log.Printf("[BrandProductsCallbackHandler] failed to get products for brand: %v", err)
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.CallbackQuery.From.ID,
				Text:   "Произошла ошибка при загрузке парфюма этого бренда.",
			})
			return
		}
		if len(products) == 0 {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.CallbackQuery.From.ID,
				Text:   "У данного бренда парфюм не предоставлен.",
			})
			return
		}
		responseText = "Парфюм от выбранного бренда:\n\n"
		for _, p := range products {
			responseText += fmt.Sprintf("%s от %s \nЦена: %d руб.\nОписание: %s\n\n", p.Title, p.Brand.Title, p.Price, p.Description)
		}
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.CallbackQuery.From.ID,
			Text:   responseText,
		})
	}
}

func (h *Handler) DefaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.CallbackQuery != nil {
		log.Printf("[DefaultHandler] unhandled callback data: %s", update.CallbackQuery.Data)
		b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: update.CallbackQuery.ID,
		})

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
