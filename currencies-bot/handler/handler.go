package handler

import (
	"context"
	"currencies/client/fxratesapi"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func New(FxRatesApiClient *fxratesapi.FxRatesApiClient) *Handler {
	return &Handler{
		FxRatesApiClient: FxRatesApiClient,
	}
}

type Handler struct {
	FxRatesApiClient *fxratesapi.FxRatesApiClient
}

const CurrencyRub = "RUB"
const CurrencyUsd = "USD"
const CurrencyTry = "TRY"
const CurrencySar = "SAR"
const CurrencyEur = "EUR"

func (h *Handler) Handle(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message.Text == "/start" {
		url := "https://api.fxratesapi.com/latest?base=RUB&currencies=USD,EUR,TRY,SAR&resolution=1m&amount=1&places=6&format=json"
		res, err := http.Get(url)
		if err != nil {
			fmt.Println("error GET request to API: ", err)
			return
		}
		defer res.Body.Close()

		var rates fxratesapi.LatestResponse
		err = json.NewDecoder(res.Body).Decode(&rates)
		if err != nil {
			fmt.Println("error unmarshall: ", err)
			return
		}
		ratesData := fmt.Sprintf(
			"USD: %.2f\n"+
				"EUR: %.2f\n"+
				"TRY: %.2f\n"+
				"SAR: %.2f",
			1/rates.Rates.Usd,
			1/rates.Rates.Eur,
			1/rates.Rates.Try,
			1/rates.Rates.Sar,
		)

		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text: fmt.Sprintf(
				"Привет, Я - бот конвертер валют.\n"+
					"Актуальный курс валют на сегодня:\n\n<code>%s</code>\n",
				ratesData),
			ParseMode: models.ParseModeHTML,
		})
		return
	}

	h.HandleCurrencyConvert(ctx, b, update)
}

func (h *Handler) HandleCurrencyConvert(ctx context.Context, b *bot.Bot, update *models.Update) {
	currencyFrom := extractCurrency((update.Message.Text))

	if currencyFrom == "" {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Поддерживаемые валюты: RUB, USD, TRY, EUR, SAR",
		})
		return
	}

	re := regexp.MustCompile(`[\d\s]+`)
	valRaw := re.FindString(update.Message.Text)
	valString := strings.ReplaceAll(valRaw, " ", "")
	val, err := strconv.Atoi(valString)
	if err != nil {
		fmt.Println("error strconv: ", err)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Ошибка!",
		})
		return
	}
	rates, err := h.FxRatesApiClient.GetCurrencyRate(ctx, currencyFrom)
	if err != nil {
		fmt.Println("error FxRatesApiClient.GetCurrencyRate: ", err)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Ошибка!",
		})
		return
	}
	var resp string = "<code>"
	var supportedCurrencies = []string{CurrencyUsd, CurrencyRub, CurrencyEur, CurrencyTry, CurrencySar}
	for _, cur := range supportedCurrencies {
		if currencyFrom != cur {
			resp += fmt.Sprintf("%s: %.2f\n", cur, rates[cur]*float64(val))
		}
	}
	resp += "</code>"
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      resp,
		ParseMode: models.ParseModeHTML,
	})
}

func extractCurrency(s string) string {
	if strings.Contains(strings.ToLower(s), "руб") ||
		strings.Contains(s, "р") ||
		strings.Contains(s, "₽") {
		return CurrencyRub
	} else if strings.Contains(strings.ToLower(s), "usd") ||
		strings.Contains(s, "$") {
		return CurrencyUsd
	} else if strings.Contains(strings.ToLower(s), "try") ||
		strings.Contains(s, "лир") ||
		strings.Contains(s, "₺") {
		return CurrencyTry
	} else if strings.Contains(strings.ToLower(s), "sar") ||
		strings.Contains(s, "риял") {
		return CurrencySar
	} else if strings.Contains(strings.ToLower(s), "eur") ||
		strings.Contains(s, "евро") ||
		strings.Contains(s, "€") {
		return CurrencyEur
	}
	return ""
}
