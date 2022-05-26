package service

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/eciccone/rh/api/repo/profile"
)

var (
	ErrNoProfile         = errors.New("profile not found")
	ErrProfileData       = errors.New("must provide username for profile")
	ErrUsernameForbidden = errors.New("username not available")
	ErrProfileExists     = errors.New("profile already created")
)

type ProfileService interface {
	// Returns ErrProfileData if username is empty.
	// Returns ErrProfileExists if profile already exists.
	// Returns ErrUsernameForbidden if username is in use.
	CreateProfile(args profile.Profile) error

	// Returns ErrNoProfile if profile does not exist.
	FetchProfile(id string) (profile.Profile, error)
}

type profileService struct {
	profileRepo profile.ProfileRepository
}

func NewProfileService(profileRepo profile.ProfileRepository) ProfileService {
	return &profileService{profileRepo}
}

// Returns ErrNoProfile if profile does not exist.
func (s *profileService) FetchProfile(id string) (profile.Profile, error) {
	result, err := s.profileRepo.SelectProfileById(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return profile.Profile{}, ErrNoProfile
		}

		return profile.Profile{}, err
	}

	return result, nil
}

// Returns ErrProfileData if username is empty.
// Returns ErrProfileExists if profile already exists.
// Returns ErrUsernameForbidden if username is in use.
func (s *profileService) CreateProfile(args profile.Profile) error {
	if args.Username == "" {
		return ErrProfileData
	}

	// check if profile exists
	_, err := s.profileRepo.SelectProfileById(args.Id)
	if !errors.Is(err, sql.ErrNoRows) {
		if err != nil {
			return fmt.Errorf("CreateProfile failed to get profile by id: %w", err)
		}
		// profile exists
		return ErrProfileExists
	}

	// check if username is in use
	_, err = s.profileRepo.SelectProfileByUsername(args.Username)
	if !errors.Is(err, sql.ErrNoRows) {
		if err != nil {
			return fmt.Errorf("CreateProfile failed to get profile by username: %w", err)
		}
		// username exists
		return ErrUsernameForbidden
	}

	if err = s.profileRepo.InsertProfile(args); err != nil {
		return fmt.Errorf("CreateProfile failed to create profile: %w", err)
	}

	return nil
}
