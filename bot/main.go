package main

import (
	"log"

	"github.com/gustawillg/bot-storage/oauth"

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
		if update.Message == nil {
			continue
		}

		switch update.Message.Text {
		case "/start":
			reply := "Olá! Envie um vídeo no formato MP4, AVI, WMV, MOV, QT, MKV, AVCHD, FLV, SWF ou REALVIDEO e eu vou enviá-lo para o Google Drive."
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
			bot.Send(msg)
		case "/upload":
			if update.Message.Document != nil {
				switch update.Message.Document.MimeType {
				case "video/mp4", "video/avi", "video/wmv", "video/mov", "video/qt", "video/mkv", "video/avchd", "video/flv", "video/swf", "video/realvideo":
					if !oauth.IsLoggedIn(update.Message.Chat.ID) {
						loginURL := oauth.GetGoogleLoginURL()
						reply := "Por favor, faça login no Google para continuar: " + loginURL
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
						bot.Send(msg)
					} else {
						err := oauth.UploadToGoogleDrive(update.Message.Document.FileID)
						if err != nil {
							log.Println("Erro ao fazer upload do video para o Google Drive: ", err)
							reply := "Desculpe, ocorreu um erro ao fazer o upload do vídeo."
							msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
							bot.Send(msg)

						} else {
							reply := "O vídeo foi enviado para o Google Drive com sucesso!"
							msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
							bot.Send(msg)
						}
					}
				default:
					reply := "Desculpe, apenas vídeos nos formatos MP4, AVI, WMV, MOV, QT, MKV, AVCHD, FLV, SWF e REALVIDEO são suportados."
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
					bot.Send(msg)
				}
			} else {
				reply := "Por favor, envie um vídeo para ser enviado para o Google Drive."
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
				bot.Send(msg)
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
