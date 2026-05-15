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

var nowDataFormat = "2006年01月02日 15時04分05秒"
var nowDayFormat = "2006年01月02日"
var nowTimeFormat = "15時04分05秒"

var (
	session *discordgo.Session
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
					channels, err := s.GuildChannels(guild.ID)
					if err != nil {
						log.Printf("ギルド %s のチャンネル一覧取得に失敗: %v", guild.ID, err)
						continue
					}

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
