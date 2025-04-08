package repository

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/oklog/ulid/v2"
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

func (r *SqlxRepository[T]) FindByID(id string) (T, error) {
	var item T
	err := r.DB.Get(&item, fmt.Sprintf("SELECT * FROM %s WHERE id = ?", r.TableName), id)
	return item, err
}

func (r *SqlxRepository[T]) Create(item T) error {
	fields := strings.Join(r.Fields, ", ")
	values := strings.Join(r.Fields, ", :")

	dataMap, err := r.convertToMap(item)
	if err != nil {
		return err
	}

	id := ulid.Make()

	dataMap["id"] = id.String()
	dataMap["created_at"] = time.Now().Format("2006-01-02 15:04:05")
	dataMap["updated_at"] = time.Now().Format("2006-01-02 15:04:05")

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (:%s)", r.TableName, fields, values)

	_, err = r.DB.NamedExec(query, dataMap)
	return err
}

func (r *SqlxRepository[T]) convertToMap(item T) (map[string]interface{}, error) {
	data, err := json.Marshal(item)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *SqlxRepository[T]) Update(id string, item T) error {
	setClauses := generateQueryFields(r.Fields)

	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = :id", r.TableName, setClauses)

	dataMap, err := r.convertToMap(item)
	if err != nil {
		return err
	}

	dataMap["updated_at"] = time.Now().Format("2006-01-02 15:04:05")

	_, err = r.DB.NamedExec(query, dataMap)

	return err
}

func (r *SqlxRepository[T]) Delete(id string) error {
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
