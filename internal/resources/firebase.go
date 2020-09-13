package fb

import (
	"os"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/db"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
)

var fb *db.Client

func init() {
	if fb == nil {
		opt := option.WithCredentialsFile("internal/resources/account.json")
		config := &firebase.Config{
			DatabaseURL: os.Getenv("DB_URL"),
		}

		f, _ := firebase.NewApp(context.Background(), config, opt)
		fb, _ = f.Database(context.Background())
	}
}

func UserExists(phone string) bool {
	var User User
	ref := fb.NewRef("users").Child(phone)
	err := ref.Get(context.Background(), &User)
	if err != nil {
		return false
	}

	if User.PersonalInfo.Name != "" {
		if User.Verified {
			return false
		} else {
			ref.Child("verified").Set(context.Background(), true)
			return true
		}
	}

	return false
}