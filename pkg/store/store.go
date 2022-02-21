package store

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
)

type Database struct {
	Conn *pgx.Conn
}

// type Database interface {
// 	GetDescriptionId(context.Context, string) (int, bool, error)
// 	CreateDescription(context.Context, string) (int, error)
// 	CreateExpense(context.Context, time.Time, int, float64, string) (int, error)
// 	GetCategoryId(context.Context, string) (int, bool, error)
// 	CreateCategory(context.Context, string) (int, error)
// 	LinkExpenseCategory(context.Context, int, int) (int, error)
// }

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
