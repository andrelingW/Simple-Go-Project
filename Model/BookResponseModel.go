package Model

type BookResponse struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Author    string `json:"author"`
	Available bool   `json:"available"`
}
