package expense

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vapor05/financeview/graph/model"
)

type MockDatabase struct {
	desc map[int]string
	exp  map[int]struct {
		Id      int
		Date    time.Time
		Did     int
		Amount  float64
		Comment string
	}
	cat  map[int]string
	link map[int]struct {
		Id  int
		Eid int
		Cid int
	}
}

func (mdb *MockDatabase) GetDescriptionId(ctx context.Context, d string) (int, bool, error) {
	for k, v := range mdb.desc {
		if v == d {
			return k, true, nil
		}
	}
	return 0, false, nil
}

func (mdb *MockDatabase) CreateDescription(ctx context.Context, d string) (int, error) {
	id := rand.Int()
	mdb.desc[id] = d
	return id, nil
}

func (mdb *MockDatabase) CreateExpense(ctx context.Context, dt time.Time, did int, amt float64, cmt string) (int, error) {
	id := rand.Int()
	r := struct {
		Id      int
		Date    time.Time
		Did     int
		Amount  float64
		Comment string
	}{id, dt, did, amt, cmt}
	mdb.exp[id] = r
	return id, nil
}

func (mdb *MockDatabase) GetCategoryId(ctx context.Context, cat string) (int, bool, error) {
	for k, v := range mdb.cat {
		if v == cat {
			return k, true, nil
		}
	}
	return 0, false, nil
}

func (mdb *MockDatabase) CreateCategory(ctx context.Context, cat string) (int, error) {
	id := rand.Int()
	mdb.cat[id] = cat
	return id, nil
}

func (mdb *MockDatabase) LinkExpenseCategory(ctx context.Context, eid int, cid int) (int, error) {
	id := rand.Int()
	r := struct {
		Id  int
		Eid int
		Cid int
	}{id, eid, cid}
	mdb.link[id] = r
	return id, nil
}

func TestSaveExpense(t *testing.T) {
	t.Run("new and existing cat, existing desc", func(t *testing.T) {
		mock := MockDatabase{
			desc: map[int]string{2: "test desc"},
			cat:  map[int]string{5: "test cat"},
			exp: make(map[int]struct {
				Id      int
				Date    time.Time
				Did     int
				Amount  float64
				Comment string
			}),
			link: make(map[int]struct {
				Id  int
				Eid int
				Cid int
			}),
		}
		cmt := "test comment"
		input := model.NewExpense{
			Date:        "02-21-2022",
			Description: "test desc",
			Amount:      12.45,
			Categories:  []string{"test cat", "a new cat"},
			Comment:     &cmt,
		}
		want := model.Expense{
			Id:          -1,
			Date:        "02-21-2022",
			Description: "test desc",
			Amount:      12.45,
			Categories: []model.Category{
				{Id: 5, Name: "test cat"},
				{Id: -1, Name: "a new cat"},
			},
			Comment: cmt,
		}
		actual, err := SaveExpense(context.Background(), input, &mock)
		if err != nil {
			t.Fatalf("error running SaveExpense func, %v", err)
		}
		// id randomly assign
		assert.NotEqual(t, -1, actual.Id)
		want.Id = actual.Id
		var w int
		for i, c := range want.Categories {
			if c.Name == "a new cat" {
				w = i
			}
		}
		for i, c := range actual.Categories {
			if c.Name == "a new cat" {
				assert.NotEqual(t, -1, c.Id)
			}
			want.Categories[w].Id = actual.Categories[i].Id
		}
		assert.Equal(t, want, actual)
	})
	t.Run("existing cat, new desc", func(t *testing.T) {
		mock := MockDatabase{
			desc: make(map[int]string),
			cat:  map[int]string{5: "test cat"},
			exp: make(map[int]struct {
				Id      int
				Date    time.Time
				Did     int
				Amount  float64
				Comment string
			}),
			link: make(map[int]struct {
				Id  int
				Eid int
				Cid int
			}),
		}
		cmt := "test comment"
		input := model.NewExpense{
			Date:        "02-21-2022",
			Description: "test desc",
			Amount:      12.45,
			Categories:  []string{"test cat"},
			Comment:     &cmt,
		}
		want := model.Expense{
			Id:          -1,
			Date:        "02-21-2022",
			Description: "test desc",
			Amount:      12.45,
			Categories: []model.Category{
				{Id: 5, Name: "test cat"},
			},
			Comment: cmt,
		}
		actual, err := SaveExpense(context.Background(), input, &mock)
		if err != nil {
			t.Fatalf("error running SaveExpense func, %v", err)
		}
		// id randomly assign
		assert.NotEqual(t, -1, actual.Id)
		want.Id = actual.Id
		assert.Equal(t, want, actual)
	})
}
