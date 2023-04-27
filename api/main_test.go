package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/perebaj/ondehj/event"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockSQLRepository struct {
	mock.Mock
}

func NewMockSQLRepository() *MockSQLRepository {
	return &MockSQLRepository{}
}

func (m *MockSQLRepository) Create(a context.Context, b event.Event) (*event.Event, error) {
	args := m.Called(a, b)
	return args.Get(0).(*event.Event), args.Error(1)
}

func (m *MockSQLRepository) Update(ctx context.Context, id int64, newEvent event.Event) (*event.Event, error) {
	args := m.Called(ctx, id, newEvent)
	return args.Get(0).(*event.Event), args.Error(1)
}

func (m *MockSQLRepository) GetByID(ctx context.Context, id int64) (*event.Event, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*event.Event), args.Error(1)

}

func (m *MockSQLRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockSQLRepository) Migrate() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockSQLRepository) All(ctx context.Context) ([]event.Event, error) {
	args := m.Called(ctx)
	return args.Get(0).([]event.Event), args.Error(1)
}

type MockEvent interface {
	Create(ctx context.Context, event event.Event) (*event.Event, error)
	Update(ctx context.Context, id int64, newEvent event.Event) (*event.Event, error)
	GetByID(ctx context.Context, id int64) (*event.Event, error)
	Delete(ctx context.Context, id int64) error
	Migrate() error
	All(ctx context.Context) ([]event.Event, error)
}

func Test_postCreateEventHandler(t *testing.T) {
	testCases := []struct {
		name               string
		event              *event.Event
		extectedStatusCode int
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
			extectedStatusCode: 200,
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
			extectedStatusCode: 405,
			method:             "GET",
		},
		{
			name:               "Empty event",
			event:              nil,
			extectedStatusCode: 400,
			method:             "POST",
		},
		{
			name:               "Empty event2",
			event:              &event.Event{},
			extectedStatusCode: 400,
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
			assert.Equal(t, tc.extectedStatusCode, res.StatusCode)
		})
	}
}
