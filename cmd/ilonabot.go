package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/slack-go/slack"
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
		channelID := os.Getenv("SLACK_CHANNEL_ID")
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
		socketClient.Run()

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
	rootCmd.AddCommand(ilonabotCmd)
}
