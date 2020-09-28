package main

import (
	fb "discord/internal/resources"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

// Variables used for command line parameters
var (
	Token string
	Emojis map[string]string
	Roles map[string]string
)

func init() {
	flag.StringVar(&Token, "t", os.Getenv("DISCORD_TOKEN"), "Bot Token")
	flag.Parse()

	Emojis = map[string]string{}
	Roles = map[string]string{}
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

	if val, ok := Emojis[m.Emoji.Name]; ok {
		if val == m.MessageID {
			fmt.Println(m.UserID, m.MessageReaction.UserID)
			err := s.GuildMemberRoleAdd(m.GuildID, Roles[m.MessageID], m.MessageReaction.UserID)
			fmt.Println(err)
		}
	}
	
	//guild ID, role id, memberid
	//s.GuildMemberRoleAdd()
}

func discordJoin(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
	user, _ := s.UserChannelCreate(m.User.ID)

	s.ChannelMessageSend(user.ID, "Welcome to the Lionheart beta!\n\n" +

		"Our beta program starts on **October 1st**, on a first-come, first-serve basis. Please schedule a time to onboard where we'll explain how it works and next steps.\n" +
		"After onboarding, you'll be given your login information and instructions to start the apprenticeship accelerator in the #bot-room channel.\n\n" +

		"Feel free to ask any questions in #questions and ping any of the available mods.\n" +
		"We're excited to have you be part of our accelerator! Please abide by our #code-of-conduct and #rules while at Lionheart.\n\n" +

		"See you soon!\n\n" +

		"Link to onboard: https://calendly.com/juan-lionheart/welcome-to-lionheart")
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	channel, _ := s.Channel(m.ChannelID)

	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == ".alive" {
		s.ChannelMessageSend(m.ChannelID, "Hello, yes, I'm alive, good sir.")
	}

	if m.Content == ".hello" {
		s.ChannelMessageSend(m.ChannelID, "Hello, good sir.")
	}

	if m.Content == ".db" {
		val := fb.GetNumUsers()
		num := strconv.Itoa(val)
		s.ChannelMessageSend(m.ChannelID, "There are currently " + num + " users who have taken the test.")
	}

	if strings.Contains(m.Content, ".clear") {
		s.ChannelMessageDelete(m.ChannelID, m.ID)

		values := strings.Split(m.Content, " ")

		if len(values) > 1 {
			v, _ := strconv.Atoi(values[1])
			messages, _ := s.ChannelMessages(m.ChannelID, v, "", "", "")
			for _, message := range messages {
				s.ChannelMessageDelete(m.ChannelID, message.ID)
			}
		}
	}

	if channel.Name == "feedback" {
		c, _ := s.GuildChannels(m.GuildID)
		for _, channel := range c {
			if channel.Name == "mod-feedback" {
				s.ChannelMessageDelete(m.ChannelID, m.ID)
				s.ChannelMessageSend(channel.ID, "User ID: " + m.Author.ID + "(" + m.Author.Username + ") has submitted feedback: \n" + m.Content)
				return
			}
		}
	}

	if channel.Name == "3-skill-selection-🙋" {
		s.ChannelMessageSend(m.ChannelID, "That ^ message ID is: " + m.Message.ID)
	}

	if channel.Name == "bot-room" {
		if strings.Contains(m.Content, ".verify") {
			s.ChannelMessageDelete(m.ChannelID, m.ID)
			values := strings.Split(m.Content, " ")

			if len(values) > 1 {
				user, message := fb.UserExists(values[1])

				if user {
					roles, _ := s.GuildRoles(m.GuildID)
					for _, role := range roles {
						if role.Name == "Users" {
							s.GuildMemberRoleAdd(m.GuildID, m.Author.ID, role.ID)
							s.ChannelMessageSend(m.ChannelID, m.Author.Mention()+" has been verified. Users role granted.")
						}
					}
				}
				s.ChannelMessageSend(m.ChannelID, message)
			}
		}
	}

	if strings.Contains(m.Content, ".addrole") {
		values := strings.Split(m.Content, " ")
		message := values[1]
		emoji := values[2]
		role := values[3]

		fmt.Println(message, emoji, role)

		Emojis[emoji] = message

		//.addrole <messageID> <emoji> <roleID>

		roles, _ := s.GuildRoles(m.GuildID)
		for _, r := range roles {
			fmt.Println(r)
			if strings.EqualFold(role, r.Name){
				Roles[message] = r.ID
			}
		}
	}

	if m.Content == ".done" {
		fb.WriteData("emojis", Emojis)
		fb.WriteData("roles", Roles)
	}

	//updateLeaderboards(s, m)
}

func truncatePhoneNumber(number string) {

}

func updateLeaderboards(s *discordgo.Session, m *discordgo.MessageCreate) {
	fb.WriteLeaderboards(m.Author.ID)
}