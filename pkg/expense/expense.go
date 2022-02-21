package expense

import (
	"context"
	"fmt"
	"time"

	"github.com/vapor05/financeview/graph/model"
)

type Database interface {
	GetDescriptionId(context.Context, string) (int, bool, error)
	CreateDescription(context.Context, string) (int, error)
	CreateExpense(context.Context, time.Time, int, float64, string) (int, error)
	GetCategoryId(context.Context, string) (int, bool, error)
	CreateCategory(context.Context, string) (int, error)
	LinkExpenseCategory(context.Context, int, int) (int, error)
}

func SaveExpense(ctx context.Context, ne model.NewExpense, db Database) (model.Expense, error) {
	did, ok, err := db.GetDescriptionId(ctx, ne.Description)
	if err != nil {
		return model.Expense{}, fmt.Errorf("failed to get description_id for new expense, %w", err)
	}
	if !ok {
		did, err = db.CreateDescription(ctx, ne.Description)
		if err != nil {
			return model.Expense{}, fmt.Errorf("failed to create new description, %w", err)
		}
	}
	dt, err := time.Parse("01-02-2006", ne.Date)
	if err != nil {
		return model.Expense{}, fmt.Errorf("failed to parse new expense date, %w", err)
	}
	eid, err := db.CreateExpense(ctx, dt, did, ne.Amount, *ne.Comment)
	if err != nil {
		return model.Expense{}, fmt.Errorf("failed to save new expense data, %w", err)
	}
	var cats []model.Category
	for _, c := range ne.Categories {
		cid, ok, err := db.GetCategoryId(ctx, c)
		if err != nil {
			return model.Expense{}, fmt.Errorf("failed to get category_id, %w", err)
		}
		if !ok {
			cid, err = db.CreateCategory(ctx, c)
			if err != nil {
				return model.Expense{}, fmt.Errorf("failed to create new category, %w", err)
			}
		}
		_, err = db.LinkExpenseCategory(ctx, eid, cid)
		if err != nil {
			return model.Expense{}, fmt.Errorf("failed to link expense and category, %w", err)
		}
		cats = append(cats, model.Category{Id: cid, Name: c})
	}
	e := model.Expense{
		Id:          eid,
		Date:        dt.Format("01-02-2006"),
		Description: ne.Description,
		Amount:      ne.Amount,
		Categories:  cats,
		Comment:     *ne.Comment,
	}
	return e, nil
}
