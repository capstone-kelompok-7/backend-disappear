package dto

import "time"

type CreateChallengeRequest struct {
	Title       string    `form:"title" json:"title" validate:"required"`
	Photo       string    `form:"photo" json:"photo"`
	StartDate   time.Time `form:"start_date" json:"start_date" validate:"required"`
	EndDate     time.Time `form:"end_date" json:"end_date" validate:"required"`
	Description string    `form:"description" json:"description" validate:"required"`
	Exp         uint64    `form:"exp" json:"exp" validate:"required"`
	Status      string    `form:"status" json:"status"`
}

type UpdateChallengeRequest struct {
	Title       string    `form:"title" json:"title"`
	Photo       string    `form:"photo" json:"photo"`
	StartDate   time.Time `form:"start_date" json:"start_date"`
	EndDate     time.Time `form:"end_date" json:"end_date"`
	Description string    `form:"description" json:"description"`
	Exp         uint64    `form:"exp" json:"exp"`
	Status      string    `form:"status" json:"status"`
}

type CreateChallengeFormRequest struct {
	ChallengeID uint64 `form:"challenge_id" json:"challenge_id" validate:"required"`
	Username    string `form:"username" json:"username" validate:"required"`
	Photo       string `form:"photo" json:"photo"`
}

type UpdateChallengeFormStatusRequest struct {
	Status string `form:"status" json:"status" validate:"required"`
}
