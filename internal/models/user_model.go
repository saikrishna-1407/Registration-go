package model

import "time"

type User struct {
	ID            uint      `gorm:"primaryKey"`
	FullName      string    `gorm:"size:50;not null"`
	Email         string    `gorm:"unique;not null"`
	Password      string    `gorm:"not null"`
	Username      string    `gorm:"unique;size:20;not null"`
	DateOfBirth   time.Time `gorm:"not null"`
	PhoneNumber   string    `gorm:"not null"`
	Gender        string
	Country       string `gorm:"not null"`
	State         string
	PinCode       string
	ReferralCode  string
	TermsAccepted bool `gorm:"not null"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
