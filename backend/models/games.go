package models

import "time"

type Game struct {
	ID           uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	Title        string     `json:"title"        gorm:"type:varchar(200);not null;index:idx_title_platform,priority:1"`
	Platform     string     `json:"platform"     gorm:"type:varchar(80);not null;index:idx_title_platform,priority:2"`
	Genre        string     `json:"genre"        gorm:"type:varchar(80);index"`
	Status       string     `json:"status"       gorm:"type:varchar(32);index"`
	Progress     int        `json:"progress"     gorm:"type:int;check:progress_between_0_100,progress >= 0 AND progress <= 100"`
	HoursPlayed  float64    `json:"hoursPlayed"  gorm:"type:decimal(10,2);default:0"` 
	PersonalNote string     `json:"personalNote" gorm:"type:text"`
	Score        int        `json:"score"        gorm:"type:int;check:score_between_0_10,score >= 0 AND score <= 10"`
	StartedAt    *time.Time `json:"startedAt"    gorm:"index"`
	FinishedAt   *time.Time `json:"finishedAt"   gorm:"index"`
	CoverURL     string     `json:"coverURL"     gorm:"type:varchar(500)"`
	CreatedAt    time.Time  `json:"createdAt"    gorm:"not null"`
	UpdatedAt    time.Time  `json:"updatedAt"    gorm:"not null"`
}

type GameStats struct {
	TotalGames      int            `json:"total_games"`
	ByStatus        map[string]int `json:"by_status"`
	AverageHours    float64        `json:"average_hours_played"`
	MostPlayedGenre string         `json:"most_played_genre"`
	PendingGames    int            `json:"pending_games"`
}