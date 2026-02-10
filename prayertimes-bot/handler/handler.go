package handler

import (
	"context"
	"fmt"
	"prayertimes/client/aladhan"
	appModels "prayertimes/models"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type Handler struct {
}

func Handle(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message.Text == "/start" {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Просто напиши мне название своего города и я пришлю расписание молитв на сегодня.",
		})
		return
	}
	city := update.Message.Text
	prayers, err := aladhan.GetTodayPrayerTimesByCity(ctx, city)
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Не удалось найти этот город",
		})
		return
	}
	resp := FormatPrayerTimesResp(prayers, city)

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      resp,
		ParseMode: models.ParseModeHTML,
	})
}

func FormatPrayerTimesResp(prayers appModels.PrayerTimes, city string) string {
	today := time.Now().Format("02.01.2006")
	timings := fmt.Sprintf(
		"Расписание на %s\n\n"+
			"<code>Фаджр: %s\n"+
			"Восход: %s\n"+
			"Зухр: %s\n"+
			"Аср: %s\n"+
			"Магриб: %s\n"+
			"Иша: %s</code>",
		today,
		prayers.Fajr,
		prayers.Sunrise,
		prayers.Dhuhr,
		prayers.Asr,
		prayers.Maghrib,
		prayers.Isha,
	)
	resp := fmt.Sprintf("Город: <b>%s</b>\n%s", city, timings)
	return resp
}
