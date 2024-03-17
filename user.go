package vk_test_task

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type ChangeDataFilm struct {
	Title       string   `json:"title,omitempty"`
	Description string   `json:"description,omitempty"`
	Date        string   `json:"release_date,omitempty"`
	Actors      []string `json:"actors,omitempty"`
	Rating      string   `json:"rating,omitempty"`
}

type NewActor struct {
	IdActor   int64    `json:"id"`
	FirstName string   `json:"first_name"`
	LatName   string   `json:"last_name"`
	Surname   string   `json:"surname"`
	Sex       string   `json:"sex"`
	Birthday  string   `json:"birthday"`
	Films     []string `json:"films"`
}

type ChangeDataActor struct {
	IdActor   int64    `json:"id"`
	FirstName string   `json:"first_name,omitempty"`
	LastName  string   `json:"last_name,omitempty"`
	Surname   string   `json:"surname,omitempty"`
	Sex       string   `json:"sex,omitempty"`
	Birthday  string   `json:"birthday,omitempty"`
	Films     []string `json:"films,omitempty"`
}

type CreateFilm struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Date        string   `json:"release_date"`
	Actors      []string `json:"actors"`
	Rating      string   `json:"rating"`
}

type Tokens struct {
	Access_token  string
	Refresh_token string
}
