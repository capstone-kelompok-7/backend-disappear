package dto

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"time"
)

type ChallengeFormatter struct {
	ID          uint64 `json:"id"`
	Title       string `json:"title"`
	Photo       string `json:"photo"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
	Description string `json:"description"`
	Winner      string `json:"winner"`
	Status      string `json:"status"`
	Exp         uint64 `json:"exp"`
}

func FormatChallenge(challenge entities.ChallengeModels) ChallengeFormatter {
	challengeFormatter := ChallengeFormatter{}
	challengeFormatter.ID = challenge.ID
	challengeFormatter.Title = challenge.Title
	challengeFormatter.Photo = challenge.Photo
	challengeFormatter.StartDate = challenge.StartDate.Format("2006-01-02")
	challengeFormatter.EndDate = challenge.EndDate.Format("2006-01-02")
	challengeFormatter.Description = challenge.Description
	challengeFormatter.Winner = challenge.Winner
	challengeFormatter.Exp = challenge.Exp

	now := time.Now()
	endd := challenge.EndDate

	if challenge.Winner != "" {
		challengeFormatter.Status = "berakhir"
	} else if now.After(endd) {
		challengeFormatter.Status = "berakhir"
	} else if now.Before(endd) {
		challengeFormatter.Status = "berlangsung"
	}

	return challengeFormatter
}

func FormatterChallenge(challenge []entities.ChallengeModels) []ChallengeFormatter {
	var challengeFormatter []ChallengeFormatter

	for _, challenge := range challenge {
		FormatChallenge := FormatChallenge(challenge)
		challengeFormatter = append(challengeFormatter, FormatChallenge)
	}

	return challengeFormatter
}
