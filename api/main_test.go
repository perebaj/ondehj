package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/perebaj/ondehj/event"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockSQLRepository struct {
	mock.Mock
}

func NewMockSQLRepository() *MockSQLRepository {
	return &MockSQLRepository{}
}

func (m *MockSQLRepository) Create(a context.Context, b event.Event, log zerolog.Logger) (*event.Event, error) {
	args := m.Called(a, b)
	return args.Get(0).(*event.Event), args.Error(1)
}

func (m *MockSQLRepository) Update(ctx context.Context, id int64, newEvent event.Event, log zerolog.Logger) (*event.Event, error) {
	args := m.Called(ctx, id, newEvent)
	return args.Get(0).(*event.Event), args.Error(1)
}

func (m *MockSQLRepository) GetByID(ctx context.Context, id int64, log zerolog.Logger) (*event.Event, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*event.Event), args.Error(1)

}

func (m *MockSQLRepository) Delete(ctx context.Context, id int64, log zerolog.Logger) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockSQLRepository) Migrate() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockSQLRepository) All(ctx context.Context, log zerolog.Logger) ([]event.Event, error) {
	args := m.Called(ctx)
	return args.Get(0).([]event.Event), args.Error(1)
}

type MockEvent interface {
	Create(ctx context.Context, event event.Event, log zerolog.Logger) (*event.Event, error)
	Update(ctx context.Context, id int64, newEvent event.Event, log zerolog.Logger) (*event.Event, error)
	GetByID(ctx context.Context, id int64, log zerolog.Logger) (*event.Event, error)
	Delete(ctx context.Context, id int64, log zerolog.Logger) error
	Migrate() error
	All(ctx context.Context, log zerolog.Logger) ([]event.Event, error)
}

func Test_postCreateEventHandler(t *testing.T) {
	testCases := []struct {
		name               string
		event              *event.Event
		expectedStatusCode int
		method             string
	}{
		{
			name: "Create event",
			event: &event.Event{
				Title:         "Example Event JOjo is here",
				Description:   "This is an example event.",
				Location:      "New York City",
				StartTime:     time.Now(),
				EndTime:       time.Now(),
				InstagramPage: "example_event",
			},
			expectedStatusCode: 200,
			method:             "POST",
		},
		{
			name: "Invalid method",
			event: &event.Event{
				Title:         "JOJO is awesome",
				Description:   "Jojo mage",
				Location:      "Jojo Town",
				StartTime:     time.Now(),
				EndTime:       time.Now(),
				InstagramPage: "jojo",
			},
			expectedStatusCode: 405,
			method:             "GET",
		},
		{
			name:               "Empty event",
			event:              nil,
			expectedStatusCode: 400,
			method:             "POST",
		},
		{
			name:               "Empty event2",
			event:              &event.Event{},
			expectedStatusCode: 400,
			method:             "POST",
		},
		{
			name: "Empty title",
			event: &event.Event{
				Title:         "",
				Description:   "Jojo mage",
				Location:      "Jojo Town",
				StartTime:     time.Now(),
				EndTime:       time.Now(),
				InstagramPage: "jojo",
			},
			expectedStatusCode: 400,
			method:             "POST",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// convert the event to a JSON string
			eventJson, err := json.Marshal(tc.event)
			if err != nil {
				fmt.Println("Error marshaling event:", err)
				return
			}
			eventString := string(eventJson)

			mockRepo := NewMockSQLRepository()
			mockRepo.On("Create", mock.Anything, mock.Anything).Return(tc.event, nil)
			resultHandlerFunc := postCreateEventHandler(mockRepo)
			req := httptest.NewRequest(tc.method, "/events", strings.NewReader(eventString))
			w := httptest.NewRecorder()
			resultHandlerFunc.ServeHTTP(w, req) // doing the fake request
			res := w.Result()                   // capturing the response
			assert.Equal(t, tc.expectedStatusCode, res.StatusCode)
		})
	}
}

func Test_deleteEventHandler(t *testing.T) {
	type deleteReturn struct {
		err error // error to be returned by the mock repo when Delete is called
	}
	type getByIdReturn struct {
		event *event.Event // event to be returned by the mock repo when GetByID is called
		err   error        // error to be returned by the mock repo when GetByID is called
	}
	testCases := []struct {
		name               string        //subtest name
		deleteReturn       deleteReturn  // return values for the mock repo's Delete method
		getByIdReturn      getByIdReturn // return values for the mock repo's GetByID method
		expectedStatusCode int
		method             string
		requestIdParam     string
		deleteError        error
	}{
		{
			name: "Delete event",
			getByIdReturn: getByIdReturn{
				err: nil,
				event: &event.Event{
					ID:            2,
					Title:         "Example Event JOjo is here",
					Description:   "This is an example event.",
					Location:      "New York City",
					StartTime:     time.Now(),
					EndTime:       time.Now(),
					InstagramPage: "example_event",
				},
			},
			deleteReturn: deleteReturn{
				err: nil,
			},
			expectedStatusCode: 200,
			method:             "DELETE",
			requestIdParam:     "2",
		},
		{
			name:               "Invalid method",
			expectedStatusCode: 405,
			method:             "GET",
			requestIdParam:     "2",
			deleteReturn: deleteReturn{
				err: nil,
			},
			getByIdReturn: getByIdReturn{
				err:   nil,
				event: &event.Event{},
			},
		},
		{
			name:               "Invalid parameter value",
			method:             "DELETE",
			requestIdParam:     "invalid",
			expectedStatusCode: 400,
			deleteReturn: deleteReturn{
				err: nil,
			},
			getByIdReturn: getByIdReturn{
				err:   nil,
				event: &event.Event{},
			},
		},
		{
			name:               "Event not found",
			method:             "DELETE",
			requestIdParam:     "2",
			expectedStatusCode: 404,
			getByIdReturn: getByIdReturn{
				err:   sql.ErrNoRows,
				event: nil,
			},
			deleteReturn: deleteReturn{
				err: nil,
			},
		},
		{
			name: "Delete failed",

			deleteReturn: deleteReturn{
				err: event.ErrDeleteFailed,
			},
			getByIdReturn: getByIdReturn{
				err: nil,
				event: &event.Event{
					ID:            2,
					Title:         "Example Event JOjo is here",
					Description:   "This is an example event.",
					Location:      "New York City",
					StartTime:     time.Now(),
					EndTime:       time.Now(),
					InstagramPage: "example_event",
				},
			},

			method:             "DELETE",
			requestIdParam:     "2",
			expectedStatusCode: 500,
			deleteError:        event.ErrDeleteFailed,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := NewMockSQLRepository()
			mockRepo.On("GetByID", mock.Anything, mock.Anything).Return(tc.getByIdReturn.event, tc.getByIdReturn.err)
			mockRepo.On("Delete", mock.Anything, mock.Anything).Return(tc.deleteReturn.err)

			resultHandlerFunc := deleteEventHandler(mockRepo)
			fmt.Printf("/events/%s", tc.requestIdParam)

			req := httptest.NewRequest(tc.method, "/event/"+tc.requestIdParam, nil)
			req = mux.SetURLVars(req, map[string]string{"id": tc.requestIdParam})
			w := httptest.NewRecorder()
			resultHandlerFunc.ServeHTTP(w, req) // doing the fake request
			res := w.Result()                   // capturing the response
			assert.Equal(t, tc.expectedStatusCode, res.StatusCode)
		})
	}

}
