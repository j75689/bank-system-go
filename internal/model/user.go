package model

import "time"

type User struct {
	ID             uint64     `json:"-" gorm:"primarykey"`
	UUID           string     `json:"uuid" gorm:"nchar(36)"`
	Name           string     `json:"name" gorm:"nvarchar(255)"`
	Account        string     `json:"account" gorm:"nvarchar(255);uniqueIndex"`
	Password       []byte     `json:"password" gorm:"text"`
	LatestAccessAt *time.Time `json:"latest_access_at"`
	CreatedAt      time.Time  `json:"created_at" gorm:"not null;default:now()"`
	UpdatedAt      time.Time  `json:"updated_at" gorm:"not null;default:now()"`
	DeletedAt      *time.Time `json:"deleted_at" gorm:"index"`
}

type AccessLog struct {
	ID        uint64    `json:"id" gorm:"primarykey"`
	UserID    uint64    `json:"user_id" gorm:"index"`
	IP        string    `json:"ip" gorm:"varchar(255)"`
	CreatedAt time.Time `json:"created_at" gorm:"not null;default:now()"`
}
