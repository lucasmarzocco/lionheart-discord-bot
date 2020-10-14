package main

import (
	fb "discord/internal/resources"
	"flag"
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
	Token string
	Emojis map[string]fb.Emoji
	BotRoom string
	Pods string
	Session *discordgo.Session
)

func init() {
	flag.StringVar(&Token, "t", os.Getenv("DISCORD_TOKEN"), "Bot Token")
	flag.Parse()

	Emojis = fb.LoadData()
	BotRoom = "765438490007175179"
	Pods = "765514925531070502"
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

	startMatching()

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

func startMatching() {
	fmt.Println("Create new cron")
	c := cron.New()

	//should be 4am everyday, make that project level, ENV variables
	c.AddFunc("45 11 * * *", matchUsers)

	// Start cron with one scheduled job
	fmt.Println("Start cron")
	c.Start()
}

func TestPods() map[string][]fb.Pod {

	pods := make(map[string][]fb.Pod)

	pods["Adventurousness"] = []fb.Pod {
		{
			Limit: 5,
			Size: 5,
			Members: []fb.PodMember{
				{
					DiscordInfo: fb.Discord{
						"132332",
						"Sally",
						"Sally",
					},
					PhoneNumber: "2393284294",
				},
				{
					DiscordInfo: fb.Discord{
						"132332",
						"Bob",
						"Bob",
					},
					PhoneNumber: "2393284294",
				},
				{
					DiscordInfo: fb.Discord{
						"132332",
						"Jake",
						"Jake",
					},
					PhoneNumber: "2393284294",
				},
				{
					DiscordInfo: fb.Discord{
						"132332",
						"Martin",
						"Martin",
					},
					PhoneNumber: "2393284294",
				},
				{
					DiscordInfo: fb.Discord{
						"132332",
						"Juan",
						"Juan",
					},
					PhoneNumber: "2393284294",
				},


			},
			Skill: "Adventurousness",
			RoomName: "Adventurousness_1",
		},
	}

	pods["Liberalism"] = []fb.Pod {
		{
			Limit: 5,
			Size: 4,
			Members: []fb.PodMember{
				{
					DiscordInfo: fb.Discord{
						"132332",
						"Sally",
						"Sally",
					},
					PhoneNumber: "2393284294",
				},
				{
					DiscordInfo: fb.Discord{
						"132332",
						"Jason",
						"Jason",
					},
					PhoneNumber: "2393284294",
				},
				{
					DiscordInfo: fb.Discord{
						"132332",
						"Duke",
						"Duke",
					},
					PhoneNumber: "2393284294",
				},
				{
					DiscordInfo: fb.Discord{
						"132332",
						"Bob",
						"Bob",
					},
					PhoneNumber: "2393284294",
				},
			},
			Skill: "Liberalism",
			RoomName: "Liberalism_1",
		},
	}


	return pods
}

func matchUsers() {
	fmt.Println("MATCHING USERS NOW====================")
	Session.ChannelMessageSend(Pods, "Group matching begun, stay tuned...")

	pods := make(map[string][]fb.Pod)
	users := fb.GetUsers()

	for phone, user := range users {
		fmt.Println(user)
		if !user.Verified {
			continue
		}

		for _, skill := range user.Roles {
			fmt.Println("CHECKING SKILL: " + skill)

			if val, ok := pods[skill]; ok {

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
						Skill:    skill,
						RoomName: skill + "_" + string(len(val)),
					}
					pods[skill] = append(pods[skill], newPod)

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
					pods[skill][len(pods[skill])-1] = pod
				}
			} else {
				pods[skill] = []fb.Pod {
					{
						Limit: 5,
						Size: 1,
						Members: []fb.PodMember{
							{
								DiscordInfo: fb.Discord{
									user.Discord.ID,
									user.Discord.Username,
									user.Discord.Nickname,
								},
								PhoneNumber: phone,
							}},
						Skill: skill,
						RoomName: skill + "_1",
					},
				}
			}
		}
	}

	Session.ChannelMessageSend(Pods, "Group matching over, this is what I found...")
	Session.ChannelMessageSend(Pods, prettyPrintPods(pods))
	fb.WritePods(pods)
}

func prettyPrintPods(pods map[string][]fb.Pod) string {
	str := ""

	for skill, pod := range pods {
		for i, p := range pod {
			str +=
				"Skill: " + skill + "- Pod: " + strconv.Itoa(i + 1) + "\n" +
				"# Members in this pod: " + strconv.Itoa(p.Size) + "\n" +
				"Who? " + fmt.Sprint(p.Members) + "\n" +
				"------------" + "\n"
		}
	}
	return str
}

func messageReactAdd(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
	fmt.Println("There was a message react.")

	if val, ok := Emojis[m.Emoji.Name]; ok {
		if val.MessageID == m.MessageID {
			member, _ := s.GuildMember(m.GuildID, m.UserID)

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
					user.Roles = append(user.Roles, role.Name)
					roleName = role.Name
					break
				}
			}

			fb.WriteUser(phone, user)
			s.GuildMemberEdit(m.GuildID, m.UserID, append(member.Roles, val.RoleID))
			s.ChannelMessageSend(BotRoom, "DEBUG: " + member.Nick + " added role: " + roleName + " on: " + time.Now().Format("2006-01-02-15:04:05"))
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

	if m.Content == ".id" {
		s.ChannelMessageSend(m.ChannelID, "Channel ID is: " + m.ChannelID)
	}

	if m.Content == ".db" {
		fb.GetUsers()
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
				s.ChannelMessageSend(channel.ID, "User ID: " + m.Author.ID + " (" + m.Author.String() + ") has submitted feedback: \n" + m.Content)
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
				mem, _ := s.GuildMember(m.GuildID, m.Author.ID)
				user, message := fb.UserExists("1" + fixPhoneNumber(values[1]), m.Author.ID, m.Author.Username, mem.Nick)

				if user {
					roles, _ := s.GuildRoles(m.GuildID)
					for _, role := range roles {
						if role.Name == "Levelers" {
							s.GuildMemberRoleAdd(m.GuildID, m.Author.ID, role.ID)
							s.ChannelMessageSend(m.ChannelID, m.Author.Mention()+ " has been verified. Levelers role granted.")
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