package fb

import (
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/db"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
	"os"
)

var fb *db.Client

func init() {
	if fb == nil {
		//opt := option.WithCredentialsFile("wow.json")
		opt := option.WithCredentialsJSON([]byte(os.Getenv("ACCOUNT")))
		config := &firebase.Config{
			DatabaseURL: os.Getenv("DB_URL"),
				//DatabaseURL: "https://lionheart-7b6c9.firebaseio.com/",
		}

		f, _ := firebase.NewApp(context.Background(), config, opt)
		fb, _ = f.Database(context.Background())
	}
}

func UserExists(phone string, discordID string, discordName string, nickname string) (bool, string) {
	var User User
	ref := fb.NewRef("users").Child(phone)
	err := ref.Get(context.Background(), &User)
	if err != nil {
		return false, "Error: searching user has failed"
	}

	if User.PersonalInfo.Name != "" {
		if User.Verified {
			return false, "Someone has already been verified with that phone number."
		} else {
			ref.Child("Verified").Set(context.Background(), true)
			ref.Child("Discord/ID").Set(context.Background(), discordID)
			ref.Child("Discord/Username").Set(context.Background(), discordName)
			ref.Child("Discord/Nickname").Set(context.Background(), nickname)
			fb.NewRef("lookup").Child(discordID).Set(context.Background(), phone)
			return true, ""
		}
	}

	return false, "The phone number cannot be found. Please check to make sure you typed in the correct phone number."
}


func WriteUser(phone string, user User) {
	fb.NewRef("users").Child(phone).Set(context.Background(), user)
}

func GetNumUsers() int {
	v := map[string]User{}
	ref := fb.NewRef("users")
	err := ref.Get(context.Background(), &v)
	if err != nil {
		fmt.Println(err)
	}

	return len(v)
}

func GetUsers() map[string]User {
	v := map[string]User{}
	ref := fb.NewRef("users")
	err := ref.Get(context.Background(), &v)
	if err != nil {
		fmt.Println(err)
	}

	return v
}

func WriteLeaderboards(user string) {
	boardUser := LeaderboardUser{}
	ref := fb.NewRef("leaderboards").Child(user)
	err := ref.Get(context.Background(), &boardUser)

	if err != nil {
		fmt.Println(err)
	}

	boardUser.Points += 1
	boardUser.Messages += 1

	ref.Set(context.Background(), boardUser)
}

func WriteData(child string, data map[string]Emoji) {
	ref := fb.NewRef(child)
	ref.Set(context.Background(), data)
}

func LoadData() map[string]Emoji {
	v := map[string]Emoji{}
	ref := fb.NewRef("emojis")
	err := ref.Get(context.Background(), &v)
	if err != nil {
		fmt.Println(err)
	}

	return v
}

func GetUserByNumber(phone string) User {
	var User User
	ref := fb.NewRef("users").Child(phone)
	err := ref.Get(context.Background(), &User)
	if err != nil {
		fmt.Println(err)
	}

	return User
}

func GetUserPhoneNumber(ID string) string {
	var number string
	ref := fb.NewRef("lookup").Child(ID)
	err := ref.Get(context.Background(), &number)
	if err != nil {
		fmt.Println(err)
	}

	return number
}

func WritePods(pods map[string][]Pod) {
	fb.NewRef("pods").Set(context.Background(), pods)
}

func DeleteChild(child string) {
	fb.NewRef(child).Delete(context.Background())
}