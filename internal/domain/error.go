package domain

type (
	ErrorCode uint
	ErrorMsg  string
)

const (
	NonExistentTableError  ErrorMsg = "NonExistentPostgreSQLTableError"
	NonExistentColumnError ErrorMsg = "NonExistentPostgreSQLColumnError"
	DuplicateKeyError      ErrorMsg = "DuplicateKeyPostgreSQLError"
	ForeignKeyError        ErrorMsg = "ForeignKeyPostgreSQLError"
	NotFoundError          ErrorMsg = "NotFoundPostgreSQLError"
	UnknownError           ErrorMsg = "UnkwnownPostgreSQLError"
)
