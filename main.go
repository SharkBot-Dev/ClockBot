package main

import (
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/SharkBot-Dev/ClockBot/lib"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var nowDataFormat = "2006年01月02日 15時04分"
var nowDayFormat = "2006年01月02日"
var nowTimeFormat = "15時04分"

var (
	session  *discordgo.Session
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "help",
			Description: "Botの使い方を知ります",
		},
	}
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		log.Fatal("環境変数 DISCORD_TOKEN が設定されていません")
	}

	sessionManager := &lib.DiscordSessionManager{}
	session = sessionManager.InitializeSession(token)

	session.Identify.Intents = discordgo.IntentsGuilds

	session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Print("起動しました。")

		go func() {
			s.ApplicationCommandBulkOverwrite(s.State.Application.ID, "", commands)
			log.Print("スラッシュコマンドを同期しました。")
		}()

		go func() {
			for {
				nowtime := time.Now().Format(nowDataFormat)
				s.UpdateCustomStatus(nowtime)
				time.Sleep(10 * time.Second)
			}
		}()

		go func() {
			time.Sleep(5 * time.Second)

			for {
				nowtime := time.Now().Format(nowTimeFormat)
				nowday := time.Now().Format(nowDayFormat)
				nowclock := time.Now().Format(nowDataFormat)

				guilds := s.State.Guilds
				for _, guild := range guilds {
					channels := guild.Channels

					for _, channel := range channels {
						var newName string

						if strings.Contains(channel.Topic, "time-ch") {
							newName = nowtime
						} else if strings.Contains(channel.Topic, "day-ch") {
							newName = nowday
						} else if strings.Contains(channel.Topic, "clock-ch") {
							newName = nowclock
						}

						if newName != "" && channel.Name != newName {
							_, err := s.ChannelEdit(channel.ID, &discordgo.ChannelEdit{Name: newName})
							if err != nil {
								log.Printf("チャンネル %s (%s) の編集に失敗: %v", channel.Name, channel.ID, err)
							}
							time.Sleep(1 * time.Second)
						}
					}
				}

				time.Sleep(30 * time.Minute)
			}
		}()
	})

	session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type != discordgo.InteractionApplicationCommand {
			return
		}

		commandName := i.ApplicationCommandData().Name
		switch commandName {
		case "help":
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						{
							Title: "時計Botの使い方",
							Color: 16769280,
							Fields: []*discordgo.MessageEmbedField{
								{
									Name:   "基本的な時計",
									Value:  "時計Botのステータスに現在時刻が表示されています",
									Inline: false,
								},
								{
									Name:   "チャンネル表示時計",
									Value:  "特定のチャンネルのチャンネル名が現在時刻に編集されます。\n以下の文字列をチャンネルトピックに含ませると編集されるようになります。\n・時刻 (`time-ch`)\n・日付 (`day-ch`)\n・日付と時刻 (`clock-ch`)",
									Inline: false,
								},
								{
									Name:   "編集の周期",
									Value:  "編集の周期は以下になっています。\n・ステータス (10秒に一回)\n・チャンネル (30分に一回)",
									Inline: false,
								},
							},
						},
					},
					Flags: discordgo.MessageFlagsEphemeral,
				},
			})
		}
	})

	if err := session.Open(); err != nil {
		log.Fatalf("Discordセッションのオープンに失敗: %v", err)
	}
	defer session.Close()

	log.Println("ボットが起動しました。Ctrl+Cで終了します。")

	waitForExitSignal()
}

func waitForExitSignal() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
