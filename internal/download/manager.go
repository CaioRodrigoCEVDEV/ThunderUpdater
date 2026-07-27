package download

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type DownloadManager struct {
	log      *slog.Logger
	destDir  string
	mu       sync.Mutex
	callback ProgressCallback
}

func NewDownloadManager(log *slog.Logger) *DownloadManager {
	if log == nil {
		log = slog.Default()
	}
	return &DownloadManager{
		log:      log,
		destDir:  filepath.Join(os.TempDir(), "ThunderUpdater", "downloads"),
		callback: nil,
	}
}

func (m *DownloadManager) SetCallback(cb ProgressCallback) {
	m.callback = cb
}

func (m *DownloadManager) Download(ctx context.Context, fileName, downloadURL string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	destPath := filepath.Join(m.destDir, fileName)

	var partialSize int64
	if fi, err := os.Stat(destPath); err == nil && fi.Size() > 0 {
		partialSize = fi.Size()
	}

	var lastErr error

	for attempt := 1; attempt <= 3; attempt++ {
		if attempt > 1 {
			m.log.Info(fmt.Sprintf("Tentativa %d de 3", attempt))
			time.Sleep(backoffDuration(attempt))
		}

		if partialSize > 0 {
			m.log.Info("Continuando download em",
				"tamanho", fmt.Sprintf("%.0f MB", float64(partialSize)/(1024*1024)),
			)

			cb := m.callback
			if cb == nil {
				cb = makeResumeCallback(m.log)
			}
			err := ResumeDownload(ctx, downloadURL, destPath, partialSize, cb)
			if err == nil {
				return verifyAndReturn(destPath, m.log)
			}

			if errors.Is(err, ErrRangeNotSupported) {
				m.log.Info("Servidor não suporta resume. Reiniciando download.")
				os.Remove(destPath)
				partialSize = 0
			} else {
				lastErr = err
				if !isRetryable(err) {
					return "", err
				}
				continue
			}
		}

		cb := m.callback
		if cb == nil {
			cb = makeFreshCallback(m.log)
		}
		err := Download(ctx, downloadURL, destPath, cb)
		if err == nil {
			return verifyAndReturn(destPath, m.log)
		}

		lastErr = err

		if fi, statErr := os.Stat(destPath); statErr == nil && fi.Size() > 0 {
			partialSize = fi.Size()
		}

		if !isRetryable(err) {
			return "", err
		}
	}

	return "", lastErr
}

func backoffDuration(attempt int) time.Duration {
	switch attempt {
	case 2:
		return 2 * time.Second
	default:
		return 5 * time.Second
	}
}

func isRetryable(err error) bool {
	if err == nil {
		return false
	}
	return errors.Is(err, ErrNoConnection) ||
		errors.Is(err, ErrDownloadAborted) ||
		errors.Is(err, ErrSizeMismatch) ||
		errors.Is(err, ErrTimeout)
}

func verifyAndReturn(destPath string, log *slog.Logger) (string, error) {
	fi, err := os.Stat(destPath)
	if err != nil {
		return "", ErrFileNotFound
	}

	log.Info("Download concluído.",
		"arquivo", destPath,
		"tamanho", fmt.Sprintf("%.0f MB", float64(fi.Size())/(1024*1024)),
	)

	return destPath, nil
}

func makeFreshCallback(log *slog.Logger) ProgressCallback {
	lastPercent := -1
	return func(info ProgressInfo) {
		if info.Percent == 0 && info.BytesTotal > 0 {
			totalMB := float64(info.BytesTotal) / (1024 * 1024)
			log.Info("Tamanho do arquivo",
				"total", fmt.Sprintf("%.0f MB", totalMB),
			)
		}

		p := int(info.Percent)
		if p > lastPercent && p%5 == 0 {
			lastPercent = p
			log.Info("Progresso",
				"percent", fmt.Sprintf("%d%%", p),
				"speed", fmt.Sprintf("%.1f MB/s", info.SpeedMBps),
			)
		}
	}
}

func makeResumeCallback(log *slog.Logger) ProgressCallback {
	lastPercent := -1
	return func(info ProgressInfo) {
		p := int(info.Percent)
		if p > lastPercent && p%5 == 0 {
			lastPercent = p
			log.Info("Progresso",
				"percent", fmt.Sprintf("%d%%", p),
				"speed", fmt.Sprintf("%.1f MB/s", info.SpeedMBps),
			)
		}
	}
}
