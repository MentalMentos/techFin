package mock

import (
	"context"
	"github.com/MentalMentos/techFin/internal/clients/db"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"

	"github.com/stretchr/testify/mock"
)

type MockTx struct {
	mock.Mock
}

func (m *MockTx) Commit(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockTx) Rollback(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

type MockTxManager struct {
	mock.Mock
}

func (m *MockTxManager) ReadCommitted(ctx context.Context, f db.TxHandler) error {
	args := m.Called(ctx, f)
	return args.Error(0)
}

func (m *MockTxManager) RepeatableRead(ctx context.Context, f db.TxHandler) error {
	args := m.Called(ctx, f)
	return args.Error(0)
}

func (m *MockTxManager) Serializable(ctx context.Context, f db.TxHandler) error {
	args := m.Called(ctx, f)
	return args.Error(0)
}

func (m *MockTx) Begin(ctx context.Context) (pgx.Tx, error) {
	args := m.Called(ctx)
	return args.Get(0).(pgx.Tx), args.Error(1)
}

func (m *MockTx) BeginFunc(ctx context.Context, f func(pgx.Tx) error) error {
	args := m.Called(ctx, f)
	return args.Error(0)
}

func (m *MockTx) CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error) {
	args := m.Called(ctx, tableName, columnNames, rowSrc)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockTx) SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults {
	args := m.Called(ctx, b)
	return args.Get(0).(pgx.BatchResults)
}

func (m *MockTx) LargeObjects() pgx.LargeObjects {
	args := m.Called()
	return args.Get(0).(pgx.LargeObjects)
}

func (m *MockTx) Prepare(ctx context.Context, name string, sql string) (*pgconn.StatementDescription, error) {
	args := m.Called(ctx, name, sql)
	return args.Get(0).(*pgconn.StatementDescription), args.Error(1)
}

// Добавляем другие необходимые методы
func (m *MockTx) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	args = append([]interface{}{ctx, sql}, args...)
	res := m.Called(args...)
	return res.Get(0).(pgx.Rows), res.Error(1)
}

func (m *MockTx) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	args = append([]interface{}{ctx, sql}, args...)
	res := m.Called(args...)
	return res.Get(0).(pgx.Row)
}

func (m *MockTx) Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	args = append([]interface{}{ctx, sql}, args...)
	res := m.Called(args...)
	return res.Get(0).(pgconn.CommandTag), res.Error(1)
}

func (m *MockTx) QueryFunc(ctx context.Context, sql string, args []interface{}, scans []interface{}, f func(pgx.QueryFuncRow) error) (pgconn.CommandTag, error) {
	args = append([]interface{}{ctx, sql, args, scans, f})
	res := m.Called(args...)
	return res.Get(0).(pgconn.CommandTag), res.Error(1)
}

// Метод Conn
func (m *MockTx) Conn() *pgx.Conn {
	args := m.Called()
	return args.Get(0).(*pgx.Conn)
}
