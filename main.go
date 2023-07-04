package main

import (
	"fmt"
	"kmuttBot/functions"
	"kmuttBot/types/payload"
	"kmuttBot/utils/config"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/go-co-op/gocron"
)

var dg *discordgo.Session
var prevGrade *payload.Welcome
var channel *discordgo.Channel

func main() {
	// Discord bot token
	token := config.C.BotToken

	// Create a new Discord session using the provided bot token.
	var err error
	dg, err = discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("error creating Discord session:", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// Receive message events in guilds (server)
	//dg.Identify.Intents = discordgo.IntentGuildMessages
	dg.Identify.Intents = discordgo.IntentDirectMessages

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection:", err)
		return
	}

	// Schedule gradeCheck
	s := gocron.NewScheduler(time.UTC)
	_, err = s.Every(3).Seconds().Do(func() {
		GradeCheck(channel, s)
	})
	if err != nil {
		fmt.Println("error scheduling gradeCheck:", err)
		return
	}
	s.StartAsync()

	// Wait here until CTRL-C or other termination signal is received.
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == config.C.Text1 {
		// Create the private channel with the user who sent the message.
		var err error
		channel, err = s.UserChannelCreate(m.Author.ID)
		if err != nil {
			// If an error occurred, we failed to create the channel.
			fmt.Println("error creating channel:", err)
			s.ChannelMessageSend(
				m.ChannelID,
				"Something went wrong while sending the DM!",
			)
			return
		}

		// Send the initial message
		_, err = s.ChannelMessageSend(channel.ID, "Grade 📌: year 2/2")
		if err != nil {
			fmt.Println("error sending DM message:", err)
			return
		}

		// GetGrade
		grade, gradeErr := functions.GetGrade()
		if gradeErr != nil {
			fmt.Println(gradeErr)
			fmt.Println("error getting grade")
			return
		}
		for _, v := range grade.GradeInfo.Grades {
			for _, v2 := range v.Courses {
				_, err = s.ChannelMessageSend(channel.ID, v2.CourseCode+" "+v2.CourseNameEn+" Grade : "+"**"+v2.CourseGrade+"**")
			}
		}

		if err != nil {
			fmt.Println("error sending DM message:", err)
			s.ChannelMessageSend(
				m.ChannelID,
				"Failed to send you a DM. "+
					"Did you disable DM in your privacy settings?",
			)
		}
	}

	if m.Content == config.C.Text2 {

		// Create the private channel with the user who sent the message.
		var err error
		channel, err = s.UserChannelCreate(m.Author.ID)
		if err != nil {
			// If an error occurred, we failed to create the channel.
			fmt.Println("error creating channel:", err)
			s.ChannelMessageSend(
				m.ChannelID,
				"Something went wrong while sending the DM!",
			)
			return
		}

		// Check grade
		GradeCheckDM(channel)

		if err != nil {
			fmt.Println("error sending DM message:", err)
			s.ChannelMessageSend(
				m.ChannelID,
				"Failed to send you a DM.")
		}

	}
	if m.Content != config.C.Text1 && m.Content != config.C.Text2 {
		s.ChannelMessageSend(m.ChannelID, "Aep Du lor, pai ask cherry na 😔")
	}

}

// GradeCheck checks the grades and sends notifications
func GradeCheck(channel *discordgo.Channel, s *gocron.Scheduler) {
	//fmt.Println("checking grade 📌")
	grade, err := functions.GetGrade()
	if err != nil {
		fmt.Println(err)
		fmt.Println("error getting grade")
		return
	}

	if prevGrade == nil {
		prevGrade = grade
		return
	}

	// Add a flag to track if all courses are graded
	allGraded := true

	for _, v := range grade.GradeInfo.Grades {
		for _, v2 := range v.Courses {
			for _, v3 := range prevGrade.GradeInfo.Grades {
				for _, v4 := range v3.Courses {
					if v2.CourseCode == v4.CourseCode {
						if v2.CourseGrade != v4.CourseGrade && channel != nil {
							// Send message to DM
							_, err := dg.ChannelMessageSend(channel.ID, "Grade ook laew na 🎉\n"+
								v2.CourseCode+" "+v2.CourseNameEn+" Grade : "+v2.CourseGrade)
							if err != nil {
								fmt.Println("error sending DM message:", err)
								return
							}
						}
					}

					// Check if all courses are graded
					if v2.CourseGrade == "-" {
						allGraded = false
					}
				}
			}
		}
	}

	if allGraded && channel != nil {
		_, err := dg.ChannelMessageSend(channel.ID, "All courses are graded✨")
		if err != nil {
			fmt.Println("error sending DM message:", err)
			return
		}

		// Stop gradeCheck
		s.Stop()
	}

	prevGrade = grade

}

// GradeCheck checks the grades and sends notifications
func GradeCheckDM(channel *discordgo.Channel) {
	//fmt.Println("checking gradeDM 📌")
	grade, err := functions.GetGrade()
	if err != nil {
		fmt.Println(err)
		fmt.Println("error getting grade")
		return
	}

	if prevGrade == nil {
		prevGrade = grade
		return
	}

	// Add a flag to track if all courses are graded
	allGraded := true

	for _, v := range grade.GradeInfo.Grades {
		for _, v2 := range v.Courses {
			for _, v3 := range prevGrade.GradeInfo.Grades {
				for range v3.Courses {
					// Check if all courses are graded
					if v2.CourseGrade == "-" {
						allGraded = false
					}
				}
			}
		}
	}

	if allGraded {
		_, err := dg.ChannelMessageSend(channel.ID, "grade ook kob laew na ✨")
		if err != nil {
			fmt.Println("error sending DM message:", err)
			return
		}
	}

	if !allGraded {
		_, err := dg.ChannelMessageSend(channel.ID, "grade ook yung mai kob 😔")
		if err != nil {
			fmt.Println("error sending DM message:", err)
			return
		}
	}

	prevGrade = grade
}
