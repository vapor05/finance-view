package store

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"testing"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/assert"
	"github.com/vapor05/financeview/graph/model"
)

const dbUrl = "postgres://postgres:testing@localhost:5432/postgres"

var conn *pgx.Conn

func init() {
	var err error
	conn, err = pgx.Connect(context.Background(), dbUrl)
	if err != nil {
		log.Fatalf("failed to connect to test database, %v", err)
	}
}

func TestGetDescriptionId(t *testing.T) {
	db := Database{conn}
	t.Run("desc doesn't exist", func(t *testing.T) {
		_, ok, err := db.GetDescriptionId(context.Background(), "does not exist")
		if err != nil {
			t.Fatalf("error running GetDescriptionId func, %v", err)
		}
		assert.False(t, ok)
	})
	t.Run("desc does exist", func(t *testing.T) {
		var want int
		if err := conn.QueryRow(context.TODO(), "INSERT INTO financeview.description (description) VALUES ('test desc') RETURNING id").Scan(&want); err != nil {
			t.Errorf("error inserting test description data into db, %v", err)
		}
		defer func() {
			_, err := conn.Exec(context.TODO(), "TRUNCATE TABLE financeview.description")
			if err != nil {
				t.Fatalf("error cleaning up test data")
			}
		}()
		id, ok, err := db.GetDescriptionId(context.Background(), "test desc")
		if err != nil {
			t.Errorf("error running GetDescriptionId func, %v", err)
		}
		assert.True(t, ok)
		assert.Equal(t, want, id)
	})
}

func TestCreateDescription(t *testing.T) {
	db := Database{conn}
	desc := "a test description"
	actual, err := db.CreateDescription(context.Background(), desc)
	if err != nil {
		t.Fatalf("error running CreateDescription func, %v", err)
	}
	defer func() {
		_, err := conn.Exec(context.TODO(), "TRUNCATE TABLE financeview.description")
		if err != nil {
			t.Fatalf("error cleaning up test data")
		}
	}()
	var want int
	if err = conn.QueryRow(context.TODO(), "select id from financeview.description where description=$1", desc).Scan(&want); err != nil {
		t.Fatalf("failed to get created id from db, %v", err)
	}
	assert.Equal(t, want, actual)
	var d string
	if err = conn.QueryRow(context.TODO(), "select description from financeview.description where id=$1", actual).Scan(&d); err != nil {
		t.Fatalf("failed to get created description from db, %v", err)
	}
	assert.Equal(t, desc, d)
}

func TestCreateExpense(t *testing.T) {
	db := Database{conn}
	dt := time.Date(2022, 02, 21, 0, 0, 0, 0, time.UTC)
	did := 5
	amt := 25.08
	cmt := "test comment"
	actual, err := db.CreateExpense(context.TODO(), dt, did, amt, cmt)
	if err != nil {
		t.Fatalf("error running CreateExpense func, %v", err)
	}
	defer func() {
		_, err := conn.Exec(context.TODO(), "TRUNCATE TABLE financeview.expense")
		if err != nil {
			t.Fatalf("error cleaning up test data")
		}
	}()
	var want int
	if err = conn.QueryRow(context.TODO(), "select id from financeview.expense where amount=$1", amt).Scan(&want); err != nil {
		t.Fatalf("failed to get created id from db, %v", err)
	}
	assert.Equal(t, want, actual)
	var adt time.Time
	var adid int
	var aamt string
	var acmt string
	sql := "select date, description_id, amount, comment from financeview.expense where id=$1"
	if err = conn.QueryRow(context.TODO(), sql, actual).Scan(&adt, &adid, &aamt, &acmt); err != nil {
		t.Fatalf("failed to get created expense from db, %v", err)
	}
	aamtf, err := strconv.ParseFloat(aamt[1:], 64)
	if err != nil {
		t.Fatalf("failed to convert db amount to float, %v", err)
	}
	assert.Equal(t, dt, adt)
	assert.Equal(t, adid, did)
	assert.Equal(t, aamtf, amt)
	assert.Equal(t, acmt, cmt)
}

func TestGetCategoryId(t *testing.T) {
	db := Database{conn}
	t.Run("cat doesn't exist", func(t *testing.T) {
		_, ok, err := db.GetCategoryId(context.Background(), "does not exist")
		if err != nil {
			t.Fatalf("error running GetCategoryId func, %v", err)
		}
		assert.False(t, ok)
	})
	t.Run("cat does exist", func(t *testing.T) {
		var want int
		if err := conn.QueryRow(context.TODO(), "INSERT INTO financeview.category (name) VALUES ('test category') RETURNING id").Scan(&want); err != nil {
			t.Errorf("error inserting test category data into db, %v", err)
		}
		defer func() {
			_, err := conn.Exec(context.TODO(), "TRUNCATE TABLE financeview.category")
			if err != nil {
				t.Fatalf("error cleaning up test data")
			}
		}()
		id, ok, err := db.GetCategoryId(context.Background(), "test category")
		if err != nil {
			t.Errorf("error running GetCategoryId func, %v", err)
		}
		assert.True(t, ok)
		assert.Equal(t, want, id)
	})
}

