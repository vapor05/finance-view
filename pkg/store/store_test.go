package store

import (
	"context"
	"log"
	"testing"

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
	var want int
	if err = conn.QueryRow(context.TODO(), "select id from financeview.description where description=$1", desc).Scan(&want); err != nil {
		t.Fatalf("failed to get created id from db, %v", err)
	}
	assert.Equal(t, want, actual)
}
