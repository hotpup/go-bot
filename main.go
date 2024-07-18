package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

func main() {
	sess, err := discordgo.New("Bot ")
	if err != nil {
		log.Fatal(err)
	}

	sess.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID { return }
		
		if strings.HasPrefix(m.Content, "!timer ") {
			// finds and sets duration. starts timer setting process
			parts := strings.Split(m.Content, " ")
			if len(parts) <= 3 {
				s.ChannelMessageSend(m.ChannelID, "Invalid command format. Example: `!timer 1h30m @user reminder message`")
				return
			}
			duration, err := time.ParseDuration(parts[1])
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "Invalid duration format. Example: `!timer 1h30m @user reminder message`")
				return
			}
			
			time.AfterFunc(duration, func() {
				s.ChannelMessageSend(m.ChannelID, strings.Join(parts[2:], " "))
			})

			s.ChannelMessageSend(m.ChannelID, "Reminder timer set.")
		}

		if strings.HasPrefix(m.Content, "!daily ") {
			// finds and sets duration. starts timer setting process
			parts := strings.Split(m.Content, " ")
			if len(parts) <= 3 {
				s.ChannelMessageSend(m.ChannelID, "Invalid command format. Example: `!daily 15:04 @user reminder message`")
				return
			}
			
			reminderTime, err := time.Parse("15:04", parts[1])
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "Invalid time format. Example: `!daily 15:04 @user reminder message`")
				return
			}

			// get current time
			now := time.Now()
			// get the time of the reminder today
			if now.Hour() > reminderTime.Hour() || (now.Hour() == reminderTime.Hour() && now.Minute() >= reminderTime.Minute()) {
				reminderTime = reminderTime.AddDate(0, 0, 1)
			}
			reminderTime = time.Date(now.Year(), now.Month(), now.Day(), reminderTime.Hour(), reminderTime.Minute(), 0, 0, now.Location())
			durationUntilReminder := reminderTime.Sub(now)

			// create a function that sends the reminder and sets the next reminder
			var sendReminder func()
			sendReminder = func() {
				s.ChannelMessageSend(m.ChannelID, strings.Join(parts[2:], " "))
				time.AfterFunc(24 * time.Hour, sendReminder)
			}

			// set the first reminder
			time.AfterFunc(durationUntilReminder, sendReminder)

			s.ChannelMessageSend(m.ChannelID, "Daily reminder set.")
		}
	})

	sess.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	err = sess.Open()
	if err != nil {
		log.Fatal(err)
	}

	defer sess.Close()

	fmt.Println("Online")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	fmt.Println("Offline")
}