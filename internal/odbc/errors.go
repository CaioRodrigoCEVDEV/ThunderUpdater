package odbc

import "errors"

var (
	ErrDSNNotFound         = errors.New("dsn não encontrado")
	ErrConnectionFailed    = errors.New("falha na conexão")
	ErrInvalidQuery        = errors.New("consulta inválida")
	ErrColumnNotFound      = errors.New("campo inexistente no resultado")
	ErrDatabaseUnavailable = errors.New("banco de dados indisponível")
	ErrTimeout             = errors.New("timeout na operação")
)
