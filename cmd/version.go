package cmd

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/slack-go/slack"
	"github.com/spf13/cobra"
)

var appVersion = "Version"

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print IlonaBot version",
	Long:  `Get IlonaBot version to console and Slack channel`,
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
		client := slack.New(token, slack.OptionDebug(false))

		// Create the Slack attachment that we will send to the channel
		attachment := slack.Attachment{
			// Pretext: "Slack бот Ilona запущено!",
			// Text:    "Подія сталась:",
			// Color Styles the Text, making it possible to have like Warnings etc.
			Color: "#03574B",
			// Fields are Optional extra data!
			Fields: []slack.AttachmentField{
				{
					Title: "Версія:",
					Value: appVersion,
				},
			},
		}
		// PostMessage will send the message away.
		// First parameter is just the channelID, makes no sense to accept it
		_, _, err := client.PostMessage(
			channelID,
			// uncomment the item below to add a extra Header to the message, try it out :)
			slack.MsgOptionText("Запущена команда визначення версії Slack боту Ilona", true),
			slack.MsgOptionAttachments(attachment),
		)

		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
