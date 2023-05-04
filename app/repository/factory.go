package repository

import (
	"context"
	"fmt"

	"github.com/opensaucerer/barf/app/database"
	"github.com/opensaucerer/barf/app/global"
)

type Factor struct {
	Key   string `json:"key"`
	Value int64  `json:"value"`
}

// GenerateAccountNumber generates a new account number for a new account
// by shifting the cursor by the given step. Account numbers are linearly
// random.
func GenerateAccountNumber() (string, error) {
	if global.FactoryCursor == 0 {
		cursor, err := ShiftCursorForKey(global.FactoryStep, "account_number")
		if err != nil {
			return "", err
		}
		global.FactoryCursor = cursor
		global.FactoryPointer = cursor - global.FactoryStep
	}
	if global.FactoryPointer == global.FactoryCursor {
		cursor, err := ShiftCursorForKey(global.FactoryStep, "account_number")
		if err != nil {
			return "", err
		}
		global.FactoryCursor = cursor
		global.FactoryPointer = cursor - global.FactoryStep
	}
	global.FactoryPointer++
	return fmt.Sprintf("%010d", global.FactoryPointer), nil
}

// ShiftCursorForAccountNumber shifts the cursor by the given step if the field exists else it creates it.
func ShiftCursorForKey(step int64, key string) (int64, error) {
	query := `INSERT INTO factory (key, value) VALUES ($1, $2) ON CONFLICT (key) DO UPDATE SET value = factory.value + $2 RETURNING value`
	var cursor int64
	err := database.PostgreSQLDB.QueryRow(context.Background(), query, key, step).Scan(&cursor)
	if err != nil {
		return 0, err
	}
	return cursor, nil
}
