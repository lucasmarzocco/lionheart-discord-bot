package main

import (
	fb "discord/internal/resources"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/robfig/cron/v3"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// Variables used for command line parameters
var (
	Token   string
	Emojis  map[string]fb.Emoji
	BotRoom string
	Pods    string
	Errors string
	Session *discordgo.Session
	Roles   map[string]string
	Quotes string
	Commands map[string]string
)

func init() {
	//Emojis = fb.LoadData()
	BotRoom = "765438490007175179"
	Pods = "765514925531070502"
	Errors = "768057592115101746"
	Quotes = "803545216464322570"

	Roles = map[string]string{
		"Adventurousness":      "Perception",
		"Artistic Interest":    "Perception",
		"Emotionality":         "Perception",
		"Imagination":          "Perception",
		"Intellect":            "Perception",
		"Liberalism":           "Perception",
		"Cautiousness":         "Productivity",
		"Dutifulness":          "Productivity",
		"Achievement Striving": "Productivity",
		"Orderliness":          "Productivity",
		"Self-Discipline":      "Productivity",
		"Self-Efficacy":        "Productivity",
		"Activity":             "Fulfillment",
		"Assertiveness":        "Fulfillment",
		"Cheerfulness":         "Fulfillment",
		"Excitement Seeking":   "Fulfillment",
		"Friendliness":         "Fulfillment",
		"Gregariousness":       "Fulfillment",
		"Altruism":             "Collaboration",
		"Cooperation":          "Collaboration",
		"Modesty":              "Collaboration",
		"Morality":             "Collaboration",
		"Sympathy":             "Collaboration",
		"Trust":                "Collaboration",
		"Peacefulness":         "Reslience",
		"Calmness":             "Reslience",
		"Lionheartedness":      "Reslience",
		"Moderation":           "Reslience",
		"Self-Confidence":      "Reslience",
		"Invulnerability":      "Reslience",
	}

	Commands = map[string]string {
		".alive" : "Check to see if the bot is functional",
		".id"    : "Check the ID of the current channel",
		".db"    :  "See details about the database",
		".clear <#>" : "Clear # of messages in the channel",
		".createrole apprentice_NameHere" : "Create a role and channel named NameHere",
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	token := os.Getenv("DISCORD_TOKEN")

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + token)

	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	Session = dg

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)
	//dg.AddHandler(discordJoin)
	//dg.AddHandler(messageReactAdd)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)

	randomQuotes()

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

func randomQuotes() {

	fmt.Println("Create new cron")
	c := cron.New()

	c.AddFunc("0 13 * * 0,1,3,5", quotes)

	fmt.Println("Start cron")
	c.Start()
}

func quotes() {
	Session.ChannelMessageSend(Quotes, GetQuote())
}

func messageReactAdd(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
	if val, ok := Emojis[m.Emoji.Name]; ok {
		if val.MessageID == m.MessageID {
			member, _ := s.GuildMember(m.GuildID, m.UserID)
			roleFlag := false
			for _, role := range member.Roles {
				if strings.EqualFold(role, "756572611516825600") {
					roleFlag = true
				}
			}

			if roleFlag {
				if len(member.Roles) >= 3 {
					s.MessageReactionRemove(m.ChannelID, m.MessageID, m.Emoji.Name, m.UserID)
					user, _ := s.UserChannelCreate(m.UserID)
					s.ChannelMessageSend(user.ID, "Currently users can only have 2 skills to level up at once. If you chose a skill on accident, please post in #questions. Thanks!")
					return
				}

				phone := fb.GetUserPhoneNumber(m.UserID)
				user := fb.GetUserByNumber(phone)
				roles, _ := s.GuildRoles(m.GuildID)
				roleName := ""

				for _, role := range roles {
					if role.ID == val.RoleID {
						if !isRoleThere(role.Name, user.Roles) {
							user.Roles = append(user.Roles, role.Name)
							roleName = role.Name
							break
						} else {
							return
						}
					}
				}

				fb.WriteUser(phone, user)
				s.GuildMemberEdit(m.GuildID, m.UserID, append(member.Roles, val.RoleID))
				s.ChannelMessageSend(BotRoom, member.User.Username + "#" + member.User.Discriminator + " added role: "+roleName+" on: "+time.Now().Format("2006-01-02-15:04:05"))

			} else {
				s.MessageReactionRemove(m.ChannelID, m.MessageID, m.Emoji.Name, m.UserID)
				user, _ := s.UserChannelCreate(m.UserID)
				s.ChannelMessageSend(user.ID, "You have not verified. Please verify in order to select your skills!")
			}
		}
	}
}

func discordJoin(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
	user, _ := s.UserChannelCreate(m.User.ID)

	s.ChannelMessageSend(user.ID, "Welcome to Lionheart! \n\n "+

		"#rules - Explains our rules and code of conduct while in Lionheart. \n "+
		"#questions - Got questions? Ask here! \n\n "+

		"#1-start-here - An overview of Lionheart \n "+
		"#2-verify - Verify your phone number after you've taken the skills assessment. No bots please. \n "+
		"#3-skill-selection - After you've verified your phone number, select the skill to level and get grouped with a pod. \n\n "+

		"Link to skills assessment: https://join.lionheart-institution.app/latent-potential \n "+
		"If you need personalized assistance, schedule a time to chat with Juan: \n "+
		"https://calendly.com/juan-lionheart/welcome-to-lionheart \n\n"+

		"Thanks for joining. Reach out to any of the Admins or Mods if you need any assistance during your journey.")


	roles, _ := s.GuildRoles(m.GuildID)
	for _, role := range roles {
		if role.Name == "Guests" {
			fmt.Println("Adding guest role to new user")
			s.GuildMemberRoleAdd(m.GuildID, m.User.ID, role.ID)
		}
	}
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	channel, _ := s.Channel(m.ChannelID)
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == ".test" {
		roles, _ := s.GuildRoles(m.GuildID)
		for _, role := range roles {
			if val, ok := Roles[role.Name]; ok {
				s.ChannelMessageSend(m.ChannelID, "Role: " + role.Name + "found: " + val)
			}
		}
	}

	if m.Content == ".alive" {
		s.ChannelMessageSend(m.ChannelID, "Hello, yes, I'm alive, good sir.")
	}

	if m.Content == ".clean" {
		rooms, _ := s.GuildChannels(m.GuildID)
		roles, _ := s.GuildRoles(m.GuildID)

		for _, room := range rooms {
			if strings.Contains(room.Name, "_") {
				s.ChannelDelete(room.ID)
			}
		}

		for _, role := range roles {
			if strings.Contains(role.Name, "_") {
				s.GuildRoleDelete(m.GuildID, role.ID)
			}
		}

		fb.DeleteChild("pods")
	}

	if m.Content == ".id" {
		s.ChannelMessageSend(m.ChannelID, "Channel ID is: "+m.ChannelID)
	}

	if strings.Contains(m.Content, ".embed") {
		s.ChannelMessageDelete(m.ChannelID, m.ID)
		messageEmbed := &discordgo.MessageEmbed{}
		values := strings.Split(m.Content, "-")

		messageEmbed.Title = values[1]
		messageEmbed.URL = values[2]

		messageEmbed.Fields = []*discordgo.MessageEmbedField{
			{
				Name: values[3],
				Value: values[4],
			},
		}

		message := &discordgo.MessageSend {
			Embed: messageEmbed,
		}
		_, err := s.ChannelMessageSendComplex(m.ChannelID, message)

		if err != nil {
			fmt.Println(err)
		}
	}

	if m.Content == ".db" {
		fb.GetUsers()
		val := fb.GetNumUsers()
		num := strconv.Itoa(val)
		s.ChannelMessageSend(m.ChannelID, "There are currently "+num+" users who have taken the test.")
	}

	/*if m.Content == ".help" {
		help := ""
		for key, val := range Commands {
			help += key + " => " + val + "\n"
		}
		s.ChannelMessageSend(m.ChannelID, help)
	} */

	if strings.Contains(m.Content, ".createrole") {
		values := strings.Split(m.Content, " ")
		fmt.Println(values)

		if strings.Contains(values[1], "apprentice") {

			values := strings.Split(values[1], "_")
			fmt.Println(values)
			roleName := values[1]

			role, _ := s.GuildRoleCreate(m.GuildID)
			s.GuildRoleEdit(m.GuildID, role.ID, roleName, 0, false, 0, false)
			createDiscordRoom(m.GuildID, roleName, role.ID)
		}
	}

	if strings.Contains(m.Content, ".clear") {
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
			if channel.Name == "internal-feedback" {
				s.ChannelMessageDelete(m.ChannelID, m.ID)
				s.ChannelMessageSend(channel.ID, "User ID: "+m.Author.ID+" ("+m.Author.String()+") has submitted feedback: \n"+m.Content)
				return
			}
		}
	}

	if channel.Name == "3-bot-room" {
		if strings.Contains(m.Content, ".verify") {
			values := strings.Split(m.Content, " ")

			if len(values) > 1 {
				mem, _ := s.GuildMember(m.GuildID, m.Author.ID)
				user, message := fb.UserExists("1"+fixPhoneNumber(values[1]), m.Author.ID, m.Author.Username, mem.Nick)

				if user {
					roles, _ := s.GuildRoles(m.GuildID)
					for _, role := range roles {

						if role.Name == "Guests" {
							s.GuildMemberRoleRemove(m.GuildID, m.Author.ID, role.ID)
						}
						if role.Name == "Apprentices" {
							s.GuildMemberRoleAdd(m.GuildID, m.Author.ID, role.ID)
							s.ChannelMessageSend(m.ChannelID, m.Author.Mention()+" has been verified. Apprentices role granted. Please go to #3-skill-selection to continue.")
						}
					}
				}
				s.ChannelMessageSend(m.ChannelID, message)
				return
			}
		}
	}

	if channel.Name == "4-skill-selection" {
		fmt.Println("That ^ message ID is: " + m.Message.ID)
	}
}

