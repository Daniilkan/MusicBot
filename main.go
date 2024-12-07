package main

import (
	"TgBot/core"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	BOT_TOKEN, exists := os.LookupEnv("BOT_TOKEN")
	if !exists {
		log.Fatal("BOT_TOKEN not found in environment variables")
	}

	bot, err := telego.NewBot(BOT_TOKEN, telego.WithDefaultDebugLogger())
	if err != nil {
		log.Fatalf("Error creating bot: %v", err)
	}

	updates, err := bot.UpdatesViaLongPolling(nil)
	if err != nil {
		log.Fatalf("Error getting updates: %v", err)
	}

	bh, err := th.NewBotHandler(bot, updates)
	if err != nil {
		log.Fatalf("Error creating bot handler: %v", err)
	}

	defer bot.StopLongPolling()
	defer bh.Stop()

	user_states := make(map[int64]string)

	bh.HandleMessage(func(bot *telego.Bot, message telego.Message) {
		log.Printf("Received message: %s from: %s", message.Text, message.From.FirstName)

		_, err = bot.SendMessage(tu.Messagef(
			tu.ID(message.Chat.ID),
			"Hello %s!", message.From.FirstName,
		).WithReplyMarkup(tu.InlineKeyboard(
			tu.InlineKeyboardRow(tu.InlineKeyboardButton("Search song by name").WithCallbackData("searchName")), tu.InlineKeyboardRow(tu.InlineKeyboardButton("Get top tracks").WithCallbackData("getTop"))),
		))

		if err != nil {
			log.Printf("Error sending message: %v", err)
		}
		user_states[message.Chat.ChatID().ID] = "None"
	}, th.CommandEqual("start"))

	bh.HandleCallbackQuery(func(bot *telego.Bot, query telego.CallbackQuery) {
		_, err := bot.SendMessage(tu.Message(tu.ID(query.From.ID), "Type song name you search"))
		if err == nil {
			user_states[query.From.ID] = "WaitName"
		} else {
			log.Fatal("failed to send message")
		}
	}, th.AnyCallbackQueryWithMessage(), th.CallbackDataEqual("searchName"))

	bh.HandleMessage(func(bot *telego.Bot, message telego.Message) {
		if user_states[message.Chat.ChatID().ID] == "WaitName" {
			text := ""
			for _, track := range core.SearchName(message.Text) {
				text += fmt.Sprintf(`<a href="%s">%s</a> by %s`, track.Url, track.Name, track.Artist) + "\n"
			}
			_, err := bot.SendMessage(tu.Message(tu.ID(message.GetChat().ID), text).
				WithParseMode("HTML").WithLinkPreviewOptions(&telego.LinkPreviewOptions{IsDisabled: true}))
			if err != nil {
				log.Printf("Error sending message on callback: %v", err)
			}
			user_states[message.Chat.ChatID().ID] = "None"
		}
	}, th.AnyMessage())

	bh.HandleCallbackQuery(func(bot *telego.Bot, query telego.CallbackQuery) {
		text := ""
		for i, track := range core.GetTopTracks() {
			if i == 10 {
				break
			}
			text += fmt.Sprintf(`<a href="%s">%s</a> by %s`, track.URL, track.Name, track.Artist.Name) + "\n"
		}
		_, err := bot.SendMessage(tu.Message(tu.ID(query.From.ID), text).WithParseMode("HTML").WithLinkPreviewOptions(&telego.LinkPreviewOptions{IsDisabled: true}))
		if err != nil {
			log.Fatal("Error")
		}
	})

	bh.Start()
}
