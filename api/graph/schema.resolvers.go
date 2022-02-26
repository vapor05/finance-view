package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/vapor05/financeview/graph/generated"
	"github.com/vapor05/financeview/graph/model"
	"github.com/vapor05/financeview/pkg/expense"
)

func (r *mutationResolver) CreateExpense(ctx context.Context, input model.NewExpense) (*model.Expense, error) {
	ex, err := expense.SaveExpense(ctx, input, r.Db)
	if err != nil {
		return nil, fmt.Errorf("failed to save input new expense, %w", err)
	}
	return &ex, nil
}

func (r *queryResolver) Expenses(ctx context.Context) ([]*model.Expense, error) {
	exps, err := expense.ListExpenses(ctx, r.Db)
	if err != nil {
		return nil, fmt.Errorf("failed to get expenses, %w", err)
	}
	return exps, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
