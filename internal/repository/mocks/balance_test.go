package repo_mocks

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v4"
	"reflect"
)

type MockBalanceRepository struct {
	ctrl     *gomock.Controller
	recorder *MockBalanceRepositoryMockRecorder
}

type MockBalanceRepositoryMockRecorder struct {
	mock *MockBalanceRepository
}

func NewMockBalanceRepository(ctrl *gomock.Controller) *MockBalanceRepository {
	mock := &MockBalanceRepository{ctrl: ctrl}
	mock.recorder = &MockBalanceRepositoryMockRecorder{mock}
	return mock
}

func (m *MockBalanceRepository) EXPECT() *MockBalanceRepositoryMockRecorder {
	return m.recorder
}

func (m *MockBalanceRepository) GetBalance(ctx context.Context, userID int) (float64, error) {
	ret := m.ctrl.Call(m, "GetBalance", ctx, userID)
	ret0, _ := ret[0].(float64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockBalanceRepositoryMockRecorder) GetBalance(ctx, userID interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBalance", reflect.TypeOf((*MockBalanceRepository)(nil).GetBalance), ctx, userID)
}

func (m *MockBalanceRepository) UpdateBalance(ctx context.Context, tx pgx.Tx, userID int, amount float64) (float64, error) {
	ret := m.ctrl.Call(m, "UpdateBalance", ctx, tx, userID, amount)
	ret0, _ := ret[0].(float64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockBalanceRepositoryMockRecorder) UpdateBalance(ctx, tx, userID, amount interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateBalance", reflect.TypeOf((*MockBalanceRepository)(nil).UpdateBalance), ctx, tx, userID, amount)
}
