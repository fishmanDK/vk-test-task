package handlers

import (
	"errors"
	vk_test_task "github.com/fishmanDK/vk-test-task"
	"github.com/fishmanDK/vk-test-task/internal/service"
	"github.com/fishmanDK/vk-test-task/internal/service/mocks"
	"testing"
)

func TestHandler_createFilm(t *testing.T) {
	type fewFields error

	var e fewFields = errors.New("insufficient number of fields")

	cases := []struct {
		name    string
		newFilm vk_test_task.CreateFilm
		wantErr error
	}{
		{
			name: "Success",
			newFilm: vk_test_task.CreateFilm{
				Title:       "GTO",
				Description: "very good",
				Date:        "1997-01-01",
				Actors:      []string{"1"},
				Rating:      "10",
			},
			wantErr: nil,
		},
		{
			name:    "No id",
			wantErr: e,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			serviceFilmMock := mocks.NewFilms(t)

			serviceFilmMock.
				On("CreateFilm", tc.newFilm).
				Return(func(newFilm vk_test_task.CreateFilm) error {
					if newFilm.Rating == "" || newFilm.Date == "" || newFilm.Title == "" {
						return e
					}
					return nil
				})

			s := service.Service{
				Films: serviceFilmMock,
			}

			err := s.CreateFilm(tc.newFilm)
			if err != tc.wantErr {
				t.Errorf("CreateActor() error = %v, wantErr = %v", err, tc.wantErr)
				return
			}

		})
	}
}

func TestHandler_changeFilm(t *testing.T) {
	type fewFields error

	var e fewFields = errors.New("insufficient number of fields")

	cases := []struct {
		name           string
		idFilm         string
		changeDataFilm vk_test_task.ChangeDataFilm
		wantErr        error
	}{
		{
			name:   "Success",
			idFilm: "23",
			changeDataFilm: vk_test_task.ChangeDataFilm{
				Title:       "GTO",
				Description: "very good",
				Date:        "1997-01-01",
				Actors:      []string{"1"},
				Rating:      "10",
			},
			wantErr: nil,
		},
		{
			name: "No id",
			changeDataFilm: vk_test_task.ChangeDataFilm{
				Title:       "GTO",
				Description: "very good",
				Date:        "1997-01-01",
				Actors:      []string{"1"},
				Rating:      "10",
			},
			wantErr: e,
		},
		{
			name:    "No id",
			wantErr: e,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			serviceFilmMock := mocks.NewFilms(t)

			serviceFilmMock.
				On("ChangeFilm", tc.idFilm, tc.changeDataFilm).
				Return(func(idFilm string, changeDataFilm vk_test_task.ChangeDataFilm) error {
					if idFilm == "" {
						return e
					}
					return nil
				})

			s := service.Service{
				Films: serviceFilmMock,
			}

			err := s.ChangeFilm(tc.idFilm, tc.changeDataFilm)
			if err != tc.wantErr {
				t.Errorf("CreateActor() error = %v, wantErr = %v", err, tc.wantErr)
				return
			}

		})
	}
}

func TestHandler_deleteFilm(t *testing.T) {
	type fewFields error

	var e fewFields = errors.New("insufficient number of fields")

	cases := []struct {
		name    string
		idFilm  string
		wantErr error
	}{
		{
			name:    "Success",
			idFilm:  "23",
			wantErr: nil,
		},
		{
			name:    "No id",
			wantErr: e,
		},
		{
			name:    "No id",
			wantErr: e,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			serviceFilmMock := mocks.NewFilms(t)

			serviceFilmMock.
				On("DeleteFilmById", tc.idFilm).
				Return(func(idFilm string) error {
					if idFilm == "" {
						return e
					}
					return nil
				})

			s := service.Service{
				Films: serviceFilmMock,
			}

			err := s.DeleteFilmById(tc.idFilm)
			if err != tc.wantErr {
				t.Errorf("CreateActor() error = %v, wantErr = %v", err, tc.wantErr)
				return
			}

		})
	}
}
