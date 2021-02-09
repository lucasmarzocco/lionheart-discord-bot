package main

import (
	fb "discord/internal/resources"
	"flag"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/robfig/cron/v3"
	"net/http"
	"net/url"
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
)

func init() {
	flag.StringVar(&Token, "t", os.Getenv("DISCORD_TOKEN"), "Bot Token")
	flag.Parse()

	Emojis = fb.LoadData()
	BotRoom = "765438490007175179"
	Pods = "765514925531070502"
	Errors = "768057592115101746"

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

	Session = dg

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)
	dg.AddHandler(discordJoin)
	dg.AddHandler(messageReactAdd)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)

	//startMatching()

	randomQuotes()

	happyBirthday()

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

	c.AddFunc("0 14 * * *", quotes)

	// Start cron with one scheduled job
	fmt.Println("Start cron")
	c.Start()
}

func quotes() {
	channelID := "803545216464322570"
	Session.ChannelMessageSend(channelID, GetQuote())
}

func text() {

	now := time.Now()
	newLayout := "15:04"
	ns, _ := time.Parse(newLayout, strconv.Itoa(now.Hour())+":"+strconv.Itoa(now.Minute()))
	srt, _ := time.Parse(newLayout, "09:20")

	if ns.After(srt) {
		accountSid := os.Getenv("ACCOUNT_SID")
		token := os.Getenv("TOKEN")
		urlStr := "https://api.twilio.com/2010-04-01/Accounts/" + accountSid + "/Messages.json"

		msgData := url.Values{}
		msgData.Set("To", "9254467645")
		msgData.Set("From", os.Getenv("PHONE"))
		msgData.Set("Body", "HEHE UR CUTE!")
		msgDataReader := *strings.NewReader(msgData.Encode())

		client := &http.Client{}
		req, err := http.NewRequest("POST", urlStr, &msgDataReader)
		if err != nil {
			panic(err)
		}

		req.SetBasicAuth(accountSid, token)
		req.Header.Add("Accept", "application/json")
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		client.Do(req)
	}


}

func happyBirthday() {
	fmt.Println("Create new cron")
	c := cron.New()

	c.AddFunc("*/21 * * * *", text)

	// Start cron with one scheduled job
	fmt.Println("Start cron")
	c.Start()
}

func startMatching() {
	fmt.Println("Create new cron")
	c := cron.New()

	//should be 4am everyday, make that project level, ENV variables
	c.AddFunc("0 11 * * *", matchUsers)

	// Start cron with one scheduled job
	fmt.Println("Start cron")
	c.Start()
}

func matchUsers() {
	Session.ChannelMessageSend(Pods, "Group matching initiated.")

	guildID := "666130712734466051"
	pods := make(map[string][]fb.Pod)
	users := fb.GetUsers()

	for phone, user := range users {
		if !user.Verified {
			continue
		}

		for _, skill := range user.Roles {
			bigFive := Roles[skill]

			if val, ok := pods[bigFive]; ok {
				pod := val[len(val)-1]
				if pod.Size == pod.Limit {
					newPod := fb.Pod{
						Limit: 5,
						Size:  1,
						Members: []fb.PodMember{
							{
								DiscordInfo: fb.Discord{
									user.Discord.ID,
									user.Discord.Username,
									user.Discord.Nickname,
								},
								PhoneNumber: phone,
							}},
						Skill:    bigFive,
						RoomName: bigFive + "_" + string(len(val)),
					}
					pods[bigFive] = append(pods[bigFive], newPod)

				} else {
					pod.Members = append(pod.Members, fb.PodMember{
						DiscordInfo: fb.Discord{
							user.Discord.ID,
							user.Discord.Username,
							user.Discord.Nickname,
						},
						PhoneNumber: phone,
					})

					pod.Size += 1
					pods[bigFive][len(pods[bigFive])-1] = pod
				}
			} else {
				pods[bigFive] = []fb.Pod{
					{
						Limit: 5,
						Size:  1,
						Members: []fb.PodMember{
							{
								DiscordInfo: fb.Discord{
									user.Discord.ID,
									user.Discord.Username,
									user.Discord.Nickname,
								},
								PhoneNumber: phone,
							}},
						Skill:    bigFive,
						RoomName: bigFive + "_1",
					},
				}
			}
		}
	}

	Session.ChannelMessageSend(Pods, prettyPrintPods(pods))
	fb.WritePods(pods)
	putInRoomsAndSetRoles(pods, guildID)
	Session.ChannelMessageSend(Pods, "Group matching complete")
}

func putInRoomsAndSetRoles(pods map[string][]fb.Pod, guildID string) {
	for _, pod := range pods {
		for _, ele := range pod {
			role, _ := Session.GuildRoleCreate(guildID)
			Session.GuildRoleEdit(guildID, role.ID, ele.RoomName, 0, false, 0, false)
			createDiscordRoom(guildID, ele.RoomName, role.ID)

			for _, member := range ele.Members {
				m, _ := Session.GuildMember(guildID, member.DiscordInfo.ID)
				Session.GuildMemberEdit(guildID, member.DiscordInfo.ID, append(m.Roles, role.ID))
			}
		}
	}
}

func prettyPrintPods(pods map[string][]fb.Pod) string {
	str := ""

	for skill, pod := range pods {
		for i, p := range pod {
			str +=
				"Skill: " + skill + "- Pod: " + strconv.Itoa(i+1) + "\n" +
					"# Members in this pod: " + strconv.Itoa(p.Size) + "\n" +
					"Who? " + fmt.Sprint(p.Members) + "\n" +
					"------------" + "\n"
		}
	}
	return str
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
		s.ChannelMessageDelete(m.ChannelID, m.ID)

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

	if m.Content == ".create" {
		role, _ := s.GuildRoleCreate(m.GuildID)

		s.GuildRoleEdit(m.GuildID, role.ID, "Test-123", 0, false, 0, false)
		createDiscordRoom(m.GuildID, "TEST-HERE", role.ID)
	}

	if m.Content == ".match" {
		s.ChannelMessageDelete(m.ChannelID, m.ID)
		matchUsers()
	}

	if m.Content == ".db" {
		fb.GetUsers()
		val := fb.GetNumUsers()
		num := strconv.Itoa(val)
		s.ChannelMessageSend(m.ChannelID, "There are currently "+num+" users who have taken the test.")
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
				s.ChannelMessageSend(channel.ID, "User ID: "+m.Author.ID+" ("+m.Author.String()+") has submitted feedback: \n"+m.Content)
				return
			}
		}
	}

	if channel.Name == "4-skill-selection" {
		fmt.Println("That ^ message ID is: " + m.Message.ID)
	}

	if channel.Name == "3-bot-room" {
		if strings.Contains(m.Content, ".verify") {
			s.ChannelMessageDelete(m.ChannelID, m.ID)
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
						if role.Name == "Levelers" {
							s.GuildMemberRoleAdd(m.GuildID, m.Author.ID, role.ID)
							s.ChannelMessageSend(m.ChannelID, m.Author.Mention()+" has been verified. Levelers role granted. Please go to #3-skill-selection to continue.")
						}
					}
				}
				s.ChannelMessageSend(m.ChannelID, message)
				return
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
			if strings.EqualFold(role, r.Name) {
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
