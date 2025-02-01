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
	return &manager{db: db}
}

func (m *manager) transaction(ctx context.Context, opts pgx.TxOptions, fn db.TxHandler) (err error) {
	// Проверяем, есть ли уже активная транзакция в контексте
	if tx, ok := ctx.Value(pg.TxKey).(pgx.Tx); ok {
		// Если транзакция уже существует, выполняем код в её контексте
		return fn(ctx, tx)
	}

	// Создаём новую транзакцию
	tx, err := m.db.BeginTx(ctx, opts)
	if err != nil {
		return errors.Wrap(err, "can't begin transaction")
	}

	// Добавляем транзакцию в контекст
	ctx = pg.MakeContextTx(ctx, tx)

	// Отложенный откат/коммит
	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback(ctx)
			err = errors.Errorf("panic recovered: %v", r)
		}

		if err != nil {
			if errRollback := tx.Rollback(ctx); errRollback != nil {
				err = errors.Wrapf(err, "rollback failed: %v", errRollback)
			}
			return
		}

		if commitErr := tx.Commit(ctx); commitErr != nil {
			err = errors.Wrap(commitErr, "transaction commit failed")
		}
	}()

	// Выполняем переданную функцию внутри транзакции
	err = fn(ctx, tx)
	if err != nil {
		err = errors.Wrap(err, "failed executing code inside transaction")
	}

	return err
}

func (m *manager) ReadCommitted(ctx context.Context, f db.TxHandler) error {
	return m.transaction(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted}, f)
}

func (m *manager) RepeatableRead(ctx context.Context, f db.TxHandler) error {
	return m.transaction(ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead}, f)
}

func (m *manager) Serializable(ctx context.Context, f db.TxHandler) error {
	return m.transaction(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable}, f)
}
