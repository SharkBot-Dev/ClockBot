package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"slices"
	"strconv"
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
		{
			Name:        "about",
			Description: "Botの情報を取得します",
		},
	}
	startTime time.Time
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

					var usedTypes []string

					for _, channel := range channels {
						var newName string

						if strings.Contains(channel.Topic, "time-ch") {
							if slices.Contains(usedTypes, "time-ch") {
								continue
							}
							newName = nowtime
							usedTypes = append(usedTypes, "time-ch")
						} else if strings.Contains(channel.Topic, "day-ch") {
							if slices.Contains(usedTypes, "day-ch") {
								continue
							}
							newName = nowday
							usedTypes = append(usedTypes, "day-ch")
						} else if strings.Contains(channel.Topic, "clock-ch") {
							if slices.Contains(usedTypes, "clock-ch") {
								continue
							}
							newName = nowclock
							usedTypes = append(usedTypes, "clock-ch")
						}

						if strings.HasSuffix(channel.Name, "#time-ch") {
							if slices.Contains(usedTypes, "time-ch") {
								continue
							}
							newName = nowtime + "#time-ch"
							usedTypes = append(usedTypes, "time-ch")
						} else if strings.HasSuffix(channel.Name, "#day-ch") {
							if slices.Contains(usedTypes, "day-ch") {
								continue
							}
							newName = nowday + "#day-ch"
							usedTypes = append(usedTypes, "day-ch")
						} else if strings.HasSuffix(channel.Name, "#clock-ch") {
							if slices.Contains(usedTypes, "clock-ch") {
								continue
							}
							newName = nowclock + "#clock-ch"
							usedTypes = append(usedTypes, "clock-ch")
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

	session.AddHandler(func(s *discordgo.Session, channel *discordgo.ChannelUpdate) {
		nowtime := time.Now().Format(nowTimeFormat)
		nowday := time.Now().Format(nowDayFormat)
		nowclock := time.Now().Format(nowDataFormat)

		var newName string

		if strings.Contains(channel.Topic, "time-ch") {
			newName = nowtime
		} else if strings.Contains(channel.Topic, "day-ch") {
			newName = nowday
		} else if strings.Contains(channel.Topic, "clock-ch") {
			newName = nowclock
		}

		if strings.HasSuffix(channel.Name, "#time-ch") {
			newName = nowtime + "#time-ch"
		} else if strings.HasSuffix(channel.Name, "#day-ch") {
			newName = nowday + "#day-ch"
		} else if strings.HasSuffix(channel.Name, "#clock-ch") {
			newName = nowclock + "#clock-ch"
		}

		if newName != "" && channel.Name != newName {
			_, err := s.ChannelEdit(channel.ID, &discordgo.ChannelEdit{Name: newName})
			if err != nil {
				log.Printf("チャンネル %s (%s) の編集に失敗: %v", channel.Name, channel.ID, err)
			}
			time.Sleep(1 * time.Second)
		}
	})

	session.AddHandler(func(s *discordgo.Session, channel *discordgo.ChannelCreate) {
		nowtime := time.Now().Format(nowTimeFormat)
		nowday := time.Now().Format(nowDayFormat)
		nowclock := time.Now().Format(nowDataFormat)

		var newName string

		if strings.Contains(channel.Topic, "time-ch") {
			newName = nowtime
		} else if strings.Contains(channel.Topic, "day-ch") {
			newName = nowday
		} else if strings.Contains(channel.Topic, "clock-ch") {
			newName = nowclock
		}

		if strings.HasSuffix(channel.Name, "#time-ch") {
			newName = nowtime + "#time-ch"
		} else if strings.HasSuffix(channel.Name, "#day-ch") {
			newName = nowday + "#day-ch"
		} else if strings.HasSuffix(channel.Name, "#clock-ch") {
			newName = nowclock + "#clock-ch"
		}

		if newName != "" && channel.Name != newName {
			_, err := s.ChannelEdit(channel.ID, &discordgo.ChannelEdit{Name: newName})
			if err != nil {
				log.Printf("チャンネル %s (%s) の編集に失敗: %v", channel.Name, channel.ID, err)
			}
			time.Sleep(1 * time.Second)
		}
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
									Name:   "チャンネル表示時計 (1)",
									Value:  "特定のチャンネルのチャンネル名が現在時刻に編集されます。\n以下の文字列をチャンネルトピックに含ませると編集されるようになります。\n・時刻 (`time-ch`)\n・日付 (`day-ch`)\n・日付と時刻 (`clock-ch`)",
									Inline: false,
								},
								{
									Name:   "チャンネル表示時計 (2)",
									Value:  "特定のチャンネルもチャンネル名が現在時刻に編集されます。\n以下の文字列をチャンネル名の最後に置くと編集されるようになります。\n・時刻 (`#time-ch`)\n・日付 (`#day-ch`)\n・日付と時刻 (`#clock-ch`)",
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
		case "about":
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						{
							Title: "時計Botの情報",
							Color: 16769280,
							Fields: []*discordgo.MessageEmbedField{
								{
									Name:   "サーバー数",
									Value:  strconv.Itoa(len(s.State.Guilds)) + "サーバー",
									Inline: false,
								},
								{
									Name:   "起動時間",
									Value:  strconv.Itoa(len(s.State.Guilds)) + "サーバー",
									Inline: false,
								},
							},
						},
					},
					Flags: discordgo.MessageFlagsEphemeral,
				},
			})

			uptime := time.Since(startTime)
			uptimeStr := fmt.Sprintf("%d時間%d分", int(uptime.Hours()), int(uptime.Minutes())%60)
			s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Embeds: &[]*discordgo.MessageEmbed{
					{
						Title: "時計Botの情報",
						Color: 16769280,
						Fields: []*discordgo.MessageEmbedField{
							{
								Name:   "サーバー数",
								Value:  strconv.Itoa(len(s.State.Guilds)) + " サーバー",
								Inline: false,
							},
							{
								Name:   "起動時間",
								Value:  uptimeStr,
								Inline: false,
							},
						},
					},
				},
			})
		}
	})

	if err := session.Open(); err != nil {
		log.Fatalf("Discordセッションのオープンに失敗: %v", err)
	}

	startTime = time.Now()

	defer session.Close()

	log.Println("ボットが起動しました。Ctrl+Cで終了します。")

	waitForExitSignal()
}

func waitForExitSignal() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
