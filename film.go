package vk_test_task

type ChangeDataFilm struct {
	Title       string   `json:"title,omitempty"`
	Description string   `json:"description,omitempty"`
	Date        string   `json:"release_date,omitempty"`
	Actors      []string `json:"actors,omitempty"`
	Rating      string   `json:"rating,omitempty"`
}

type CreateFilm struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Date        string   `json:"release_date"`
	Actors      []string `json:"actors"`
	Rating      string   `json:"rating"`
}