func fixPhoneNumber(number string) string {
	reg, err := regexp.Compile("[^0-9]+")
	if err != nil {
		fmt.Println(err)
	}

	return reg.ReplaceAllString(number, "")
}


func createDiscordRoom(guildID string, roomName string, roleID string) {
	var everyoneID string
	var adminID string

	roles, _ := Session.GuildRoles(guildID)
	for _, role := range roles {
		if role.Name == "@everyone" {
			everyoneID = role.ID
		}
		if role.Name == "Admin" {
			adminID = role.ID
		}
	}

	data := discordgo.GuildChannelCreateData{
		Name: roomName,
		Type: discordgo.ChannelTypeGuildText,
		PermissionOverwrites: []*discordgo.PermissionOverwrite{
			{
				ID:    adminID,
				Type:  "role",
				Allow: 8,
			},
			{
				ID:    roleID,
				Type:  "role",
				Allow: 523328,
			},
			{
				ID:   everyoneID,
				Type: "role",
				Deny: 523328,
			},
		},
		ParentID: "767315418687078411",
	}

	_, err := Session.GuildChannelCreateComplex(guildID, data)
	if err != nil {
		Session.ChannelMessageSend(Errors, "Failed to create channel with room name: " + roomName)
		fmt.Println(err)
	}
}

func isRoleThere(role string, roles []string) bool {
	for _, r := range roles {
		if r == role {
			return true
		}
	}

	return false
}



/*

FEATURES TO BE ADDED EVENTUALLY

 */
func updateLeaderboards(s *discordgo.Session, m *discordgo.MessageCreate) {
	fb.WriteLeaderboards(m.Author.ID)
}
