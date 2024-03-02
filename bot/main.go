package main

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	bot, err := tgbotapi.NewBotAPI("7038585765:AAEAhlGInnGdG59xHffGubPyR8HN07lotoo")
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true

	log.Printf("Autorizado como %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatal(err)
	}

	for update := range updates {
		if update.Message.Document != nil && update.Message.Document.MimeType == "video/mp4" {
			err := handleVideoChunks(bot, update.Message)
			if err != nil {
				log.Println("Erro ao lidar com os chunks de vídeo:", err)
			}
		}
	}
}

func handleVideoChunks(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) error {
	replyMsg := tgbotapi.NewMessage(msg.Chat.ID, "O vídeo foi unido com sucesso!")
	_, err := bot.Send(replyMsg)
	if err != nil {
		return err
	}

	return nil
}
