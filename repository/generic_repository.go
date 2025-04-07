package repository

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

type SqlxRepository[T any] struct {
	DB        *sqlx.DB
	TableName string
	Fields    []string
}

func NewSqlxRepository[T any](
	db *sqlx.DB,
	table string,
	fields []string,
) *SqlxRepository[T] {
	return &SqlxRepository[T]{DB: db, TableName: table, Fields: fields}
}

func (r *SqlxRepository[T]) FindAll() ([]T, error) {
	var items []T
	err := r.DB.Select(&items, fmt.Sprintf("SELECT * FROM %s", r.TableName))
	return items, err
}

func (r *SqlxRepository[T]) FindByID(id int64) (T, error) {
	var item T
	err := r.DB.Get(&item, fmt.Sprintf("SELECT * FROM %s WHERE id = ?", r.TableName), id)
	return item, err
}

func (r *SqlxRepository[T]) Create(item T) error {
	fields := strings.Join(r.Fields, ", ")
	values := strings.Join(r.Fields, ", :")

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (:%s)", r.TableName, fields, values)

	_, err := r.DB.NamedExec(query, item)
	return err
}

func (r *SqlxRepository[T]) Update(id int64, item T) error {
	setClauses := generateQueryFields(r.Fields)

	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = :id", r.TableName, setClauses)

	_, err := r.DB.NamedExec(query, item)

	return err
}

func (r *SqlxRepository[T]) Delete(id int64) error {
	_, err := r.DB.Exec(fmt.Sprintf("DELETE FROM %s WHERE id = ?", r.TableName), id)
	return err
}

func generateQueryFields(fields []string) string {
	var setClauses []string
	for _, f := range fields {
		if f == "id" {
			continue
		}

		setClauses = append(setClauses, fmt.Sprintf("%s = :%s", f, f))
	}

	return strings.Join(setClauses, ", ")
}
