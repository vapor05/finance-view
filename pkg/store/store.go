package store

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/vapor05/financeview/graph/model"
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

func moneyToFloat(m string) (float64, error) {
	amt, err := strconv.ParseFloat(strings.ReplaceAll(m, "$", ""), 64)
	if err != nil {
		return 0, fmt.Errorf("failed to convert money amount to float, %v", err)
	}
	return amt, nil
}

func (db *Database) ListAllExpenses(ctx context.Context) ([]model.Expense, error) {
	expSql := `
		SELECT e.id, e.date, d.description, e.amount, e.comment
		FROM financeview.expense AS e
		INNER JOIN financeview.description AS d
		ON e.description_id = d.id
	`
	var exps []model.Expense
	rows, err := db.Conn.Query(ctx, expSql)
	if err != nil {
		return exps, fmt.Errorf("failed to select expenses from database, %w", err)
	}
	for rows.Next() {
		var e Expense
		if err := rows.Scan(&e.Id, &e.Date, &e.Description, &e.Amount, &e.Comment); err != nil {
			if err == pgx.ErrNoRows {
				return exps, nil
			}
			return exps, fmt.Errorf("failed to scan response from database, %w", err)
		}
		amt, err := moneyToFloat(e.Amount.String)
		if err != nil {
			return exps, fmt.Errorf("failed to covert amount, %w", err)
		}
		exps = append(exps, model.Expense{
			Id:          int(e.Id.Int),
			Date:        e.Date.Time.Format("01-02-2006"),
			Description: e.Description.String,
			Amount:      amt,
			Comment:     e.Comment.String,
		})
	}
	for i := range exps {
		cats, err := GetCategories(ctx, exps[i].Id, db)
		if err != nil {
			return exps, fmt.Errorf("failed to get expense's categories from database, %w", err)
		}
		exps[i].Categories = cats
	}
	return exps, nil
}

func GetCategories(ctx context.Context, eid int, db *Database) ([]model.Category, error) {
	catSql := `
		SELECT c.id, c.name
		FROM financeview.category AS c
		INNER JOIN financeview.expense_category AS ec
		ON c.id = ec.category_id AND ec.expense_id = $1
	`
	var cats []model.Category
	rows, err := db.Conn.Query(ctx, catSql, eid)
	if err != nil {
		return cats, fmt.Errorf("failed to select categories for expense_id=%v from database, %w", eid, err)
	}
	for rows.Next() {
		var c Category
		if err := rows.Scan(&c.Id, &c.Name); err != nil {
			if err == pgx.ErrNoRows {
				return cats, nil
			}
			return cats, fmt.Errorf("failed to scan categories for expense_id=%v from database, %w", eid, err)
		}
		cats = append(cats, model.Category{
			Id:   int(c.Id.Int),
			Name: c.Name.String,
		})
	}
	return cats, nil
}

type Expense struct {
	Id          pgtype.Int4
	Date        pgtype.Date
	Description pgtype.Text
	Amount      pgtype.Text
	Comment     pgtype.Text
}

type Category struct {
	Id   pgtype.Int4
	Name pgtype.Text
}
