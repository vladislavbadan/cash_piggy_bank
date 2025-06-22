package botpkg

import (
	"context"
	"log"
	"time"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

func SendText(bot *telego.Bot, ctx *th.Context, message telego.Message, textMessage string, cases string, ChatId int64) {
	switch cases {
	case "нотификация":
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_, err := bot.SendMessage(ctx,
			tu.Message(
				tu.ID(ChatId),
				textMessage,
			),
		)
		if err != nil {
			log.Println(err)
		}

	default:
		_, err := bot.SendMessage(ctx, &telego.SendMessageParams{
			ChatID:    message.Chat.ChatID(),
			Text:      textMessage,
			ParseMode: telego.ModeHTML,
		})
		if err != nil {
			log.Println(err)
		}
	}

}
