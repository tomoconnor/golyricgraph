package main

import (
	// "github.com/getsentry/sentry-go"
	// sentryecho "github.com/getsentry/sentry-go/echo"

	"gorm.io/gorm"
)

type TokenData struct {
	gorm.Model
	VerificationToken string `json:"verification_token"`
	Email             string `json:"email"`
	TokenExpired      bool   `json:"token_expired"`
}

type LyricGraph struct {
	gorm.Model
	artist   string `json:"artist"`
	title    string `json:"title"`
	filename string `json:"filename"`
	filepath string `json:"filepath"`
	rank     int    `json:"rank"`
}

type LyricGraphComparison struct {
	gorm.Model
	artist1  string `json:"artist1"`
	artist2  string `json:"artist2"`
	title1   string `json:"title1"`
	title2   string `json:"title2"`
	filename string `json:"filename"`
	filepath string `json:"filepath"`
	rank     int    `json:"rank"`
}
