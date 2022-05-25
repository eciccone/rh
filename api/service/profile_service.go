package service

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/eciccone/rh/api/repo/profile"
	"github.com/eciccone/rh/api/rherr"
)

type ProfileService interface {
	CreateProfile(args profile.Profile) error
}

type profileService struct {
	profileRepo profile.ProfileRepository
}

func NewProfileService(profileRepo profile.ProfileRepository) ProfileService {
	return &profileService{profileRepo}
}

func (s *profileService) CreateProfile(args profile.Profile) error {
	if args.Id == "" || args.Username == "" {
		return rherr.ErrBadRequest
	}

	// check if profile exists
	_, err := s.profileRepo.SelectProfileById(args.Id)
	if !errors.Is(err, sql.ErrNoRows) {
		if err != nil {
			return fmt.Errorf("CreateProfile failed to get profile by id: %w", err)
		}
		// profile exists
		return rherr.ErrBadRequest
	}

	// check if username is in use
	_, err = s.profileRepo.SelectProfileByUsername(args.Username)
	if !errors.Is(err, sql.ErrNoRows) {
		if err != nil {
			return fmt.Errorf("CreateProfile failed to get profile by username: %w", err)
		}
		// username exists
		return rherr.ErrBadRequest
	}

	if err = s.profileRepo.InsertProfile(args); err != nil {
		return fmt.Errorf("CreateProfile failed to create profile: %w", err)
	}

	return nil
}
