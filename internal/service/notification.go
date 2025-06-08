package service

import (
	quote "cash_piggy_bank/internal/api/quotes"
	"cash_piggy_bank/internal/botpkg"
	"cash_piggy_bank/internal/repository/sqlite"
	"fmt"
	"time"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
)

func Notification(text string, bot *telego.Bot, ctx *th.Context, message telego.Message) {
	//Задачу вызывает горутина, когда приходит время
	task := func() {

		quote := quote.Quote()
		userIDs := sqlite.DbQuery("40d9f51c-810c-4023-ae92-c314d90aaf30", bot, ctx, message, nil)

		for _, val := range userIDs {
			go botpkg.SendText(bot, nil, message, fmt.Sprintf("Помни о своих целях 😉\n\n%s", quote), "нотификация", val)
		}

	}
	scheduleDailyAt(18, 03, task)
}

func scheduleDailyAt(hour, min int, task func()) {
	go func() {
		for {
			now := time.Now()
			next := time.Date(now.Year(), now.Month(), now.Day(), hour, min, 0, 0, time.Local)
			if now.After(next) {
				next = next.Add(24 * time.Hour)

			}

			timer := time.NewTimer(next.Sub(now))
			<-timer.C
			task()
		}
	}()
}
