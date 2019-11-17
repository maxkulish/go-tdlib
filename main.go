package main

import (
	"fmt"
	"log"
	"path/filepath"
	"strconv"

	"github.com/zelenin/go-tdlib/client"
)

const (
	appID   = 187786
	appHash = "e782045df67ba48e441ccb105da8fc85"
)

func main() {
	authorizer := client.ClientAuthorizer()
	go client.CliInteractor(authorizer)

	authorizer.TdlibParameters <- &client.TdlibParameters{
		UseTestDc:              false,
		DatabaseDirectory:      filepath.Join(".tdlib", "database"),
		FilesDirectory:         filepath.Join(".tdlib", "files"),
		UseFileDatabase:        true,
		UseChatInfoDatabase:    true,
		UseMessageDatabase:     true,
		UseSecretChats:         false,
		ApiId:                  appID,
		ApiHash:                appHash,
		SystemLanguageCode:     "en",
		DeviceModel:            "MacBookPro15",
		SystemVersion:          "10.15.1",
		ApplicationVersion:     "1.0.0",
		EnableStorageOptimizer: true,
		IgnoreFileNames:        false,
	}

	logVerbosity := client.WithLogVerbosity(&client.SetLogVerbosityLevelRequest{NewVerbosityLevel: 1})

	tdlibClient, err := client.NewClient(authorizer, logVerbosity)
	if err != nil {
		log.Fatalf("NewClient error: %s", err)
	}

	optionValue, err := tdlibClient.GetOption(&client.GetOptionRequest{Name: "version"})
	if err != nil {
		log.Fatalf("GetOption errors: %s", err)
	}

	log.Printf("TDLib version: %s", optionValue.(*client.OptionValueString).Value)

	me, err := tdlibClient.GetMe()
	if err != nil {
		log.Fatalf("GetMe errors: %s", err)
	}

	log.Printf("Me: %s %s [%s]", me.FirstName, me.PhoneNumber, me.Status.UserStatusType())

	// Updates
	listener := tdlibClient.GetListener()
	defer listener.Close()

	for update := range listener.Updates {
		fmt.Println("[UT]:", update.GetType())
		if update.GetClass() == client.ClassUpdate {

			switch update.GetType() {
			case client.TypeUpdateNewMessage:
				msgRaw := update.(*client.UpdateNewMessage).Message
				log.Printf("id: %d, chat_id: %d, date: %d, message: %+v\n", msgRaw.Id, msgRaw.ChatId, msgRaw.Date, msgRaw.Content)

				ParseMessage(msgRaw)
			}

		}
	}
}

type NormMessage struct {
	ID     int64
	ChatID int64
	Date   int32
	Text   string
}

func ParseMessage(msgRaw *client.Message) {

	log.Println("[MT]", msgRaw.GetType())
	log.Println("[CT]", msgRaw.Content.MessageContentType())

	nm := &NormMessage{
		ID:     msgRaw.Id,
		ChatID: msgRaw.ChatId,
		Date:   msgRaw.Date,
	}

	switch msgRaw.Content.MessageContentType() {
	case client.TypeMessageText:
		mt := msgRaw.Content.(*client.MessageText)
		nm.Text = mt.Text.Text
	}

	rDate, err := strconv.ParseInt(string(nm.Date), 10, 64)
	if err != nil {
		log.Println("can't convert time from unix to string")
	}

	log.Printf("mesID: %d, date: %s, chatID: %d, message: %s",
		nm.ID, rDate, nm.ChatID, nm.Text)
}

