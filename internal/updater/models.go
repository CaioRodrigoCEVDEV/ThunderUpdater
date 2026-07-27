package updater

import "errors"

type UpdateResult struct {
	Success bool
	Version int
}

var (
	ErrExtractFailed = errors.New("falha na extração")
)
