package pg

import (
	"context"
	"github.com/MentalMentos/techFin/internal/clients/db"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
)

const (
	dsnEnvName     = "PG_DSN"
	TxKey      key = "tx"
)

type key string

type pgClient struct {
	dbc *pgxpool.Pool
}

func New(ctx context.Context) (db.Client, error) {
	dsn := "host=localhost port=5432 dbname=postgres user=user password=1234 sslmode=disable"
	if len(dsn) == 0 {
		return nil, errors.New("pg dsn not found")
	}

	dbc, err := pgxpool.Connect(ctx, dsn)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to db")
	}

	return &pgClient{dbc: dbc}, nil
}

func (c *pgClient) DB() db.DB {
	return &pg{dbc: c.dbc}
}

func (c *pgClient) Close() error {
	if c.dbc != nil {
		c.dbc.Close()
	}
	return nil
}

// pg реализует интерфейс db.DB для работы с базой данных
type pg struct {
	dbc *pgxpool.Pool
}

func (p *pg) ScanOneContext(ctx context.Context, dest interface{}, q db.Query, args ...interface{}) error {
	row, err := p.QueryContext(ctx, q, args...)
	if err != nil {
		return err
	}
	return pgxscan.ScanOne(dest, row)
}

func (p *pg) ScanAllContext(ctx context.Context, dest interface{}, q db.Query, args ...interface{}) error {
	rows, err := p.QueryContext(ctx, q, args...)
	if err != nil {
		return err
	}
	return pgxscan.ScanAll(dest, rows)
}

func (p *pg) ExecContext(ctx context.Context, q db.Query, args ...interface{}) (pgconn.CommandTag, error) {
	tx, ok := ctx.Value(TxKey).(pgx.Tx)
	if ok {
		return tx.Exec(ctx, q.QueryRaw, args...)
	}
	return p.dbc.Exec(ctx, q.QueryRaw, args...)
}

func (p *pg) QueryContext(ctx context.Context, q db.Query, args ...interface{}) (pgx.Rows, error) {
	tx, ok := ctx.Value(TxKey).(pgx.Tx)
	if ok {
		return tx.Query(ctx, q.QueryRaw, args...)
	}
	return p.dbc.Query(ctx, q.QueryRaw, args...)
}

func (p *pg) QueryRowContext(ctx context.Context, q db.Query, args ...interface{}) pgx.Row {
	tx, ok := ctx.Value(TxKey).(pgx.Tx)
	if ok {
		return tx.QueryRow(ctx, q.QueryRaw, args...)
	}
	return p.dbc.QueryRow(ctx, q.QueryRaw, args...)
}

func (p *pg) BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error) {
	return p.dbc.BeginTx(ctx, txOptions)
}

func (p *pg) Ping(ctx context.Context) error {
	return p.dbc.Ping(ctx)
}

func (p *pg) Close() {
	p.dbc.Close()
}

func MakeContextTx(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, TxKey, tx)
}
