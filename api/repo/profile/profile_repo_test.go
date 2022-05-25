package profile

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func Test_SelectProfileById(t *testing.T) {
	data := []struct {
		Id          string
		P           Profile
		ExpectedSQL func(mock sqlmock.Sqlmock, profile Profile)
		Pass        bool
		Assert      func(mock sqlmock.Sqlmock, expected Profile, actual Profile, err error)
	}{
		{
			Id: "test-id",
			P:  Profile{Id: "test-id", Username: "Test User"},
			ExpectedSQL: func(mock sqlmock.Sqlmock, profile Profile) {
				mock.ExpectQuery("SELECT id, username FROM profile WHERE id = ?").
					WithArgs(profile.Id).
					WillReturnRows(sqlmock.NewRows([]string{"id", "username"}).AddRow(profile.Id, profile.Username))
			},
			Pass: true,
			Assert: func(mock sqlmock.Sqlmock, expected, actual Profile, err error) {
				assert.NoError(t, err)
				assert.Equal(t, expected, actual)
			},
		},
		{
			Id: "test-id",
			P:  Profile{},
			ExpectedSQL: func(mock sqlmock.Sqlmock, profile Profile) {
				mock.ExpectQuery("SELECT id, username FROM profile WHERE id = ?").
					WillReturnError(errors.New("failed"))
			},
			Pass: false,
			Assert: func(mock sqlmock.Sqlmock, expected, actual Profile, err error) {
				assert.Error(t, err)
				assert.Equal(t, expected, actual)
			},
		},
	}

	db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	for _, d := range data {
		d.ExpectedSQL(mock, d.P)
		pr := NewRepo(db)
		p, err := pr.SelectProfileById(d.Id)
		d.Assert(mock, d.P, p, err)
	}
}

func Test_SelectProfileByUsername(t *testing.T) {
	data := []struct {
		Username    string
		P           Profile
		ExpectedSQL func(mock sqlmock.Sqlmock, profile Profile)
		Pass        bool
		Assert      func(mock sqlmock.Sqlmock, expected Profile, actual Profile, err error)
	}{
		{
			Username: "Test User",
			P:        Profile{Id: "test-id", Username: "Test User"},
			ExpectedSQL: func(mock sqlmock.Sqlmock, profile Profile) {
				mock.ExpectQuery("SELECT id, username FROM profile WHERE username = ?").
					WithArgs(profile.Username).
					WillReturnRows(sqlmock.NewRows([]string{"id", "username"}).AddRow(profile.Id, profile.Username))
			},
			Pass: true,
			Assert: func(mock sqlmock.Sqlmock, expected, actual Profile, err error) {
				assert.NoError(t, err)
				assert.Equal(t, expected, actual)
			},
		},
		{
			Username: "Test User",
			P:        Profile{},
			ExpectedSQL: func(mock sqlmock.Sqlmock, profile Profile) {
				mock.ExpectQuery("SELECT id, username FROM profile WHERE username = ?").
					WillReturnError(errors.New("failed"))
			},
			Pass: false,
			Assert: func(mock sqlmock.Sqlmock, expected, actual Profile, err error) {
				assert.Error(t, err)
				assert.Equal(t, expected, actual)
			},
		},
	}

	db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	for _, d := range data {
		d.ExpectedSQL(mock, d.P)
		pr := NewRepo(db)
		p, err := pr.SelectProfileByUsername(d.Username)
		d.Assert(mock, d.P, p, err)
	}
}

func Test_InsertProfile(t *testing.T) {
	data := []struct {
		P           Profile
		ExpectedSQL func(mock sqlmock.Sqlmock, profile Profile)
		Pass        bool
		Assert      func(mock sqlmock.Sqlmock, err error)
	}{
		{
			P: Profile{Id: "test-id", Username: "Test User"},
			ExpectedSQL: func(mock sqlmock.Sqlmock, profile Profile) {
				mock.ExpectExec("INSERT INTO profile(id, username) VALUES (?, ?)").
					WithArgs(profile.Id, profile.Username).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			Pass: true,
			Assert: func(mock sqlmock.Sqlmock, err error) {
				assert.NoError(t, err)
			},
		},
		{
			P: Profile{Id: "test-id", Username: "Test User"},
			ExpectedSQL: func(mock sqlmock.Sqlmock, profile Profile) {
				mock.ExpectExec("INSERT INTO profile(id, username) VALUES (?, ?)").
					WithArgs(profile.Id, profile.Username).
					WillReturnError(errors.New("failed"))
			},
			Pass: true,
			Assert: func(mock sqlmock.Sqlmock, err error) {
				assert.Error(t, err)
			},
		},
	}

	db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	for _, d := range data {
		d.ExpectedSQL(mock, d.P)
		pr := NewRepo(db)
		err := pr.InsertProfile(d.P)
		d.Assert(mock, err)
	}
}
