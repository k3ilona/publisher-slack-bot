package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/slack-go/slack"
	"github.com/spf13/cobra"
)

// ilonabotCmd represents the ilonabot command
var ilonabotCmd = &cobra.Command{
	Use: "go",
	// Aliases: []string{"go"},
	Short: "Start main functions of ilonabot",
	Long: `This type of command bot is triggered when a command begins with a slash. 
			  It is the bot for slack command interface.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("ilonabot %s started \n", appVersion)

		// Load Env variables from .dot file
		godotenv.Load(".env")

		// Ініціалізація глобальних змінних
		token := os.Getenv("SLACK_AUTH_TOKEN")
		channelID := os.Getenv("SLACK_CHANNEL_ID")
		// Передача змінних як аргументів у функцію cmd.Execute()

		// Створення нового клієнта для Slack за допомогою токена
		// Встановлення debug в true під час розробки
		client := slack.New(token, slack.OptionDebug(true))

		// Create the Slack attachment that we will send to the channel
		attachment := slack.Attachment{
			// Pretext: "Slack бот Ilona запущено!",
			// Text:    "Подія сталась:",
			// Color Styles the Text, making it possible to have like Warnings etc.
			Color: "#F78166",
			// Fields are Optional extra data!
			Fields: []slack.AttachmentField{
				{
					Title: "Версія:",
					Value: appVersion,
				},
				{
					Title: "Дата:",
					Value: time.Now().Format("02-01-2006"),
				},
				{
					Title: "Час:",
					Value: time.Now().Format("15:04:05"),
				},
			},
		}
		// PostMessage will send the message away.
		// First parameter is just the channelID, makes no sense to accept it
		_, timestamp, err := client.PostMessage(
			channelID,
			// uncomment the item below to add a extra Header to the message, try it out :)
			slack.MsgOptionText("Slack бот Ilona запущено!", true),
			slack.MsgOptionAttachments(attachment),
		)

		if err != nil {
			panic(err)
		}

		fmt.Printf("Message sent at %s", timestamp)
	},
}

func init() {
	rootCmd.AddCommand(ilonabotCmd)
}
