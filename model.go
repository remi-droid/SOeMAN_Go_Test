package main

import "time"

type Document struct {
	ID         int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Name       string    `gorm:"unique;	not null" json:"name"`
	Url        string    `gorm:"unique; not null" json:"url"`
	UploadDate time.Time `gorm:"autoUpdateTime" json:"uploaded_at"`
}
