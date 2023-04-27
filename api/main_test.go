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
	eventExample := &event.Event{
		Title:         "Example Event JOjo is here",
		Description:   "This is an example event.",
		Location:      "New York City",
		StartTime:     time.Now(),
		EndTime:       time.Now(),
		InstagramPage: "example_event",
	}
	// convert the event to a JSON string
	eventJson, err := json.Marshal(eventExample)
	if err != nil {
		fmt.Println("Error marshaling event:", err)
		return
	}
	eventString := string(eventJson)

	mockRepo := NewMockSQLRepository()
	mockRepo.On("Create", mock.Anything, mock.Anything).Return(eventExample, nil)
	resultHandlerFunc := postCreateEventHandler(mockRepo)
	req := httptest.NewRequest("POST", "/events", strings.NewReader(eventString))
	w := httptest.NewRecorder()
	resultHandlerFunc.ServeHTTP(w, req) // doing the fake request
	res := w.Result()                   // capturing the response
	assert.Equal(t, 200, res.StatusCode)
	// t.Log(res.StatusCode)
	// defer res.Body.Close()
	// data, err := ioutil.ReadAll(res.Body)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// t.Log(string(data))
}

// func Test_postCreateEventHandler(t *testing.T) {
// 	type args struct {
// 		eventRepo event.Repository
// 	}
// 	type want struct {
// 		status int
// 	}
// 	tests := []struct {
// 		name  string
// 		args  args
// 		event event.Event
// 		want  want
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := postCreateEventHandler(tt.args.eventRepo); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("postCreateEventHandler() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
