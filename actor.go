package vk_test_task

type NewActor struct {
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
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
