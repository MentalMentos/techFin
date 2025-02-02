package mock

import (
	"context"
	"github.com/MentalMentos/techFin/internal/clients/db"
	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/mock"
)

// MockPgClient реализует интерфейс db.Client
type MockPgClient struct {
	mock.Mock
}

func (m *MockPgClient) DB() db.DB {
	args := m.Called()
	return args.Get(0).(db.DB)
}

func (m *MockPgClient) Close() error {
	args := m.Called()
	return args.Error(0)
}

// Мок для Ping метода
func (m *MockPgClient) Ping(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// Мок для создания нового транзакции
func (m *MockPgClient) BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error) {
	args := m.Called(ctx, txOptions)
	return args.Get(0).(pgx.Tx), args.Error(1)
}
