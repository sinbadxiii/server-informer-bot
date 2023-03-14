package main

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"
)

type DiskStatus struct {
	All  uint64 `json:"all"`
	Used uint64 `json:"used"`
	Free uint64 `json:"free"`
}

func DiskUsage(path string) (disk DiskStatus) {
	fs := syscall.Statfs_t{}
	err := syscall.Statfs(path, &fs)
	if err != nil {
		return
	}
	disk.All = fs.Blocks * uint64(fs.Bsize)
	disk.Free = fs.Bavail * uint64(fs.Bsize)
	disk.Used = disk.All - disk.Free
	return
}

const (
	B  = 1
	KB = 1024 * B
	MB = 1024 * KB
	GB = 1024 * MB
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(".env dont exist")
	}
	token := os.Getenv("TOKEN")

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	is_debug := os.Getenv("DEBUG")
	server_name := os.Getenv("SERVER_NAME")

	bot.Debug = false
	if is_debug == "true" {
		bot.Debug = true
	}

	is_log_messages := os.Getenv("LOG_MESSAGES")
	if is_log_messages == "true" {
		log.Printf("Authorized on account %s", bot.Self.UserName)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	//updates, err := bot.GetUpdatesChan(u)

	for {
		disk := DiskUsage("/")
		title := "*" + server_name + "* \n\n" + "Место на жестком диске"
		total := fmt.Sprintf("Всего %.2f GB", float64(disk.All)/float64(GB))
		used := fmt.Sprintf("Используется %.2f GB", float64(disk.Used)/float64(GB))
		free := fmt.Sprintf("Свободно %.2f GB", float64(disk.Free)/float64(GB))

		sysinfo := syscall.Sysinfo_t{}

		err := syscall.Sysinfo(&sysinfo)

		if err != nil {
			fmt.Println("Error:", err)
		}

		titleMemory := "Оперативная память"
		memoryTotal := fmt.Sprintf("Всего %.2f GB", float64(sysinfo.Totalram)/float64(GB))
		memoryFree := fmt.Sprintf("Свободно %.2f GB", float64(sysinfo.Freeram)/float64(GB))

		out := fmt.Sprintf("%s\n\n%s\n%s\n%s\n\n%s\n\n%s\n%s\n", title, total, used, free,
			titleMemory, memoryTotal, memoryFree)

		//msg := tgbotapi.NewMessage(update.Message.Chat.ID, out)

		chat_ids := os.Getenv("CHAT_ID")

	    ids := strings.Split(chat_ids, ",")

	    for _, id := range ids {
			convertId, _ := strconv.ParseInt(id,10, 64)

			msg := tgbotapi.NewMessage(convertId, out)
			msg.ParseMode = "markdown"

			//msg.ReplyToMessageID = update.Message.MessageID

			bot.Send(msg)
		}

		time.Sleep(time.Hour * 1)
	}
}