func TestCreateCategory(t *testing.T) {
	db := Database{conn}
	name := "test category"
	actual, err := db.CreateCategory(context.Background(), name)
	if err != nil {
		t.Fatalf("error running CreateCategory func, %v", err)
	}
	defer func() {
		_, err := conn.Exec(context.TODO(), "TRUNCATE TABLE financeview.category")
		if err != nil {
			t.Fatalf("error cleaning up test data")
		}
	}()
	var want int
	if err = conn.QueryRow(context.TODO(), "select id from financeview.category where name=$1", name).Scan(&want); err != nil {
		t.Fatalf("failed to get created id from db, %v", err)
	}
	assert.Equal(t, want, actual)
	var aname string
	if err = conn.QueryRow(context.TODO(), "select name from financeview.category where id=$1", actual).Scan(&aname); err != nil {
		t.Fatalf("failed to get name from db, %v", err)
	}
	assert.Equal(t, name, aname)
}

func TestLinkExpenseCategory(t *testing.T) {
	db := Database{conn}
	eid := 15
	cid := 8
	actual, err := db.LinkExpenseCategory(context.TODO(), eid, cid)
	if err != nil {
		t.Fatalf("error running LinkExpenseCategory func, %v", err)
	}
	defer func() {
		_, err := conn.Exec(context.TODO(), "TRUNCATE TABLE financeview.expense_category")
		if err != nil {
			t.Fatalf("error cleaning up test data")
		}
	}()
	var want int
	if err = conn.QueryRow(context.TODO(), "select id from financeview.expense_category where expense_id=$1", eid).Scan(&want); err != nil {
		t.Fatalf("error getting created id from db, %v", err)
	}
	assert.Equal(t, actual, want)
	var aeid, acid int
	if err = conn.QueryRow(context.TODO(), "select expense_id, category_id from financeview.expense_category where id=$1", actual).Scan(&aeid, &acid); err != nil {
		t.Fatalf("error getting created data from db, %v", err)
	}
	assert.Equal(t, eid, aeid)
	assert.Equal(t, cid, acid)
}

func TestNewDatabase(t *testing.T) {
	db, err := NewDatabase(context.Background(), dbUrl)
	if err != nil {
		t.Fatalf("error running NewDatabase func, %v", err)
	}
	defer db.Conn.Close(context.Background())
	err = db.Conn.Ping(context.Background())
	assert.Nil(t, err)
	_, err = NewDatabase(context.Background(), "not a database")
	assert.Error(t, err)
}

func Test_moneyToFloat(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  float64
	}{
		{name: "$45", input: "$45", want: 45.0},
		{name: "$105.35", input: "$105.35", want: 105.35},
		{name: "$0", input: "$0", want: 0},
		{name: "$0.0004", input: "$0.0004", want: 0.0004},
		{name: "-$5.68", input: "-$5.68", want: -5.68},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual, err := moneyToFloat(c.input)
			if err != nil {
				t.Fatalf("error running moneyToFloat func, %v", err)
			}
			assert.Equal(t, c.want, actual)
		})
	}
}

func cleanUpDb() error {
	_, err := conn.Exec(context.TODO(), "TRUNCATE TABLE financeview.description")
	if err != nil {
		return fmt.Errorf("error cleaning up test data, %w", err)
	}
	_, err = conn.Exec(context.TODO(), "TRUNCATE TABLE financeview.expense")
	if err != nil {
		return fmt.Errorf("error cleaning up test data, %w", err)
	}
	_, err = conn.Exec(context.TODO(), "TRUNCATE TABLE financeview.category")
	if err != nil {
		return fmt.Errorf("error cleaning up test data, %w", err)
	}
	_, err = conn.Exec(context.TODO(), "TRUNCATE TABLE financeview.expense_category")
	if err != nil {
		return fmt.Errorf("error cleaning up test data, %w", err)
	}
	return nil
}
func TestListAllExpenses(t *testing.T) {
	// create test data
	ctx := context.Background()
	desc := "test desc"
	var eid, did, cid int
	if err := conn.QueryRow(ctx, "INSERT INTO financeview.description (description) VALUES ($1) RETURNING id", desc).Scan(&did); err != nil {
		t.Fatalf("failed to setup test data, %v", err)
	}
	dt := time.Date(2022, time.February, 26, 0, 0, 0, 0, time.Local)
	amt := 105.65
	cmt := "test comment"
	if err := conn.QueryRow(ctx, "INSERT INTO financeview.expense (date, description_id, amount, comment) VALUES ($1,$2,$3,$4) RETURNING id", dt, did, amt, cmt).Scan(&eid); err != nil {
		t.Fatalf("failed to setup test data, %v", err)
	}
	cname := "test cat"
	if err := conn.QueryRow(ctx, "INSERT INTO financeview.category (name) VALUES ($1) RETURNING id", cname).Scan(&cid); err != nil {
		t.Fatalf("failed to setup test data, %v", err)
	}
	_, err := conn.Exec(ctx, "INSERT INTO financeview.expense_category (expense_id,category_id) VALUES ($1, $2)", eid, cid)
	if err != nil {
		t.Fatalf("failed to setup test data, %v", err)
	}
	defer func() {
		err := cleanUpDb()
		if err != nil {
			t.Fatalf("failed to clean up test data, %v", err)
		}
	}()
	want := []model.Expense{
		{
			Id:          eid,
			Date:        dt.Format("01-02-2006"),
			Description: desc,
			Amount:      amt,
			Categories: []model.Category{
				{Id: cid, Name: cname},
			},
			Comment: cmt,
		},
	}
	db := Database{conn}
	actual, err := db.ListAllExpenses(context.Background())
	if err != nil {
		t.Fatalf("error running ListAllExpenses func, %v", err)
	}
	assert.Equal(t, want, actual)

}
