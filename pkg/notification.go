package pkg

import "time"

type Notification struct {
	RecipientID int64          `json:"recipient_id"`
	Type        string         `json:"type"`
	Priority    int            `json:"priority"`
	Payload     VacancyPayload `json:"payload"`
	CreatedAt   time.Time      `json:"created_at"`
}

type VacancyPayload struct {
	ID       string   `json:"id"`
	Title    string   `json:"title"`
	Company  string   `json:"company"`
	Salary   string   `json:"salary"`
	Location string   `json:"location"`
	Link     string   `json:"link"`
	Keywords []string `json:"keywords"`
}
