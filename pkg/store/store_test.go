package store

import (
	"context"
	"log"
	"strconv"
	"testing"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/assert"
)

const dbUrl = "postgres://postgres:testing@localhost:5432/financeview"

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
