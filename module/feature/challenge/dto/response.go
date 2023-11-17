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

func FormatChallenge(challenge entities.ChallengeModels) ChallengeFormatter {
	challengeFormatter := ChallengeFormatter{}
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

func FormatterChallenge(challenge []entities.ChallengeModels) []ChallengeFormatter {
	var challengeFormatter []ChallengeFormatter

	for _, challenge := range challenge {
		FormatChallenge := FormatChallenge(challenge)
		challengeFormatter = append(challengeFormatter, FormatChallenge)
	}

	return challengeFormatter
}
