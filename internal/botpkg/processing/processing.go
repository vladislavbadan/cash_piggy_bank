package processing

import (
	"math/rand"
	"time"

	"cash_piggy_bank/internal/api/exchangepkg"
	"cash_piggy_bank/internal/botpkg"

	"cash_piggy_bank/internal/domain"
	"cash_piggy_bank/internal/repository/sqlite"
	"cash_piggy_bank/internal/service"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

var jobGorotine bool

// Структура курса валют для понимания, когда надо отправлять курс, а когда нет.
type Rates struct {
	Day   int
	Money string
}

// Функция которая запускает дуругие функии в зависимости от того, какая команда была отправлена в чат
func Processing(text string, bot *telego.Bot, ctx *th.Context, message telego.Message, UserMap map[int64]*domain.User, exchange *Rates) {

	if text == "Инфо ℹ️" {
		botpkg.SendText(bot, ctx, message, `<b>💰 Я твой персональный помощник в накоплениях 🚀</b>

Я помогу тебе легко копить деньги, контролировать бюджет и достигать финансовых целей!

<b>✨ Основные возможности:</b>
- 📌 Умный учет накоплений: ставь цели и следи за прогрессом
- 💸 Авторасчет бюджета: скажи сколько осталось и дней до зарплаты - я посчитаю дневной лимит
- 🔔 Мотивирующие напоминания: не дам забыть отложить деньги
- 📊 Курсы валют: всегда актуальная информация

Начни копить осознанно прямо сейчас 😊`, "", 0)

	} else if text == "Дожить 🥴" || UserMap[message.Chat.ID].Command == "Ждем сумму" || UserMap[message.Chat.ID].Command == "Считаем сколько осталось" {

		service.Calculation(message.Text, bot, ctx, message, UserMap)
	} else if text == "Накопления 🤑" || UserMap[message.Chat.ID].Command == "Ждем команду для БД" || UserMap[message.Chat.ID].Command == "Ждем сумму для цели" || UserMap[message.Chat.ID].Command == "Добавляем цель" || UserMap[message.Chat.ID].Command == "Ждем название цели" || UserMap[message.Chat.ID].Command == "Ждем ID и удаление" || text == "Добавить цель" || text == "Посмотреть цели" || text == "Обновить цель" || text == "Удалить цель" {
		sqlite.DbQuery(message.Text, bot, ctx, message, UserMap)

	} else if text == "Курс 💵" {
		_, _, day := time.Now().Date()
		if exchange.Day != day {
			exchange.Day = day
			exchange.Money = exchangepkg.ExchangeRates()
		}
		botpkg.SendText(bot, ctx, message, exchange.Money, "", 0)
	}

	stikers := [11]string{"CAACAgIAAxkBAAELjjFnhBjZruG8kzV7C3JpPMg_LwABUqEAAn0DAAJtsEIDvRMu-bZNulM2BA", "CAACAgIAAxkBAAELjjNnhBjlWVszV46bhgGrd0erbjorpAACVQMAArVx2gYTIpfj7IkJBzYE", "CAACAgIAAxkBAAEOFmln8Et4FUg_hvfN02MMvU27T27V7gACKQADWbv8JWiEdiw7SWZ7NgQ", "CAACAgEAAxkBAAEOKPJn9AIxxp5aa2S0rbWkih5ZxJR6bwACHQEAAjgOghHhhIkhaufuiTYE", "CAACAgIAAxkBAAEOKPRn9AI89vRIMNLzWBvtp3VHbPBvpQACAwEAAladvQoC5dF4h-X6TzYE", "CAACAgIAAxkBAAEOKPZn9AJMhBUD_NV7jgEjSwh8vvo6LwACJgADDbbSGRYpFH5xkFugNgQ", "CAACAgIAAxkBAAEOKPhn9AJUmsdhg8cbMscCVqbOGn1OnAACpQEAAhZCawqkjIJTRgnc1jYE", "CAACAgIAAxkBAAEOKPpn9AJePsS5OFKE05IDLZDVNnBlRgACjAADFkJrCkKO_mIXPU3iNgQ", "CAACAgIAAxkBAAEOKPxn9AJpX38NCysXrpj8ZGHd0Wa96QACBAEAAvcCyA8gD3c76avISTYE", "CAACAgIAAxkBAAEOKP5n9AJyQ49gcrJx7IDf7cMjCrfGvQACXwADwDZPExXCdY-iUFE7NgQ", "CAACAgIAAxkBAAEOKQJn9AJ-_TPoHvDkqqHCgvd6I2PuxwACaAADwDZPE0z9PaPnxGmHNgQ"}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	if r.Intn(8) == 1 {
		_, _ = bot.SendSticker(ctx, &telego.SendStickerParams{
			ChatID: tu.ID(message.Chat.ID),
			Sticker: telego.InputFile{
				FileID: stikers[rand.Intn(11)], // Ваш file_id стикера
			},
		})
	}

	if !jobGorotine {
		jobGorotine = true
		service.Notification(message.Text, bot, ctx, message)
	}
}
