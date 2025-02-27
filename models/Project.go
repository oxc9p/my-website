package models

type Project struct {
	Image       string `json:"image"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Label       string `json:"label"`
	Link        string `json:"link"`
}
