package updater

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	extractMinInterval = 200 * time.Millisecond
)

type ExtractProgressInfo struct {
	FilesExtracted int64
	FilesTotal     int64
	BytesExtracted int64
	BytesTotal     int64
	Percent        float64
	CurrentFile    string
}

type ExtractProgressCallback func(info ExtractProgressInfo)

type ExtractManager struct{}

func newExtractManager() *ExtractManager {
	return &ExtractManager{}
}

func (m *ExtractManager) ExtractTo(zipPath, destDir string, callback ExtractProgressCallback) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrExtractFailed, err)
	}
	defer r.Close()

	var filesTotal int64
	var bytesTotal int64
	for _, f := range r.File {
		if !f.FileInfo().IsDir() {
			filesTotal++
			bytesTotal += int64(f.UncompressedSize64)
		}
	}

	if callback != nil && filesTotal > 0 {
		callback(ExtractProgressInfo{
			FilesTotal: filesTotal,
			BytesTotal: bytesTotal,
			Percent:    0,
		})
	}

	var mu sync.Mutex
	var filesExtracted int64
	var bytesExtracted int64
	lastReport := time.Now()
	var lastPercent int = -1

	for _, f := range r.File {
		if f.FileInfo().IsDir() {
			if err := extractFile(f, destDir); err != nil {
				return fmt.Errorf("%w: %w", ErrExtractFailed, err)
			}
			mu.Lock()
			filesExtracted++
			mu.Unlock()
			continue
		}

		if err := extractFile(f, destDir); err != nil {
			return fmt.Errorf("%w: %w", ErrExtractFailed, err)
		}

		mu.Lock()
		filesExtracted++
		bytesExtracted += int64(f.UncompressedSize64)
		currentFile := f.Name
		mu.Unlock()

		if callback != nil {
			now := time.Now()
			if now.Sub(lastReport) >= extractMinInterval {
				var percent float64
				if bytesTotal > 0 {
					percent = (float64(bytesExtracted) / float64(bytesTotal)) * 100
				}

				pctInt := int(percent)
				if pctInt != lastPercent {
					lastPercent = pctInt
					lastReport = now
					callback(ExtractProgressInfo{
						FilesExtracted: filesExtracted,
						FilesTotal:     filesTotal,
						BytesExtracted: bytesExtracted,
						BytesTotal:     bytesTotal,
						Percent:        percent,
						CurrentFile:    currentFile,
					})
				}
			}
		}
	}

	if callback != nil {
		callback(ExtractProgressInfo{
			FilesExtracted: filesTotal,
			FilesTotal:     filesTotal,
			BytesExtracted: bytesTotal,
			BytesTotal:     bytesTotal,
			Percent:        100,
			CurrentFile:    "",
		})
	}

	return nil
}

func extractFile(f *zip.File, destDir string) error {
	destPath := filepath.Join(destDir, f.Name)

	if f.FileInfo().IsDir() {
		return os.MkdirAll(destPath, 0755)
	}

	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return err
	}

	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	out, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, rc)
	return err
}
