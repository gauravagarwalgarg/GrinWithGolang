package usecase

import (
	"context"
	"time"

	"github.com/GauravAgarwalGarg/grin-with-golang/src/08_production/clean_backend/domain"
)

type profileUsecase struct {
	userRepository domain.UserRepository
	contextTimeout time.Duration
}

// NewProfileUsecase creates a new ProfileUsecase with a given timeout.
func NewProfileUsecase(userRepository domain.UserRepository, timeout time.Duration) domain.ProfileUsecase {
	return &profileUsecase{
		userRepository: userRepository,
		contextTimeout: timeout,
	}
}

func (pu *profileUsecase) GetProfileByID(c context.Context, userID string) (*domain.Profile, error) {
	ctx, cancel := context.WithTimeout(c, pu.contextTimeout)
	defer cancel()

	user, err := pu.userRepository.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Map User → Profile (strip sensitive fields)
	return &domain.Profile{Name: user.Name, Email: user.Email}, nil
}
