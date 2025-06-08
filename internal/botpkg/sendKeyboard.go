package botpkg

import (
	"context"
	"time"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

// Создает и отправляет клавиатуру пользователю в зависимости от кейса
func SendKeyboard(bot *telego.Bot, ctx *th.Context, message telego.Message, params []string, cases string) {

	var keyboard telego.ReplyMarkup
	var textMessageKeyboard string

	switch cases {
	case "накопления":
		keyboard = tu.Keyboard(
			tu.KeyboardRow( // Row 1
				tu.KeyboardButton("Добавить цель"),
				tu.KeyboardButton("Посмотреть цели"),
			),
			tu.KeyboardRow( // Row 2
				tu.KeyboardButton("Обновить цель"),
				tu.KeyboardButton("Удалить цель"),
			),
			tu.KeyboardRow( // Row 3
				tu.KeyboardButton("Вернуться назад ⏪"),
			),
		).WithResizeKeyboard().WithInputFieldPlaceholder("Выбери дату на клавиатуре")

		textMessageKeyboard = "Что будем делать? 👇"
	case "calculationDate":
		keyboard = tu.Keyboard(
			tu.KeyboardRow( // Row 1
				tu.KeyboardButton(time.Now().AddDate(0, 0, 2).Format("02.01.2006")), // Column 1
				tu.KeyboardButton(time.Now().AddDate(0, 0, 4).Format("02.01.2006")), // Column 2
				tu.KeyboardButton(time.Now().AddDate(0, 0, 6).Format("02.01.2006")), // Column 3
			),
			tu.KeyboardRow( // Row 2
				tu.KeyboardButton(time.Now().AddDate(0, 0, 8).Format("02.01.2006")),  // Column 1
				tu.KeyboardButton(time.Now().AddDate(0, 0, 10).Format("02.01.2006")), // Column 2
				tu.KeyboardButton(time.Now().AddDate(0, 0, 12).Format("02.01.2006")), // Column 3
			),
			tu.KeyboardRow( // Row 3
				tu.KeyboardButton(time.Now().AddDate(0, 0, 14).Format("02.01.2006")), // Column 1
				tu.KeyboardButton(time.Now().AddDate(0, 0, 16).Format("02.01.2006")), // Column 2
				tu.KeyboardButton(time.Now().AddDate(0, 0, 18).Format("02.01.2006")), // Column 3
			),
			tu.KeyboardRow( // Row 4
				tu.KeyboardButton(time.Now().AddDate(0, 0, 20).Format("02.01.2006")), // Column 1
				tu.KeyboardButton(time.Now().AddDate(0, 0, 22).Format("02.01.2006")), // Column 2
				tu.KeyboardButton("Вернуться назад ⏪"),                               // Column 3
			),
		).WithResizeKeyboard().WithInputFieldPlaceholder("Выбери дату на клавиатуре")

		textMessageKeyboard = "Сколько нужно протянуть на эти деньги? Укажи дату в формате 08.05.2025 или выбери из списка👇"
	case "цели":
		switch len(params) {
		case 1:
			keyboard = tu.Keyboard(
				tu.KeyboardRow( // Row 1
					tu.KeyboardButton(params[0]),           // Column 1
					tu.KeyboardButton("Вернуться назад ⏪"), // Column 2
				),
			).WithResizeKeyboard().WithInputFieldPlaceholder("Выбери цель на клавиатуре")
		case 2:
			keyboard = tu.Keyboard(
				tu.KeyboardRow( // Row 1
					tu.KeyboardButton(params[0]), // Column 1
					tu.KeyboardButton(params[1]), // Column 2
				),
				tu.KeyboardRow( // Row 2
					tu.KeyboardButton("Вернуться назад ⏪"), // Column 1
				),
			).WithResizeKeyboard().WithInputFieldPlaceholder("Выбери цель на клавиатуре")

		case 3:
			keyboard = tu.Keyboard(
				tu.KeyboardRow( // Row 1
					tu.KeyboardButton(params[0]), // Column 1
					tu.KeyboardButton(params[1]), // Column 2
					tu.KeyboardButton(params[2]), // Column 3
				),
				tu.KeyboardRow( // Row 2
					tu.KeyboardButton("Вернуться назад ⏪"), // Column 1
				),
			).WithResizeKeyboard().WithInputFieldPlaceholder("Выбери цель на клавиатуре")
		case 4:
			keyboard = tu.Keyboard(
				tu.KeyboardRow( // Row 1
					tu.KeyboardButton(params[0]), // Column 1
					tu.KeyboardButton(params[1]), // Column 2
					tu.KeyboardButton(params[2]), // Column 3
				),
				tu.KeyboardRow( // Row 2
					tu.KeyboardButton(params[3]),           // Column 1
					tu.KeyboardButton("Вернуться назад ⏪"), // Column 2
				),
			).WithResizeKeyboard().WithInputFieldPlaceholder("Выбери цель на клавиатуре")
		case 5:
			keyboard = tu.Keyboard(
				tu.KeyboardRow( // Row 1
					tu.KeyboardButton(params[0]), // Column 1
					tu.KeyboardButton(params[1]), // Column 2
					tu.KeyboardButton(params[2]), // Column 3
				),
				tu.KeyboardRow( // Row 2
					tu.KeyboardButton(params[3]), // Column 1
					tu.KeyboardButton(params[4]), // Column 2
				),
				tu.KeyboardRow( // Row 3
					tu.KeyboardButton("Вернуться назад ⏪"), // Column 1
				),
			).WithResizeKeyboard().WithInputFieldPlaceholder("Выбери цель на клавиатуре")
		case 6:
			keyboard = tu.Keyboard(
				tu.KeyboardRow( // Row 1
					tu.KeyboardButton(params[0]), // Column 1
					tu.KeyboardButton(params[1]), // Column 2
					tu.KeyboardButton(params[2]), // Column 3
				),
				tu.KeyboardRow( // Row 2
					tu.KeyboardButton(params[3]), // Column 1
					tu.KeyboardButton(params[4]), // Column 2
					tu.KeyboardButton(params[5]), // Column 3
				),
				tu.KeyboardRow( // Row 3
					tu.KeyboardButton("Вернуться назад ⏪"), // Column 1
				),
			).WithResizeKeyboard().WithInputFieldPlaceholder("Выбери цель на клавиатуре")
		case 7:
			keyboard = tu.Keyboard(
				tu.KeyboardRow( // Row 1
					tu.KeyboardButton(params[0]), // Column 1
					tu.KeyboardButton(params[1]), // Column 2
					tu.KeyboardButton(params[2]), // Column 3
				),
				tu.KeyboardRow( // Row 2
					tu.KeyboardButton(params[3]), // Column 1
					tu.KeyboardButton(params[4]), // Column 2
					tu.KeyboardButton(params[5]), // Column 3
				),
				tu.KeyboardRow( // Row 3
					tu.KeyboardButton(params[6]),           // Column 1
					tu.KeyboardButton("Вернуться назад ⏪"), // Column 2
				),
			).WithResizeKeyboard().WithInputFieldPlaceholder("Выбери цель на клавиатуре")
		case 8:
			keyboard = tu.Keyboard(
				tu.KeyboardRow( // Row 1
					tu.KeyboardButton(params[0]), // Column 1
					tu.KeyboardButton(params[1]), // Column 2
					tu.KeyboardButton(params[2]), // Column 3
				),
				tu.KeyboardRow( // Row 2
					tu.KeyboardButton(params[3]), // Column 1
					tu.KeyboardButton(params[4]), // Column 2
					tu.KeyboardButton(params[5]), // Column 3
				),
				tu.KeyboardRow( // Row 3
					tu.KeyboardButton(params[6]), // Column 1
					tu.KeyboardButton(params[7]), // Column 2
				),
				tu.KeyboardRow( // Row 4
					tu.KeyboardButton("Вернуться назад ⏪"), // Column 1
				),
			).WithResizeKeyboard().WithInputFieldPlaceholder("Выбери цель на клавиатуре")
		case 9:
			keyboard = tu.Keyboard(
				tu.KeyboardRow( // Row 1
					tu.KeyboardButton(params[0]), // Column 1
					tu.KeyboardButton(params[1]), // Column 2
					tu.KeyboardButton(params[2]), // Column 3
				),
				tu.KeyboardRow( // Row 2
					tu.KeyboardButton(params[3]), // Column 1
					tu.KeyboardButton(params[4]), // Column 2
					tu.KeyboardButton(params[5]), // Column 3
				),
				tu.KeyboardRow( // Row 3
					tu.KeyboardButton(params[6]), // Column 1
					tu.KeyboardButton(params[7]), // Column 2
				),
				tu.KeyboardRow( // Row 4
					tu.KeyboardButton(params[8]),           // Column 1
					tu.KeyboardButton("Вернуться назад ⏪"), // Column 2
				),
			).WithResizeKeyboard().WithInputFieldPlaceholder("Выбери цель на клавиатуре")
		case 10:
			keyboard = tu.Keyboard(
				tu.KeyboardRow( // Row 1
					tu.KeyboardButton(params[0]), // Column 1
					tu.KeyboardButton(params[1]), // Column 2
					tu.KeyboardButton(params[2]), // Column 3
				),
				tu.KeyboardRow( // Row 2
					tu.KeyboardButton(params[3]), // Column 1
					tu.KeyboardButton(params[4]), // Column 2
					tu.KeyboardButton(params[5]), // Column 3
				),
				tu.KeyboardRow( // Row 3
					tu.KeyboardButton(params[6]), // Column 1
					tu.KeyboardButton(params[7]), // Column 2
				),
				tu.KeyboardRow( // Row 4
					tu.KeyboardButton(params[8]),           // Column 1
					tu.KeyboardButton(params[9]),           // Column 1
					tu.KeyboardButton("Вернуться назад ⏪"), // Column 2
				),
			).WithResizeKeyboard().WithInputFieldPlaceholder("Выбери цель на клавиатуре")
		default:
			keyboard = tu.Keyboard(
				tu.KeyboardRow( // Row 1
					tu.KeyboardButton("Не могу отображать больше 10"), // Column 1
					tu.KeyboardButton("Вернуться назад ⏪"),            // Column 2
				),
			).WithResizeKeyboard().WithInputFieldPlaceholder("Отправь название цели текстом")
		}

		textMessageKeyboard = "Выбери цель из списка👇"
	default:
		// Параметры основной клавиатуры
		keyboard = tu.Keyboard(
			tu.KeyboardRow(
				tu.KeyboardButton("Инфо ℹ️"),
				tu.KeyboardButton("Накопления 🤑"),
			),
			tu.KeyboardRow(
				tu.KeyboardButton("Дожить 🥴"),
				tu.KeyboardButton("Курс 💵"),
			),
		).WithResizeKeyboard().WithInputFieldPlaceholder("Используй клавиатуру >>")

		textMessageKeyboard = "Выбирай на клавиатуре 👇"
	}

	messageKey := tu.Message(
		tu.ID(message.Chat.ID),
		textMessageKeyboard,
	).WithReplyMarkup(keyboard)

	_, _ = bot.SendMessage(context.Background(), messageKey)
}
