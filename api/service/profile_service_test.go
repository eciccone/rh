package service

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/eciccone/rh/api/repo/profile"
	"github.com/stretchr/testify/assert"
)

type ProfileRepoMocker struct {
	SelectProfileByIdMock       func(id string) (profile.Profile, error)
	SelectProfileByUsernameMock func(username string) (profile.Profile, error)
	InsertProfileMock           func(profile profile.Profile) error
}

func (r *ProfileRepoMocker) SelectProfileById(id string) (profile.Profile, error) {
	return r.SelectProfileByIdMock(id)
}

func (r *ProfileRepoMocker) SelectProfileByUsername(username string) (profile.Profile, error) {
	return r.SelectProfileByUsernameMock(username)
}

func (r *ProfileRepoMocker) InsertProfile(profile profile.Profile) error {
	return r.InsertProfileMock(profile)
}

func Test_CreateProfile(t *testing.T) {
	rr := &ProfileRepoMocker{
		SelectProfileByIdMock: func(id string) (profile.Profile, error) {
			return profile.Profile{}, sql.ErrNoRows
		},
		SelectProfileByUsernameMock: func(username string) (profile.Profile, error) {
			return profile.Profile{}, sql.ErrNoRows
		},
		InsertProfileMock: func(profile profile.Profile) error {
			return nil
		},
	}

	rs := NewProfileService(rr)

	err := rs.CreateProfile(profile.Profile{"test-id", "test user"})

	assert.NoError(t, err)
}

func Test_CreateProfileIdExists(t *testing.T) {
	rr := &ProfileRepoMocker{
		SelectProfileByIdMock: func(id string) (profile.Profile, error) {
			return profile.Profile{}, nil
		},
		SelectProfileByUsernameMock: func(username string) (profile.Profile, error) {
			return profile.Profile{}, sql.ErrNoRows
		},
		InsertProfileMock: func(profile profile.Profile) error {
			return nil
		},
	}

	rs := NewProfileService(rr)

	err := rs.CreateProfile(profile.Profile{"test-id", "test user"})

	assert.Error(t, err)
}

func Test_CreateProfileUsernameExists(t *testing.T) {
	rr := &ProfileRepoMocker{
		SelectProfileByIdMock: func(id string) (profile.Profile, error) {
			return profile.Profile{}, sql.ErrNoRows
		},
		SelectProfileByUsernameMock: func(username string) (profile.Profile, error) {
			return profile.Profile{}, nil
		},
		InsertProfileMock: func(profile profile.Profile) error {
			return nil
		},
	}

	rs := NewProfileService(rr)

	err := rs.CreateProfile(profile.Profile{"test-id", "test user"})

	assert.Error(t, err)
}

func Test_CreateProfileError(t *testing.T) {
	rr := &ProfileRepoMocker{
		SelectProfileByIdMock: func(id string) (profile.Profile, error) {
			return profile.Profile{}, sql.ErrNoRows
		},
		SelectProfileByUsernameMock: func(username string) (profile.Profile, error) {
			return profile.Profile{}, sql.ErrNoRows
		},
		InsertProfileMock: func(profile profile.Profile) error {
			return errors.New("failed")
		},
	}

	rs := NewProfileService(rr)

	err := rs.CreateProfile(profile.Profile{"test-id", "test user"})

	assert.Error(t, err)
}

func Test_FetchProfile(t *testing.T) {
	p := profile.Profile{"test-id", "test user"}
	rr := &ProfileRepoMocker{
		SelectProfileByIdMock: func(id string) (profile.Profile, error) {
			return p, nil
		},
	}

	rs := NewProfileService(rr)

	result, err := rs.FetchProfile("test-id")

	assert.NoError(t, err)
	assert.Equal(t, p, result)
}

func Test_FetchProfileError(t *testing.T) {
	p := profile.Profile{}
	rr := &ProfileRepoMocker{
		SelectProfileByIdMock: func(id string) (profile.Profile, error) {
			return profile.Profile{}, errors.New("failed")
		},
	}

	rs := NewProfileService(rr)

	result, err := rs.FetchProfile("test-id")

	assert.Error(t, err)
	assert.Equal(t, p, result)
}
