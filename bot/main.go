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
		if update.Message.Document != nil {
			switch update.Message.Document.MimeType {
			case "video/mp4", "video/avi", "video/wmv", "video/mov", "video/qt", "video/mkv", "video/avchd", "video/flv", "video/swf", "video/realvideo":
				if !oauth.IsLoggedIn(update.Message.Chat.ID) {
					loginURL := oauth.GetGoogleLoginURL()
					replyMsg := tgbotapi.NewMessage(update.Message.Chat.ID, "Por favor, faça login no google para Continuar: "+loginURL)
					_, err := bot.Send(replyMsg)
					if err != nil {
						log.Println("Erro ao enviar mensagem de login:", err)
						continue
					}
				} else {
					err := oauth.UploadToGoogleDrive(update.Message.Document.FileID)
					if err != nil {
						log.Println("Erro ao fazer upload do video para o Google Drive: ", err)
						continue
					}

					replyMsg := tgbotapi.NewMessage(update.Message.Chat.ID, "O video foi enviado para o Google Drive com sucesso!")
					_, err = bot.Send(replyMsg)
					if err != nil {
						log.Println("Erro ao enviar mensagem de confimação:", err)
					}
				}
			default:
				replyMsg := tgbotapi.NewMessage(update.Message.Chat.ID, "Desculpe, apenas videos nos formatos MP4, AVI, WMV, MOV, QT, MKV, AVCHD, FLV, SWF e REALVIDEO são suportados.")
				_, err = bot.Send(replyMsg)
				if err != nil {
					log.Println("Erro ao enviar mensagem de formato não suportado:", err)
				}
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
