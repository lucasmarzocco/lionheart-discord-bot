package fb

type Personal struct {
	Name      string
	Email     string
	Phone     string
	Gender    string
	Ethnicity string
	Education string
	CountryFrom   string
	CountryNow      string
	State     string
	City      string
	Religion  string
	Marital   string
}

type Discord struct {
	ID       string
	Username string
	Nick     string
}

type User struct {
	PersonalInfo Personal
	Subtraits    map[string]*Trait
	Traits       map[string]*Trait
	Verified     bool
	Roles        []string
	Discord      Discord
}

type Trait struct {
	Name        string
	RawScore    float64
	NormalScore float64
	Min         float64
}

type LeaderboardUser struct {
	Points int64
	Messages int64
}

type Emoji struct {
	MessageID string
	RoleID    string
}

type PodMember struct {
	DiscordInfo Discord
	PhoneNumber string
}

type Pod struct {
	Limit int
	Size  int
	Members []PodMember
	Skill string
	RoomName string
}

func (p PodMember) String() string {
	return p.DiscordInfo.Nick
}