package handler

import (
	"context"
	"fmt"
	"log"
	"prayertimes/client/aladhan"
	appModels "prayertimes/models"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type Handler struct {
	aladhanClient *aladhan.Client
}

func New(aladhanClient *aladhan.Client) *Handler {
	return &Handler{
		aladhanClient: aladhanClient,
	}
}

func (h *Handler) Handle(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message.Text == "/start" {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Просто напиши мне название своего города и я пришлю расписание молитв на сегодня.",
		})
		return
	}
	city := update.Message.Text
	prayers, hijraDate, err := h.aladhanClient.GetTodayPrayerTimesByCity(ctx, city)
	if err != nil {
		log.Println(err)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Не удалось найти этот город",
		})
		return
	}
	resp := FormatPrayerTimesResp(prayers, hijraDate, city)

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      resp,
		ParseMode: models.ParseModeHTML,
	})
}

func FormatPrayerTimesResp(prayers appModels.PrayerTimes, hijraDate appModels.Hijri, city string) string {
	type Prayer struct {
		Name string
		Time string
	}
	prayerstime := []Prayer{
		{"Фаджр", prayers.Fajr},
		{"Восход", prayers.Sunrise},
		{"Зухр", prayers.Dhuhr},
		{"Аср", prayers.Asr},
		{"Магриб", prayers.Maghrib},
		{"Иша", prayers.Isha},
	}

	var nextPrayer string
	today := time.Now().Format("02.01.2006")
	now := time.Now()
	for _, k := range prayerstime {
		parsedTime, _ := time.Parse("15:04", k.Time)
		prayerTime := time.Date(now.Year(), now.Month(), now.Day(), parsedTime.Hour(), parsedTime.Minute(), 0, 0, now.Location())
		if now.Before(prayerTime) {
			nextPrayer = k.Name
			break
		}
	}
	var prayersResp string
	for _, k := range prayerstime {
		prayer := fmt.Sprintf("%s: %s\n", k.Name, k.Time)
		if k.Name == nextPrayer {
			prayer = fmt.Sprintf("<b>%s</b>", prayer)
		} else {
			prayer = fmt.Sprintf("<code>%s</code>", prayer)
		}
		prayersResp += prayer
	}
	timings := fmt.Sprintf(
		"Расписание на %s\n\n"+"%s",
		today,
		prayersResp,
	)
	resp := fmt.Sprintf("Город: <b>%s</b>\n%s\n%s %s %s", city, timings, hijraDate.Day, hijraDate.Month.En, hijraDate.Year)
	return resp
}
