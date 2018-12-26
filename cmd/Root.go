package cmd

import (
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strings"
)

func init() {
	rootCmd.AddCommand(cmdPrint, cmdProfile)
}

var cmdExit = &cobra.Command{
	Use:   "exit",
	Short: "Exit the program",
	Long:  `Exit the program`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		os.Exit(0)
	},
}

var cmdPrint = &cobra.Command{
	Use:   "print [string to print]",
	Short: "Print anything to the screen",
	Long:  `print is for printing anything back to the screen. For many years people have printed back to the screen.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ReplyMessage("Print: " + strings.Join(args, " "))
	},
}

var cmdProfile = &cobra.Command{
	Use:   "profile",
	Short: "",
	Long:  ``,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		if source.UserID != "" {
			profile, err := bot.GetProfile(source.UserID).Do()
			if err != nil {
				bot.ReplyMessage(
					replyToken,
					linebot.NewTextMessage(err.Error()),
				).Do()
			}
			if _, err := bot.ReplyMessage(
				replyToken,
				linebot.NewTextMessage("Display name: "+profile.DisplayName),
				linebot.NewTextMessage("Status message: "+profile.StatusMessage),
			).Do(); err != nil {
				log.Print(err)
				return
			}
		} else {
			if _, err := bot.ReplyMessage(
				replyToken,
				linebot.NewTextMessage("Bot can't use profile API without user ID"),
			).Do(); err != nil {
				log.Print(err)
				return
			}
		}
	},
}
