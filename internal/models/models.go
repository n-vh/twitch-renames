package models

import "time"

type Rename struct {
	ID             uint   `gorm:"primaryKey"`
	UserId         string `gorm:"index"`
	Login          string `gorm:"index"`
	DisplayName    string
	OldLogin       string `gorm:"index"`
	OldDisplayName string
	Date           time.Time
}

type User struct {
	ID          uint   `gorm:"primaryKey"`
	UserId      string `gorm:"uniqueIndex"`
	Login       string
	DisplayName string
}

type Worker struct {
	WorkerId int `gorm:"index"`
	Cycles   int
	Offset   int
}
