package store

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
)

type Database struct {
	Conn *pgx.Conn
}

func NewDatabase(ctx context.Context, dbUrl string) (*Database, error) {
	conn, err := pgx.Connect(ctx, dbUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to app database, %w", err)
	}
	return &Database{conn}, nil
}

func (db *Database) GetDescriptionId(ctx context.Context, d string) (int, bool, error) {
	sql := `SELECT id FROM financeview.description WHERE description=$1`
	var id int
	if err := db.Conn.QueryRow(ctx, sql, d).Scan(&id); err != nil {
		if err == pgx.ErrNoRows {
			return 0, false, nil
		}
		return 0, false, fmt.Errorf("failed to query description table, %w", err)
	}
	return id, true, nil
}

func (db *Database) CreateDescription(ctx context.Context, d string) (int, error) {
	sql := `INSERT INTO financeview.description (description, createdate) VALUES ($1, $2) RETURNING id`
	var id int
	if err := db.Conn.QueryRow(ctx, sql, d, time.Now().UTC()).Scan(&id); err != nil {
		return 0, fmt.Errorf("failed to insert new description into database, %w", err)
	}
	return id, nil
}

func (db *Database) CreateExpense(ctx context.Context, dt time.Time, did int, amt float64, cmt string) (int, error) {
	sql := `INSERT INTO financeview.expense (date, description_id, amount, comment, createdate) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	var id int
	if err := db.Conn.QueryRow(
		ctx,
		sql,
		dt,
		did,
		amt,
		cmt,
		time.Now().UTC(),
	).Scan(&id); err != nil {
		return 0, fmt.Errorf("failed to insert new expense into database, %w", err)
	}
	return id, nil
}

func (db *Database) GetCategoryId(ctx context.Context, c string) (int, bool, error) {
	sql := `SELECT id FROM financeview.category WHERE name=$1`
	var id int
	if err := db.Conn.QueryRow(ctx, sql, c).Scan(&id); err != nil {
		if err == pgx.ErrNoRows {
			return 0, false, nil
		}
		return 0, false, fmt.Errorf("failed to query database for category, %w", err)
	}
	return id, true, nil
}

func (db *Database) CreateCategory(ctx context.Context, c string) (int, error) {
	sql := `INSERT INTO financeview.category (name, createdate) VALUES ($1, $2) RETURNING id`
	var id int
	if err := db.Conn.QueryRow(ctx, sql, c, time.Now().UTC()).Scan(&id); err != nil {
		return 0, fmt.Errorf("failed to insert new category into database, %w", err)
	}
	return id, nil
}

func (db *Database) LinkExpenseCategory(ctx context.Context, eid int, cid int) (int, error) {
	sql := `INSERT INTO financeview.expense_category (expense_id, category_id, createdate) VALUES ($1, $2, $3) RETURNING id`
	var id int
	if err := db.Conn.QueryRow(ctx, sql, eid, cid, time.Now().UTC()).Scan(&id); err != nil {
		return 0, fmt.Errorf("failed to insert new expense_category into database, %w", err)
	}
	return id, nil
}
