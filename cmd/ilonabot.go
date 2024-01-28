package cmd

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
	"github.com/spf13/cobra"
)

// ilonabotCmd represents the ilonabot command
var ilonabotCmd = &cobra.Command{
	Use: "go",
	// Aliases: []string{"go"},
	Short: "Start main functions of ilonabot",
	Long:  `This type of command bot is triggered when a command begins with a slash. It is the bot for slack command interface.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("ilonabot %s started \n", appVersion)

		// Load Env variables from .dot file
		godotenv.Load(".env")

		// Ініціалізація глобальних змінних
		token := os.Getenv("SLACK_AUTH_TOKEN")
		appToken := os.Getenv("SLACK_APP_TOKEN")
		// Передача змінних як аргументів у функцію cmd.Execute()

		// Створення нового клієнта для Slack за допомогою токена
		// Встановлення debug в true під час розробки
		client := slack.New(token, slack.OptionDebug(true), slack.OptionAppLevelToken(appToken))
		// go-slack comes with a SocketMode package that we need to use that accepts a Slack client and outputs a Socket mode client instead
		socketClient := socketmode.New(
			client,
			socketmode.OptionDebug(true),
			// Option to set a custom logger
			socketmode.OptionLog(log.New(os.Stdout, "socketmode: ", log.Lshortfile|log.LstdFlags)),
		)

		// Створюємо контекст, який можна використовувати для скасування goroutine
		ctx, cancel := context.WithCancel(context.Background())
		// Зробити так, щоб це скасування викликалося належним чином у реальній програмі, плавне завершення роботи тощо
		defer cancel()

		go func(ctx context.Context, client *slack.Client, socketClient *socketmode.Client) {
			// Create a for loop that selects either the context cancellation or the events incomming
			for {
				select {
				// inscase context cancel is called exit the goroutine
				case <-ctx.Done():
					log.Println("Shutting down socketmode listener")
					return
				case event := <-socketClient.Events:
					// We have a new Events, let's type switch the event
					// Add more use cases here if you want to listen to other events.
					switch event.Type {
					// handle EventAPI events
					case socketmode.EventTypeEventsAPI:
						// The Event sent on the channel is not the same as the EventAPI events so we need to type cast it
						eventsAPIEvent, ok := event.Data.(slackevents.EventsAPIEvent)
						if !ok {
							log.Printf("Could not type cast the event to the EventsAPIEvent: %v\n", event)
							continue
						}
						// We need to send an Acknowledge to the slack server
						socketClient.Ack(*event.Request)

						// Now we have an Events API event, but this event type can in turn be many types, so we actually need another type switch
						// log.Println(eventsAPIEvent)

						// Now we have an Events API event, but this event type can in turn be many types, so we actually need another type switch
						err := handleEventMessage(eventsAPIEvent, client)
						if err != nil {
							// Replace with actual err handeling
							log.Fatal(err)
						}

					}

				}
			}
		}(ctx, client, socketClient)

		socketClient.Run()
	},
}

// handleEventMessage will take an event and handle it properly based on the type of event
func handleEventMessage(event slackevents.EventsAPIEvent, client *slack.Client) error {
	switch event.Type {
	// First we check if this is an CallbackEvent
	case slackevents.CallbackEvent:

		innerEvent := event.InnerEvent
		// Yet Another Type switch on the actual Data to see if its an AppMentionEvent
		switch ev := innerEvent.Data.(type) {
		case *slackevents.AppMentionEvent:
			// The application has been mentioned since this Event is a Mention event
			err := handleAppMentionEvent(ev, client)
			if err != nil {
				return err
			}
		}
	default:
		return errors.New("unsupported event type")
	}
	return nil
}

// handleAppMentionEvent is used to take care of the AppMentionEvent when the bot is mentioned
func handleAppMentionEvent(event *slackevents.AppMentionEvent, client *slack.Client) error {

	// Grab the user name based on the ID of the one who mentioned the bot
	user, err := client.GetUserInfo(event.User)
	if err != nil {
		return err
	}
	// Check if the user said Hello to the bot
	text := strings.ToLower(event.Text)

	// Create the attachment and assigned based on the message
	attachment := slack.Attachment{}
	// Add Some default context like user who mentioned the bot
	attachment.Fields = []slack.AttachmentField{
		{
			Title: "Дата та час звернення",
			Value: time.Now().Format("02-01-2006 15:04:05"),
		}, {
			Title: "Ініціатор звернення",
			Value: user.Name,
		},
	}
	if strings.Contains(text, "hello") || strings.Contains(text, "привіт") || strings.Contains(text, "вітаю") {
		// Greet the user
		// attachment.Text = "Привіт"
		attachment.Pretext = fmt.Sprintf("Привіт %s", user.RealName)
		attachment.Color = "#4af030"
	} else {
		// Send a message to the user
		attachment.Text = fmt.Sprintf("Як я можу Вам допомогти %s?", user.RealName)
		attachment.Pretext = "Знаходжусь в очікуванні команд"
		attachment.Color = "#3d3d3d"
	}
	// Send the message to the channel
	// The Channel is available in the event message
	_, _, err = client.PostMessage(event.Channel, slack.MsgOptionAttachments(attachment))
	if err != nil {
		return fmt.Errorf("failed to post message: %w", err)
	}
	return nil
}

func init() {
	rootCmd.AddCommand(ilonabotCmd)
}
