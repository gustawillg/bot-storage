package main

import (
	"fmt"
	"io"
	"log"
	"math"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	bot, err := tgbotapi.NewBotAPI("7038585765:AAEAhlGInnGdG59xHffGubPyR8HN07lotoo")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Autorizado como %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Panic(err)
	}

	for update := range updates {
		// Verifica se a atualização é uma mensagem de documento (vídeo, por exemplo)
		if update.Message.Document != nil && update.Message.Document.MimeType == "video/mp4" {
			// Se for um vídeo, chama a função para lidar com ele
			handleVideo(bot, update)
		} else {
			// Se não for um vídeo, continua com a lógica para lidar com outros tipos de documentos
			if update.Message != nil {
				log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
				handleDocument(bot, update)
			}
		}
	}

}

func handleVideo(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	panic("unimplemented")

}

func handleDocument(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	if update.Message.Document == nil {
		return
	}

	fileID := update.Message.Document.FileID
	fileConfig := tgbotapi.FileConfig{
		FileID: fileID,
	}

	file, err := bot.GetFile(fileConfig)
	if err != nil {
		log.Println("Erro ao obter o arquivo:", err)
		return
	}

	fileURL := file.Link(bot.Token)

	// Aqui você pode baixar o arquivo e armazená-lo onde desejar
	log.Println("URL do arquivo:", fileURL)

	// Exemplo de envio de vídeo de volta
	sendVideo(bot, update.Message.Chat.ID, fileURL)
}

func sendVideo(bot *tgbotapi.BotAPI, chatID int64, videoURL string) {
	video := tgbotapi.NewVideoUpload(chatID, videoURL)
	video.Caption = "Seu vídeo"

	_, err := bot.Send(video)
	if err != nil {
		log.Println("Erro ao enviar vídeo:", err)
	}
}

// Função para enviar chunks de um vídeo para o Telegram
func sendVideoChunks(bot *tgbotapi.BotAPI, chatID int64, videoFilePath string) error {
	// Abra o arquivo de vídeo
	videoFile, err := os.Open(videoFilePath)
	if err != nil {
		return err
	}
	defer videoFile.Close()

	// Determine o tamanho do arquivo
	fileInfo, err := videoFile.Stat()
	if err != nil {
		return err
	}

	fileSize := fileInfo.Size()

	// Defina o tamanho máximo do chunk (por exemplo, 50 MB)
	chunkSize := int64(50 * 1024 * 1024) // 50 MB
	numChunks := int(math.Ceil(float64(fileSize) / float64(chunkSize)))

	// Envie cada chunk como um documento separado
	for i := 0; i < numChunks; i++ {
		// Calcule os offsets para cada chunk
		offset := int64(i) * chunkSize
		chunkBytes := min(chunkSize, fileSize-offset)

		// Leia os bytes do chunk
		chunkBuffer := make([]byte, chunkBytes)
		_, err := videoFile.ReadAt(chunkBuffer, offset)
		if err != nil && err != io.EOF {
			return err
		}

		// Envie o chunk como um documento para o Telegram
		videoPart := tgbotapi.FileBytes{
			Name:  fmt.Sprintf("part_%d.mp4", i+1),
			Bytes: chunkBuffer,
		}

		videoConfig := tgbotapi.NewDocumentUpload(chatID, videoPart)
		_, err = bot.Send(videoConfig)
		if err != nil {
			return err
		}
	}

	return nil
}

// Função auxiliar para encontrar o mínimo entre dois inteiros
func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}
