package mock

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"
)

// MockRedisClient - структура для мока Redis клиента
type MockRedisClient struct {
	mock.Mock
}

// Get - имитирует метод Get
func (m *MockRedisClient) Get(ctx context.Context, key string) (string, error) {
	args := m.Called(ctx, key)
	return args.String(0), args.Error(1)
}

// Set - имитирует метод Set
func (m *MockRedisClient) Set(ctx context.Context, key, value string, expiration time.Duration) error {
	args := m.Called(ctx, key, value, expiration)
	return args.Error(0)
}

// SetObject - имитирует метод SetObject
func (m *MockRedisClient) SetObject(ctx context.Context, key string, obj interface{}, expiration time.Duration) error {
	args := m.Called(ctx, key, obj, expiration)
	return args.Error(0)
}

// GetObject - имитирует метод GetObject
func (m *MockRedisClient) GetObject(ctx context.Context, key string, obj interface{}) error {
	args := m.Called(ctx, key, obj)
	return args.Error(0)
}

func (m *MockRedisClient) Delete(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}
