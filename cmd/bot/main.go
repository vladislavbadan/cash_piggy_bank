package main

import (
	"context"
	"fmt"
	"os"

	"cash_piggy_bank/internal/botpkg"
	"cash_piggy_bank/internal/botpkg/processing"
	"cash_piggy_bank/internal/domain"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
)

func main() {
	ctx := context.Background()
	botToken := os.Getenv("TOKEN")

	// Создать бота с включенной отладкой
	// Регистратор по умолчанию может раскрыть конфиденциальную информацию, используйте только в разработке
	bot, err := telego.NewBot(botToken /*, telego.WithDefaultDebugLogger()*/)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Получить пользователя бота
	botUser, err := bot.GetMe(ctx)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("Bot User: %+v\n", botUser)

	// Получаем канал обновлений
	updates, _ := bot.UpdatesViaLongPolling(ctx, nil)

	// Создать обработчик бота
	bh, _ := th.NewBotHandler(bot, updates)

	//Создаем карту пользователей, чтобы знать, кто в какую команду заходит
	UserMap := make(map[int64]*domain.User)

	//Создаем структуру курса валют, чтобы не отправлять лишнего
	var exchange *processing.Rates = &processing.Rates{Day: 0, Money: ""}

	// Обработка любого сообщения
	bh.HandleMessage(func(ctx *th.Context, message telego.Message) error {

		//Проверяем структуру на пустоту и если она пустая, то заполняем нулями
		if UserMap[message.Chat.ID] == UserMap[1] {
			UserMap[message.Chat.ID] = &domain.User{Command: "", Money: 0, TargetMap: make(map[string]int), DbChangeTargetId: 0}
		}

		if message.Text == "/start" || message.Text == "/Start" || message.Text == "Вернуться назад ⏪" {
			UserMap[message.Chat.ID].Command = ""
			botpkg.SendKeyboard(bot, ctx, message, nil, "")
		}

		/* Copy sent messages back to the User
		_, _ = bot.SendMessage(ctx,
			tu.Message(
				tu.ID(message.Chat.ID),
				"Это тестовое смс",
			),
		) */

		// Обработка сообщений. Внутри нее вызываются нужные функции и команды в зависимости от сообщения
		processing.Processing(message.Text, bot, ctx, message, UserMap, exchange)
		return nil
	})

	// Прекратить обработку обновлений при выходе
	defer func() { _ = bh.Stop() }()

	// Начать обработку обновлений
	_ = bh.Start()
}
