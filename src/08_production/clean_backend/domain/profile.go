package domain

import "context"

// Profile is a read-only DTO it strips password from User.
type Profile struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// ProfileUsecase defines the business contract for fetching a user profile.
type ProfileUsecase interface {
	GetProfileByID(c context.Context, userID string) (*Profile, error)
}
