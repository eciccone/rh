package profile

import (
	"database/sql"
	"errors"
	"fmt"
)

type ProfileRepository interface {
	SelectProfileById(id string) (Profile, error)
	SelectProfileByUsername(username string) (Profile, error)
	InsertProfile(profile Profile) error
}

type profileRepo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) ProfileRepository {
	return &profileRepo{db}
}

func (r *profileRepo) SelectProfileById(id string) (Profile, error) {
	var result Profile

	row := r.db.QueryRow("SELECT id, username FROM profile WHERE id = ?", id)
	if err := row.Scan(&result.Id, &result.Username); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Profile{}, err
		}

		return Profile{}, fmt.Errorf("SelectProfileById failed to select profile: %w", err)
	}

	return result, nil
}

func (r *profileRepo) SelectProfileByUsername(username string) (Profile, error) {
	var result Profile

	row := r.db.QueryRow("SELECT id, username FROM profile WHERE username = ?", username)
	if err := row.Scan(&result.Id, &result.Username); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Profile{}, err
		}

		return Profile{}, fmt.Errorf("SelectProfileByUsername failed to select profile: %w", err)
	}

	return result, nil
}

func (r *profileRepo) InsertProfile(profile Profile) error {
	_, err := r.db.Exec("INSERT INTO profile(id, username) VALUES (?, ?)", profile.Id, profile.Username)
	if err != nil {
		return fmt.Errorf("InsertProfile failed to insert profile: %w", err)
	}

	return nil
}
