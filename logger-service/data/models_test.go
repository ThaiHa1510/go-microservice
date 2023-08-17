package data_test

import (
	"context"
	"log-service/data"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockMongoClient struct {
	mock.Mock
}

func (m *MockMongoClient) InsertOne(ctx context.Context, document interface{}) error {
	args := m.Called(ctx, document)
	return args.Error(0)
}

func (m *MockMongoClient) Find(ctx context.Context, filter interface{}) ([]*data.LogEntry, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]*data.LogEntry), args.Error(1)
}

func TestInsert(t *testing.T) {
	mockDB := new(MockMongoClient)
	models := data.Models{
		LogEntry: data.LogEntry{},
	}

	entry := data.LogEntry{
		Name:      "Test Log Entry",
		Data:      "Test data",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Đặt kỳ vọng rằng phương thức InsertOne sẽ không trả về lỗi
	mockDB.On("InsertOne", mock.Anything, mock.Anything).Return(nil)

	err := models.LogEntry.Insert(entry)
	assert.NoError(t, err, "Insert should not return an error")
	mockDB.AssertCalled(t, "InsertOne", mock.Anything, mock.Anything)
}

func TestAll(t *testing.T) {
	mockDB := new(MockMongoClient)
	models := data.Models{
		LogEntry: data.LogEntry{},
	}

	// Tạo dữ liệu mock
	mockLogs := []*data.LogEntry{
		{ID: "1", Name: "Log 1", Data: "Data 1", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: "2", Name: "Log 2", Data: "Data 2", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: "3", Name: "Log 3", Data: "Data 3", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}

	// Đặt kỳ vọng rằng phương thức Find sẽ trả về dữ liệu mockLogs
	mockDB.On("Find", mock.Anything, mock.Anything).Return(mockLogs, nil)

	logs, err := models.LogEntry.All()
	assert.NoError(t, err, "All should not return an error")
	assert.Equal(t, mockLogs, logs, "All should return correct logs")
	mockDB.AssertCalled(t, "Find", mock.Anything, mock.Anything)
}
