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

func (r *SqlxRepository[T]) FindAll(query string, filtersQuery map[string]interface{}) ([]T, error) {
	var items []T

	finalQuery := fmt.Sprintf("SELECT * FROM %s %s", r.TableName, query)
	rows, err := r.DB.NamedQuery(finalQuery, filtersQuery)

	if err != nil {
		return items, err
	}
	defer rows.Close()

	for rows.Next() {
		var item T

		err = rows.StructScan(&item)
		if err != nil {
			break
		}

		items = append(items, item)
	}

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
