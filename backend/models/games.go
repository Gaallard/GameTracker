package models

import "time"

type Game struct {
	ID           uint       `json:"id" gorm:"primaryKey"`
	Title        string     `json:"title"`
	Platform     string     `json:"platform"`
	Genre        string     `json:"genre"`
	Status       string     `json:"status"`
	Progress     int        `json:"progress"`
	HoursPlayed  float64    `json:"hoursPlayed"`
	PersonalNote string     `json:"personalNote"`
	Score        int        `json:"score"`
	StartedAt    *time.Time `json:"startedAt"`
	FinishedAt   *time.Time `json:"finishedAt"`
	CoverURL     string     `json:"coverURL"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type GameStats struct {
	TotalGames      int            `json:"total_games"`
	ByStatus        map[string]int `json:"by_status"`
	AverageHours    float64        `json:"average_hours_played"`
	MostPlayedGenre string         `json:"most_played_genre"`
	PendingGames    int            `json:"pending_games"`
}
