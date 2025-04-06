package repository

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

type TestModel struct {
	ID   int64  `db:"id"`
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

	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(1, "Item1").
		AddRow(2, "Item2")

	mock.ExpectQuery("SELECT \\* FROM test_table").
		WillReturnRows(rows)

	repo := NewSqlxRepository[TestModel](db, "test_table", []string{"id", "name"})
	items, err := repo.FindAll()

	assert.NoError(t, err)
	assert.Len(t, items, 2)
	assert.Equal(t, int64(1), items[0].ID)
}

func TestFindByID(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	row := sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Item1")

	mock.ExpectQuery("SELECT \\* FROM test_table WHERE id = ?").
		WithArgs(1).
		WillReturnRows(row)

	repo := NewSqlxRepository[TestModel](db, "test_table", []string{"id", "name"})
	item, err := repo.FindByID(1)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), item.ID)
}

func TestCreate(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO test_table (name) VALUES (?)")).
		WillReturnResult(sqlmock.NewResult(1, 1))

	repo := NewSqlxRepository[TestModel](db, "test_table", []string{"name"})
	err := repo.Create(TestModel{Name: "Test"})

	assert.NoError(t, err)
}

func TestUpdate(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	mock.ExpectExec(regexp.QuoteMeta("UPDATE test_table SET id = ?, name = ? WHERE id = ?")).
		WillReturnResult(sqlmock.NewResult(1, 1))

	repo := NewSqlxRepository[TestModel](db, "test_table", []string{"id", "name"})
	err := repo.Update(1, TestModel{ID: 1, Name: "Updated"})

	assert.NoError(t, err)
}

func TestDelete(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	mock.ExpectExec("DELETE FROM test_table WHERE id = ?").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	repo := NewSqlxRepository[TestModel](db, "test_table", []string{"id", "name"})
	err := repo.Delete(1)

	assert.NoError(t, err)
}
