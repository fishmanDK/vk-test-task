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

type Tokens struct {
	Access_token  string
	Refresh_token string
}
