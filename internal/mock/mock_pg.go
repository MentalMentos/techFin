package mock

import (
	"context"
	"github.com/MentalMentos/techFin/internal/clients/db"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/mock"
)

// MockPg реализует интерфейс db.DB
type MockPg struct {
	mock.Mock
}

func (m *MockPg) ScanOneContext(ctx context.Context, dest interface{}, q db.Query, args ...interface{}) error {
	args2 := m.Called(ctx, dest, q, args)
	return args2.Error(0)
}

func (m *MockPg) ScanAllContext(ctx context.Context, dest interface{}, q db.Query, args ...interface{}) error {
	args2 := m.Called(ctx, dest, q, args)
	return args2.Error(0)
}

func (m *MockPg) ExecContext(ctx context.Context, q db.Query, args ...interface{}) (pgconn.CommandTag, error) {
	args2 := m.Called(ctx, q, args)
	return args2.Get(0).(pgconn.CommandTag), args2.Error(1)
}

func (m *MockPg) QueryContext(ctx context.Context, q db.Query, args ...interface{}) (pgx.Rows, error) {
	args2 := m.Called(ctx, q, args)
	return args2.Get(0).(pgx.Rows), args2.Error(1)
}

func (m *MockPg) QueryRowContext(ctx context.Context, q db.Query, args ...interface{}) pgx.Row {
	args2 := m.Called(ctx, q, args)
	return args2.Get(0).(pgx.Row)
}

func (m *MockPg) BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error) {
	args := m.Called(ctx, txOptions)
	return args.Get(0).(pgx.Tx), args.Error(1)
}

func (m *MockPg) Ping(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockPg) Close() {
	m.Called()
}
