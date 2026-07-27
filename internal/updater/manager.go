package updater

import (
	"log/slog"

	"github.com/caio/ThunderUpdaterGO/internal/config"
)

type UpdateManager struct {
	log        *slog.Logger
	version    int
	extractMgr *ExtractManager
}

func NewUpdateManager(log *slog.Logger, version int) *UpdateManager {
	if log == nil {
		log = slog.Default()
	}
	return &UpdateManager{
		log:        log,
		version:    version,
		extractMgr: newExtractManager(),
	}
}

func (m *UpdateManager) Update(zipPath string, callback ExtractProgressCallback) (*UpdateResult, error) {
	thunderDir := config.ThunderPath

	if err := m.extractMgr.ExtractTo(zipPath, thunderDir, callback); err != nil {
		return nil, err
	}

	result := &UpdateResult{
		Success: true,
		Version: m.version,
	}

	return result, nil
}
