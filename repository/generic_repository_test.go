package repository

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

type TestModel struct {
	ID   string `db:"id"`
	Name string `db:"name"`
}

func setupMockDB(t *testing.T) (*sqlx.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	return sqlxDB, mock
}

func TestFindAll(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	filters := map[string]interface{}{"id": "01JW4MH8S671QVVGD0NYY1XWAP"}
	query := "WHERE id = :id"

	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow("01JW4MH8S671QVVGD0NYY1XWAP", "Item1").
		AddRow("01JW4MW2JXJQRQXCPP0T8EGPD0", "Item2")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM test_table WHERE id = ?")).
		WithArgs("01JW4MH8S671QVVGD0NYY1XWAP").
		WillReturnRows(rows)

	repo := NewSqlxRepository[TestModel](db, "test_table", []string{"name"})
	items, err := repo.FindAll(query, filters)

	assert.NoError(t, err)
	assert.Len(t, items, 2)
	assert.Equal(t, "01JW4MH8S671QVVGD0NYY1XWAP", items[0].ID)
	assert.Equal(t, "Item1", items[0].Name)
}

func TestFindByID(t *testing.T) {
	var id string = "01JW1A10MR50EPWW5QW7JKTFJE"

	db, mock := setupMockDB(t)
	defer db.Close()

	row := sqlmock.NewRows([]string{"id", "name"}).AddRow(id, "Item1")

	mock.ExpectQuery("SELECT \\* FROM test_table WHERE id = ?").
		WithArgs(id).
		WillReturnRows(row)

	repo := NewSqlxRepository[TestModel](db, "test_table", []string{"name"})
	item, err := repo.FindByID(id)

	assert.NoError(t, err)
	assert.Equal(t, "01JW1A10MR50EPWW5QW7JKTFJE", item.ID)
}

func TestCreate(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO test_table (name) VALUES (?)")).
		WithArgs("Test").
		WillReturnResult(sqlmock.NewResult(1, 1))

	repo := NewSqlxRepository[TestModel](db, "test_table", []string{"name"})
	err := repo.Create(TestModel{Name: "Test"})

	assert.NoError(t, err)
}

func TestUpdate(t *testing.T) {
	var id string = "01JW1A10MR50EPWW5QW7JKTFJE"

	db, mock := setupMockDB(t)
	defer db.Close()

	mock.ExpectExec(regexp.QuoteMeta("UPDATE test_table SET name = ? WHERE id = ?")).
		WithArgs("Updated", id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	repo := NewSqlxRepository[TestModel](db, "test_table", []string{"name"})
	err := repo.Update(id, TestModel{ID: id, Name: "Updated"})

	assert.NoError(t, err)
}

func TestDelete(t *testing.T) {
	var id string = "01JW1A10MR50EPWW5QW7JKTFJE"

	db, mock := setupMockDB(t)
	defer db.Close()

	mock.ExpectExec("DELETE FROM test_table WHERE id = ?").
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	repo := NewSqlxRepository[TestModel](db, "test_table", []string{"name"})
	err := repo.Delete(id)

	assert.NoError(t, err)
}
