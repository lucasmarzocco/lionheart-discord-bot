package main

import (
	fb "discord/internal/resources"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

// Variables used for command line parameters
var (
	Token string
)

func init() {
	flag.StringVar(&Token, "t", os.Getenv("DISCORD_TOKEN"), "Bot Token")
	flag.Parse()
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)

	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)
	dg.AddHandler(discordJoin)
	dg.AddHandler(messageReact)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

func messageReact(s *discordgo.Session, m *discordgo.MessageReactionAdd) {

	fmt.Println("There was a message react....")

	fmt.Println(m.UserID)
	fmt.Println(m.ChannelID)
	fmt.Println(m.GuildID)
	fmt.Println(m.Emoji.Name)
	fmt.Println(m.MessageID)

	a, _ := s.User(m.UserID)
	fmt.Println(a.Username)
}

func discordJoin(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
	user, _ := s.UserChannelCreate(m.User.ID)
	roles, _ := s.GuildRoles(m.GuildID)
	for _, role := range roles {
		if role.Name == "Guests" {
			s.GuildMemberRoleAdd(m.GuildID, m.User.ID, role.ID)
		}
	}

	s.ChannelMessageSend(user.ID, "Welcome! Thank you for joining the Lionheart community. \n" +
		"To begin the journey of becoming a Member, please onboard at this link:\n" +
		"https://calendly.com/juan-lionheart/welcome-to-lionheart\n\n" +

		"Please schedule a time to talk to Juan, the Leader of Lionheart. He will explain a little more about " +
		"what we're trying to accomplish here and what the assessment (your next step) is trying to achieve.\n\n" +

		"When you're done with the assessment, please visit the #bot-room channel in the Lionheart Discord.\n" +
		"To obtain the member role, I need to verify that you've completed the assessment. To do so,\n" +
		"please type the command: *.verify <phone number>* where <phone number> is the phone number you\n" +
		"entered at the start of the test. \n\n" +

		"Form: .verify 1XXXXXXXXXX \n" +
		"Example: .verify 11234567890 \n\n" +

		"Thank you for joining us. We hope you enjoy your time here. Please reach out to any Admins or Moderators " +
		"if you need any assistance.")
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}
	// If the message is "ping" reply with "Pong!"
	if m.Content == ".alive" {
		s.ChannelMessageSend(m.ChannelID, "Hello, yes, I'm alive, good sir.")
	}

	if m.Content == ".hello" {
		s.ChannelMessageSend(m.ChannelID, "Hello, good sir.")
	}

	if m.Content == ".clear" {
		s.ChannelMessageDelete(m.ChannelID, m.ID)
		
		channel, _ := s.Channel(m.ChannelID)
		messages := channel.Messages

		for _, message := range messages {
			s.ChannelMessageDelete(m.ChannelID, message.ID)
		}
	}

	channel, _ := s.Channel(m.ChannelID)
	if channel.Name == "bot-room" {
		if strings.Contains(m.Content, ".verify") {
			s.ChannelMessageDelete(m.ChannelID, m.ID)
			values := strings.Split(m.Content, " ")

			if len(values) > 1 {
				user := fb.UserExists(values[1])

				if user {
					roles, _ := s.GuildRoles(m.GuildID)
					for _, role := range roles {
						if role.Name == "Guests" {
							s.GuildMemberRoleRemove(m.GuildID, m.Author.ID, role.ID)
						}
						if role.Name == "Members" {
							s.GuildMemberRoleAdd(m.GuildID, m.Author.ID, role.ID)
							s.ChannelMessageSend(m.ChannelID, m.Author.Mention()+" has been verified. Member role granted.")
						}
					}
				}
			}
		}
	}
}
