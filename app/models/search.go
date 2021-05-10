package models

type SearchInput struct {
	UserQuery       *UserSearchQuery
	TeamQuery       *TeamSearchQuery
	TournamentQuery *TournamentSearchQuery
}

type SearchOutput struct {
	Teams       *Teams       `json:"teams,omitempty"`
	Tournaments *Tournaments `json:"tournaments,omitempty"`
	Users       *Users       `json:"users,omitempty"`
}

type SearchQueryBase struct {
	Text   string
	Offset uint
}

type UserSearchQuery struct {
	Base *SearchQueryBase
	// Roles int // заменить на битовую маску после поддержания ролей
}

type TeamSearchQuery struct {
	Base        *SearchQueryBase
	MemberCount uint
}

type TournamentSearchQuery struct {
	Base *SearchQueryBase
	// Location string
	KindOfSport string
}
