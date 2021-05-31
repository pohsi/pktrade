package dbconnection

import (
	"os"
	"testing"

	dbx "github.com/go-ozzo/ozzo-dbx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

const DSN = "postgres://127.0.0.1/pktrade?sslmode=disable&user=postgres&password=postgres"

func runDBTest(t *testing.T, f func(db *dbx.DB)) {
	dsn, ok := os.LookupEnv("APP_DSN")
	if !ok {
		dsn = DSN
	}

	db, err := dbx.MustOpen("postgres", dsn)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer db.Close()

	sqls := []string{
		"CREATE TABLE IF NOT EXISTS dbconnectiontest (id VARCHAR PRIMARY KEY, name VARCHAR)",
		"TRUNCATE dbconnectiontest",
	}

	for _, s := range sqls {
		_, err = db.NewQuery(s).Execute()
		if err != nil {
			t.Error(err, " with SQL: ", s)
			t.FailNow()
		}
	}
	f(db)
}

func TestNew(t *testing.T) {

	runDBTest(t, func(db *dbx.DB) {
		dbc := New(db)
		assert.NotNil(t, dbc)
		assert.Equal(t, db, dbc.DB())
	})

}

// func TestNewHandeler(t *testing.T) {

// }
