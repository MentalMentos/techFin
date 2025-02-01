package transaction

import (
	"context"
	"github.com/MentalMentos/techFin/internal/clients/db"
	"github.com/MentalMentos/techFin/internal/clients/db/pg"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
)

type manager struct {
	db db.DB
}

func NewTransactionManager(db db.DB) db.TxManager {
	return &manager{
		db: db,
	}
}

func (m *manager) transaction(ctx context.Context, opts pgx.TxOptions, fn db.Handler) (err error) {
	// Проверка на существующую транзакцию в контексте
	tx, ok := ctx.Value(pg.TxKey).(pgx.Tx)
	if ok {
		// Если транзакция уже существует, выполняем обработчик в текущей транзакции
		return fn(ctx)
	}

	// Инициализация новой транзакции
	tx, err = m.db.BeginTx(ctx, opts)
	if err != nil {
		return errors.Wrap(err, "can't begin transaction")
	}

	// Добавляем транзакцию в контекст
	ctx = pg.MakeContextTx(ctx, tx)

	// Отсроченная функция для управления откатом и коммитом транзакции
	defer func() {
		if r := recover(); r != nil {
			err = errors.Errorf("panic recovered: %v", r)
		}

		// Откат транзакции в случае ошибки
		if err != nil {
			if errRollback := tx.Rollback(ctx); errRollback != nil {
				err = errors.Wrapf(err, "errRollback: %v", errRollback)
			}
			return
		}

		// Коммит транзакции в случае успеха
		if nil == err {
			err = tx.Commit(ctx)
			if err != nil {
				err = errors.Wrap(err, "tx commit failed")
			}
		}
	}()

	// Выполняем код в транзакции
	if err = fn(ctx); err != nil {
		err = errors.Wrap(err, "failed executing code inside transaction")
	}

	return err
}

func (m *manager) ReadCommitted(ctx context.Context, f db.Handler) error {
	txOpts := pgx.TxOptions{IsoLevel: pgx.ReadCommitted}
	return m.transaction(ctx, txOpts, f)
}

func (m *manager) RepeatableRead(ctx context.Context, f db.Handler) error {
	txOpts := pgx.TxOptions{IsoLevel: pgx.RepeatableRead}
	return m.transaction(ctx, txOpts, f)
}

func (m *manager) Serializable(ctx context.Context, f db.Handler) error {
	txOpts := pgx.TxOptions{IsoLevel: pgx.Serializable}
	return m.transaction(ctx, txOpts, f)
}
