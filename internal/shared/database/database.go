package database

import (
	"database/sql"
	"fmt"
	"todo/internal/repositories/contracts"
)

type BaseDB struct {
}

func (d *BaseDB) RunSelectSQL(appConfig contracts.AppConfig, sql string, params ...any) ([]map[string]any, error) {
	db, err := d.ConnectToDB(appConfig)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query(sql, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	results := make([]map[string]any, 0)
	for rows.Next() {
		result := make(map[string]any, len(columns))

		values := make([]any, len(columns))
		pointers := make([]any, len(columns))

		for i := range columns {
			pointers[i] = &values[i]
		}

		err := rows.Scan(pointers...)
		if err != nil {
			return nil, err
		}

		for i, fieldName := range columns {
			result[fieldName] = values[i]
		}

		results = append(results, result)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (s *BaseDB) ExecuteSQL(appConfig contracts.AppConfig, sql string, params ...any) error {
	db, err := s.ConnectToDB(appConfig)
	if err != nil {
		return err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			err = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	_, err = tx.Exec(sql, params...)

	return err
}

func (d *BaseDB) ConnectToDB(appConfig contracts.AppConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		appConfig.DBUsername(),
		appConfig.DBPassword(),
		appConfig.DBHost(),
		appConfig.DBPort(),
		appConfig.DBDatabase(),
	)

	return sql.Open("mysql", dsn)
}
