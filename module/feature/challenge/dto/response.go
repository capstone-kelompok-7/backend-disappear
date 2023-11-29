package dto

import (
	"time"

	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
)

type ChallengeFormatter struct {
	ID          uint64    `json:"id"`
	Title       string    `json:"title"`
	Photo       string    `json:"photo"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	Exp         uint64    `json:"exp"`
}

func FormatChallenge(challenge *entities.ChallengeModels) *ChallengeFormatter {
	challengeFormatter := &ChallengeFormatter{}
	challengeFormatter.ID = challenge.ID
	challengeFormatter.Title = challenge.Title
	challengeFormatter.Photo = challenge.Photo
	challengeFormatter.StartDate = challenge.StartDate
	challengeFormatter.EndDate = challenge.EndDate
	challengeFormatter.Description = challenge.Description
	challengeFormatter.Exp = challenge.Exp
	challengeFormatter.Status = challenge.Status

	return challengeFormatter
}

func FormatterChallenge(challenge []*entities.ChallengeModels) []*ChallengeFormatter {
	var challengeFormatter []*ChallengeFormatter

	for _, challenge := range challenge {
		FormatChallenge := FormatChallenge(challenge)
		challengeFormatter = append(challengeFormatter, FormatChallenge)
	}

	return challengeFormatter
}

type ChallengeFormFormatter struct {
	ID          uint64    `json:"id"`
	UserID      uint64    `json:"user_id"`
	ChallengeID uint64    `json:"challenge_id"`
	Username    string    `json:"username"`
	Photo       string    `json:"photo"`
	Status      string    `json:"status"`
	Exp         uint64    `json:"exp"`
	CreatedAt   time.Time `json:"tanggal_berpartisipasi"`
}

func FormatChallengeForm(form *entities.ChallengeFormModels) *ChallengeFormFormatter {
	challengeFormFormatter := &ChallengeFormFormatter{}
	challengeFormFormatter.ID = form.ID
	challengeFormFormatter.UserID = form.UserID
	challengeFormFormatter.ChallengeID = form.ChallengeID
	challengeFormFormatter.Username = form.Username
	challengeFormFormatter.Photo = form.Photo
	challengeFormFormatter.Status = form.Status
	challengeFormFormatter.Exp = form.Exp
	challengeFormFormatter.CreatedAt = form.CreatedAt

	return challengeFormFormatter
}

func FormatterChallengeForm(forms []*entities.ChallengeFormModels) []*ChallengeFormFormatter {
	var formFormatter []*ChallengeFormFormatter
	for _, form := range forms {
		formattedForm := FormatChallengeForm(form)
		formFormatter = append(formFormatter, formattedForm)
	}
	return formFormatter
}
