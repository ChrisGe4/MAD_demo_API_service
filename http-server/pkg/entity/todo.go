package entity

type Todo struct {
	ID          string   `json:"id"`
	Description string   `json:"description"`
	Category    Category `json:"category"`
}

type Category struct {
	Name string `json:"name"`
}
