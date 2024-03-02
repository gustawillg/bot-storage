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
		// Verifica se a atualização é uma mensagem de documento (vídeo, por exemplo)
		if update.Message.Document != nil && update.Message.Document.MimeType == "video/mp4" {
			// Se for um vídeo, chama a função para lidar com ele
			err := handleVideoChunks(bot, update.Message)
			if err != nil {
				log.Println("Erro ao lidar com os chunks de vídeo:", err)
			}
		}
	}
}

// Função para lidar com os chunks de um vídeo
func handleVideoChunks(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) error {
	// Aqui você pode implementar a lógica para juntar os chunks de vídeo
	// Receba os chunks enviados pelos usuários e os una em um único vídeo

	// Exemplo simples: Enviar uma mensagem de confirmação
	replyMsg := tgbotapi.NewMessage(msg.Chat.ID, "O vídeo foi unido com sucesso!")
	_, err := bot.Send(replyMsg)
	if err != nil {
		return err
	}

	return nil
}
