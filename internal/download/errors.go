package download

import "errors"

var (
	ErrTimeout           = errors.New("timeout no download")
	ErrNotFound          = errors.New("arquivo não encontrado (404)")
	ErrNoConnection      = errors.New("sem conexão com o servidor")
	ErrFileNotFound      = errors.New("arquivo de destino não existe")
	ErrDownloadAborted   = errors.New("download interrompido")
	ErrCancelled         = errors.New("download cancelado")
	ErrInsufficientSpace = errors.New("espaço em disco insuficiente")
	ErrSizeMismatch      = errors.New("tamanho do arquivo difere do esperado")
	ErrRangeNotSupported = errors.New("servidor não suporta range")
)
