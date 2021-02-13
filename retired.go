package main

import (
	fb "discord/internal/resources"
	"strconv"
	"fmt"
)

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

/*
if m.Content == ".create" {
role, _ := s.GuildRoleCreate(m.GuildID)

s.GuildRoleEdit(m.GuildID, role.ID, "Test-123", 0, false, 0, false)
createDiscordRoom(m.GuildID, "TEST-HERE", role.ID)
}

if m.Content == ".match" {
s.ChannelMessageDelete(m.ChannelID, m.ID)
matchUsers()
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














 */