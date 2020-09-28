package fb

type Personal struct {
	Name string
	Email string
	Phone string
	MorF  string
	Gender string
	Ethnicity string
	Education string
	Country string
	USA     bool
	State	string
	City    string
	Live    string
	Religion string
	Marital string
}

type User struct {
	PersonalInfo Personal
	Subtraits    map[string]*Trait
	Traits       map[string]*Trait
	Verified     bool
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