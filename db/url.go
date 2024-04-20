package db

import (
	"time"

	"gorm.io/gorm"
)

type CrawleUrl struct {
	ID              string         `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Url             string         `json:"url" gorm:"unique;not null"`
	Success         bool           `json:"success" gorm:"not null"`
	CrawlDuration   time.Duration  `json:"crawleDuration"`
	ResponseCode    int            `json:"responseCode"`
	PageTitle       string         `json:"pageTitle"`
	PageDescription string         `json:"pageDescription"`
	Heading         string         `json:"heading"`
	LastTested      *time.Time     `json:"lastTested"`
	Indexed         bool           `json:"indexed" gorm:"default:false"`
	CreatedAt       *time.Time     `gorm:"autoCreateTime"`
	UpdatedAt       time.Time      `gorm:"autoUpdateTime"`
	DeletedAt       gorm.DeletedAt `gorm:"index"`
}
