package repo_mocks

import (
	"context"
	"github.com/MentalMentos/techFin/internal/models"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v4"
	"reflect"
)

type MockTransactionRepository struct {
	ctrl     *gomock.Controller
	recorder *MockTransactionRepositoryMockRecorder
}

type MockTransactionRepositoryMockRecorder struct {
	mock *MockTransactionRepository
}

func NewMockTransactionRepository(ctrl *gomock.Controller) *MockTransactionRepository {
	mock := &MockTransactionRepository{ctrl: ctrl}
	mock.recorder = &MockTransactionRepositoryMockRecorder{mock}
	return mock
}

func (m *MockTransactionRepository) EXPECT() *MockTransactionRepositoryMockRecorder {
	return m.recorder
}

func (m *MockTransactionRepository) CreateTransaction(ctx context.Context, tx pgx.Tx, userID int, amount float64, targetUserID *int) error {
	ret := m.ctrl.Call(m, "CreateTransaction", ctx, tx, userID, amount, targetUserID)
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *MockTransactionRepositoryMockRecorder) CreateTransaction(ctx, tx, userID, amount, targetUserID interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTransaction", reflect.TypeOf((*MockTransactionRepository)(nil).CreateTransaction), ctx, tx, userID, amount, targetUserID)
}

func (m *MockTransactionRepository) GetLastTransactions(ctx context.Context, userID int) ([]models.Transaction, error) {
	ret := m.ctrl.Call(m, "GetLastTransactions", ctx, userID)
	ret0, _ := ret[0].([]models.Transaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockTransactionRepositoryMockRecorder) GetLastTransactions(ctx, userID interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLastTransactions", reflect.TypeOf((*MockTransactionRepository)(nil).GetLastTransactions), ctx, userID)
}
