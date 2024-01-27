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
	Use:     "ilonabot",
	Aliases: []string{"go"},
	Short:   "Start main functions of ilonabot",
	Long:    `Start main functions of ilonabot`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("ilonabot %s started", appVersion)

		// Load Env variables from .dot file
		godotenv.Load(".env")

		token := os.Getenv("SLACK_AUTH_TOKEN")
		channelID := os.Getenv("SLACK_CHANNEL_ID")

		// Create a new client to slack by giving token
		// Set debug to true while developing
		client := slack.New(token, slack.OptionDebug(true))
		// Create the Slack attachment that we will send to the channel
		attachment := slack.Attachment{
			Pretext: "Тут буде розміщуватись назва повідомлення, що сформував бот",
			Text:    "Подія сталась:",
			// Color Styles the Text, making it possible to have like Warnings etc.
			Color: "#36a64f",
			// Fields are Optional extra data!
			Fields: []slack.AttachmentField{
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
			//slack.MsgOptionText("New message from bot", false),
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
