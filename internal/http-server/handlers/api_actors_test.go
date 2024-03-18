package handlers

import (
	"errors"
	vk_test_task "github.com/fishmanDK/vk-test-task"
	"github.com/fishmanDK/vk-test-task/internal/service"
	"github.com/fishmanDK/vk-test-task/internal/service/mocks"
	"testing"
)

func TestHandler_deleteActor(t *testing.T) {
	type fewFields error

	var e fewFields = errors.New("ID field is required")

	cases := []struct {
		name                 string
		idActor              string
		wantErr              error
		expectedResponseBody string
	}{
		{
			name:    "Success",
			idActor: "4",
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

			serviceActorMock := mocks.NewActors(t)

			serviceActorMock.
				On("DeleteActor", tc.idActor).
				Return(func(idActor string) error {
					if idActor == "" {
						return e
					}
					return nil
				})

			s := service.Service{
				Actors: serviceActorMock,
			}

			err := s.DeleteActor(tc.idActor)
			if err != tc.wantErr {
				t.Errorf("CreateActor() error = %v, wantErr = %v", err, tc.wantErr)
				return
			}

		})
	}
}

func TestHandler_changeActor(t *testing.T) {
	type fewFields error

	var e fewFields = errors.New("ID field is required")

	cases := []struct {
		name                 string
		idActor              string
		changeActor          vk_test_task.ChangeDataActor
		wantErr              error
		expectedResponseBody string
	}{
		{
			name:    "Success",
			idActor: "4",
			changeActor: vk_test_task.ChangeDataActor{
				FirstName: "Denis",
				LastName:  "Svechnikov",
				Surname:   "Evgenievich",
				Sex:       "male",
				Birthday:  "2005-29-05",
				Films:     []string{"1", "-4", "3"},
			},
			wantErr: nil,
		},
		{
			name: "No id",
			changeActor: vk_test_task.ChangeDataActor{
				FirstName: "Denis",
				Sex:       "male",
				Birthday:  "2005-29-05",
				Films:     []string{"1", "-4", "3"},
			},
			wantErr: e,
		},
		{
			name:        "No id",
			changeActor: vk_test_task.ChangeDataActor{},
			wantErr:     e,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			serviceActorMock := mocks.NewActors(t)

			serviceActorMock.
				On("ChangeActor", tc.idActor, tc.changeActor).
				Return(func(idActor string, changeActor vk_test_task.ChangeDataActor) error {
					if idActor == "" {
						return e
					}
					return nil
				})

			s := service.Service{
				Actors: serviceActorMock,
			}

			err := s.ChangeActor(tc.idActor, tc.changeActor)
			if err != tc.wantErr {
				t.Errorf("CreateActor() error = %v, wantErr = %v", err, tc.wantErr)
				return
			}

		})
	}
}

func TestHandler_createActor(t *testing.T) {
	type fewFields error

	var e fewFields = errors.New("insufficient number of fields")

	cases := []struct {
		name                 string
		newActor             vk_test_task.NewActor
		wantErr              error
		expectedStatusCode   int
		inputBody            string
		expectedResponseBody string
	}{
		{
			name: "Success",
			newActor: vk_test_task.NewActor{
				FirstName: "Denis",
				LastName:  "Svechnikov",
				Surname:   "Evgenievich",
				Sex:       "male",
				Birthday:  "2005-29-05",
				Films:     []string{"1", "-4", "3"},
			},
			wantErr:            nil,
			expectedStatusCode: 200,
		},
		{
			name: "Few fields",
			newActor: vk_test_task.NewActor{
				FirstName: "Denis",
				Sex:       "male",
				Birthday:  "2005-29-05",
				Films:     []string{"1", "-4", "3"},
			},
			wantErr:            e,
			expectedStatusCode: 400,
		},
		{
			name: "Success",
			newActor: vk_test_task.NewActor{
				FirstName: "Denis",
				LastName:  "Svechnikov",
				Sex:       "male",
				Birthday:  "2005-29-05",
				Films:     []string{"1", "-4", "3"},
			},
			wantErr:            nil,
			expectedStatusCode: 200,
		},
		{
			name:               "Few fields",
			newActor:           vk_test_task.NewActor{},
			wantErr:            e,
			expectedStatusCode: 400,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			serviceActorMock := mocks.NewActors(t)

			serviceActorMock.
				On("CreateActor", tc.newActor).
				Return(func(newActor vk_test_task.NewActor) error {
					if newActor.FirstName == "" || newActor.LastName == "" || newActor.Sex == "" || newActor.Birthday == "" {
						return e
					}
					return nil
				})

			s := service.Service{
				Actors: serviceActorMock,
			}

			err := s.CreateActor(tc.newActor)
			if err != tc.wantErr {
				t.Errorf("CreateActor() error = %v, wantErr = %v", err, tc.wantErr)
				return
			}

		})
	}
}
