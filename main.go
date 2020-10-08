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
	"regexp"

	"github.com/bwmarrin/discordgo"
)

// Variables used for command line parameters
var (
	Token string
	Emojis map[string]fb.Emoji
)

func init() {
	flag.StringVar(&Token, "t", os.Getenv("DISCORD_TOKEN"), "Bot Token")
	flag.Parse()

	Emojis = fb.LoadData()
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
	dg.AddHandler(messageReactAdd)

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




func messageReactAdd(s *discordgo.Session, m *discordgo.MessageReactionAdd) {

	fmt.Println("There was a message react....")

	if val, ok := Emojis[m.Emoji.Name]; ok {
		if val.MessageID == m.MessageID {
			member, _ := s.GuildMember(m.GuildID, m.UserID)

			if len(member.Roles) == 3 {
				user, _ := s.UserChannelCreate(m.UserID)
				s.ChannelMessageSend(user.ID, "Sorry! Currently you can only have 2 categories. If this was a mistake, please ask in #questions.")
				fmt.Println(s.MessageReactionRemove(m.ChannelID, m.MessageID, m.Emoji.APIName(), user.ID))
				return
			}

			err := s.GuildMemberEdit(m.GuildID, m.UserID, append(member.Roles, val.RoleID))
			fmt.Println(err)
			fmt.Println("Message react complete...")
		}
	}
}

func discordJoin(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
	user, _ := s.UserChannelCreate(m.User.ID)

	s.ChannelMessageSend(user.ID, "Welcome to Lionheart! \n\n " +

		"#rules - Explains our rules and code of conduct while in Lionheart. \n " +
		"#questions - Got questions? Ask here! \n\n " +

		"#1-start-here - An overview of Lionheart \n " +
		"#2-verify - Verify your phone number after you've taken the skills assessment. No bots please. \n " +
		"#3-skill-selection - After you've verified your phone number, select the skill to level and get grouped with a pod. \n\n " +

		"Link to skills assessment: https://join.lionheart-institution.app/latent-potential \n " +
		"If you need personalized assistance, schedule a time to chat with Juan: \n " +
		"https://calendly.com/juan-lionheart/welcome-to-lionheart \n\n" +

		"Thanks for joining. Reach out to any of the Admins or Mods if you need any assistance during your journey.")
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

	if channel.Name == "anon-feedback" {
		c, _ := s.GuildChannels(m.GuildID)
		for _, channel := range c {
			if channel.Name == "mod-feedback" {
				s.ChannelMessageDelete(m.ChannelID, m.ID)
				s.ChannelMessageSend(channel.ID, "User ID: " + m.Author.ID + " (" + m.Author.Username+ ") has submitted feedback: \n" + m.Content)
				return
			}
		}
	}

	if channel.Name == "3-skill-selection-🙋" {
		fmt.Println("That ^ message ID is: " + m.Message.ID)
	}

	if channel.Name == "bot-room" {
		if strings.Contains(m.Content, ".verify") {
			s.ChannelMessageDelete(m.ChannelID, m.ID)
			values := strings.Split(m.Content, " ")

			if len(values) > 1 {
				user, message := fb.UserExists("1" + fixPhoneNumber(values[1]))

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
		values := strings.Split(m.Content, "|")
		message := values[1]
		emoji := values[2]
		role := values[3]
		e := fb.Emoji{}
		e.MessageID = message

		roles, _ := s.GuildRoles(m.GuildID)
		for _, r := range roles {
			fmt.Println(r)
			if strings.EqualFold(role, r.Name){
				e.RoleID = r.ID
			}
		}

		Emojis[emoji] = e
	}

	if m.Content == ".done" {
		fb.WriteData("emojis", Emojis)
	}
}

func fixPhoneNumber(number string) string {
	reg, err := regexp.Compile("[^0-9]+")
	if err != nil {
		fmt.Println(err)
	}

	return reg.ReplaceAllString(number, "")
}

func updateLeaderboards(s *discordgo.Session, m *discordgo.MessageCreate) {
	fb.WriteLeaderboards(m.Author.ID)
}