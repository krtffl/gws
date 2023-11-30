package repository

import (
	"fmt"
	"strings"

	"github.com/krtffl/get-well-soon/internal/domain"
)

func getTableName(err error) string {
	return strings.Split(err.Error(), "\"")[1]
}

func getColumnAndTableNames(err error) (string, string) {
	column := strings.Split(err.Error(), "\"")[1]
	table := strings.Split(err.Error(), "\"")[3]
	return column, table
}

func handleErrors(err error) error {
	if strings.Contains(err.Error(), "ERROR: relation") &&
		strings.Contains(err.Error(), "does not exist (SQLSTATE 42P01)") {
		return errTableNotExistent(getTableName(err))
	}

	if strings.Contains(err.Error(), "ERROR: column") &&
		strings.Contains(err.Error(), "does not exist (SQLSTATE 42703)") {
		return errColumnNotExistent(getColumnAndTableNames(err))
	}

	if strings.Contains(err.Error(), "ERROR: duplicate key") &&
		strings.Contains(err.Error(), "(SQLSTATE 23505)") {
		return errDuplicateKey()
	}

	if strings.Contains(err.Error(), "ERROR: insert or update") &&
		strings.Contains(err.Error(), "(SQLSTATE 23503)") {
		return errForeignKeyConstraint()
	}

	if strings.Contains(err.Error(), "record not found") {
		return errNotFound()
	}

	if strings.Contains(err.Error(), "no rows in result set") {
		return errNotFound()
	}

	return errUnknown(err)
}

func errTableNotExistent(table string) error {
	return fmt.Errorf("%s: Table %s does not exist", domain.NonExistentTableError, table)
}

func errColumnNotExistent(column, table string) error {
	return fmt.Errorf(
		"%s: Column %s does not exist on table %s",
		domain.NonExistentColumnError,
		column,
		table,
	)
}

func errDuplicateKey() error {
	return fmt.Errorf("%s: Duplicate key", domain.DuplicateKeyError)
}

func errForeignKeyConstraint() error {
	return fmt.Errorf("%s: Foreign key violation", domain.ForeignKeyError)
}

func errNotFound() error {
	return fmt.Errorf("%s: Record not found", domain.NotFoundError)
}

func errUnknown(err error) error {
	return fmt.Errorf("%s: %v", domain.UnknownError, err)
}
