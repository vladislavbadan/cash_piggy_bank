package service

import (
	"cash_piggy_bank/internal/botpkg"
	"cash_piggy_bank/internal/domain"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

// функция рассчтеа остатка денег
func Calculation(text string, bot *telego.Bot, ctx *th.Context, message telego.Message, UserMap map[int64]*domain.User) {

	if UserMap[message.Chat.ID].Command == "Ждем сумму" {
		//
		text = strings.ReplaceAll(text, " ", "")
		textInt, err := strconv.Atoi(text)
		if err != nil {
			_, _ = bot.SendMessage(ctx,
				tu.Message(
					tu.ID(message.Chat.ID),
					"Кажется кажется ты ввел не число. Попробуй еще раз",
				),
			)
			return
		}

		UserMap[message.Chat.ID].Money = textInt
		UserMap[message.Chat.ID].Command = "Считаем сколько осталось"

		// Отправляем клавиатуру с датами
		botpkg.SendKeyboard(bot, ctx, message, nil, "calculationDate")

	} else if UserMap[message.Chat.ID].Command == "Считаем сколько осталось" {

		time1, err := time.Parse("02.01.2006", text)
		if err != nil {
			botpkg.SendText(bot, ctx, message, "Кажется ты ввел дату неправильно.\nВводи в формате 08.05.2025", "", 0)
			return
		}

		botpkg.SendText(bot, ctx, message, fmt.Sprintf("Каждый день можешь тратить по %d рублей", UserMap[message.Chat.ID].Money/(1+int(time.Until(time1).Hours()/24))), "", 0)
		UserMap[message.Chat.ID].Command = "0"

		// Параметры основной клавиатуры. Отправляем, чтобы скрыть клавиатуру с датами
		botpkg.SendKeyboard(bot, ctx, message, nil, "")

	} else { //Отправляем первое сообщение, когда первый раз заходим в функцию
		botpkg.SendText(bot, ctx, message, "Отправь сколько денег осталось🤔\nНапример, 25000", "", 0)
		UserMap[message.Chat.ID].Command = "Ждем сумму"
	}

}
