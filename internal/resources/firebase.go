package fb

import (
	"fmt"
	"os"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/db"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
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

func UserExists(phone string) (bool, string) {
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
			ref.Child("verified").Set(context.Background(), true)
			return true, ""
		}
	}

	return false, "The phone number cannot be found. Please check to make sure you typed in the correct phone number."
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