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
	resp := formatPrayerTimesResp(prayers, hijraDate, city)

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      resp,
		ParseMode: models.ParseModeHTML,
	})
}

func formatPrayerTimesResp(prayers appModels.PrayerTimes, hijraDate appModels.Hijri, city string) string {
	prayerstime := []appModels.PrayerItem{
		{Name: "Фаджр", Time: prayers.Fajr},
		{Name: "Восход", Time: prayers.Sunrise},
		{Name: "Зухр", Time: prayers.Dhuhr},
		{Name: "Аср", Time: prayers.Asr},
		{Name: "Магриб", Time: prayers.Maghrib},
		{Name: "Иша", Time: prayers.Isha},
	}
	today := time.Now().Format("02.01.2006")
	now := time.Now()

	nextPrayer, currentPrayer, timeRemaining := definePrayerStatus(prayerstime, now)
	TimeToNextPrayer := formatTimeToNextPrayer(timeRemaining)

	var prayersResp string
	for _, k := range prayerstime {
		prayer := fmt.Sprintf("%s: %s", k.Name, k.Time)
		switch k.Name {
		case currentPrayer:
			prayer = fmt.Sprintf("<b>%s</b>\n", prayer)
		case nextPrayer:
			prayer = fmt.Sprintf("<code>%s</code> %s\n", prayer, TimeToNextPrayer)
		default:
			prayer = fmt.Sprintf("<code>%s</code>\n", prayer)
		}
		prayersResp += prayer
	}
	timings := fmt.Sprintf("Расписание на %s\n\n"+"%s", today, prayersResp)
	resp := fmt.Sprintf("Город: <b>%s</b>\n%s\n%s %s %s",
		city, timings, hijraDate.Day, hijraDate.Month.En, hijraDate.Year)
	return resp
}

func definePrayerStatus(prayerstime []appModels.PrayerItem, now time.Time) (string, string, time.Duration) {
	var nextPrayer string
	var currentPrayer string
	var timeRemaining time.Duration
	for _, k := range prayerstime {
		parsedTime, _ := time.Parse("15:04", k.Time)
		prayerTime := time.Date(now.Year(), now.Month(), now.Day(), parsedTime.Hour(), parsedTime.Minute(), 0, 0, now.Location())
		if now.After(prayerTime) {
			currentPrayer = k.Name
		} else {
			nextPrayer = k.Name
			timeRemaining = time.Until(prayerTime)
			break
		}
	}
	return nextPrayer, currentPrayer, timeRemaining
}

func formatTimeToNextPrayer(timeRemaining time.Duration) string {
	totalSeconds := int(timeRemaining.Seconds())
	hours := totalSeconds / 3600
	minutes := (totalSeconds % 3600) / 60
	TimeToNextPrayer := fmt.Sprintf("(через %dч %dм)", hours, minutes)
	return TimeToNextPrayer
}
