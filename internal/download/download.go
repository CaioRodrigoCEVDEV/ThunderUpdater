package download

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func Download(ctx context.Context, url, destPath string, callback ProgressCallback) error {
	dir := filepath.Dir(destPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("erro ao criar diretório %s: %w", dir, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("erro ao criar requisição: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if ctx.Err() != nil {
			return ErrCancelled
		}
		return fmt.Errorf("%w: %w", ErrNoConnection, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return ErrNotFound
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status HTTP inesperado: %d", resp.StatusCode)
	}

	total := resp.ContentLength

	f, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("erro ao criar arquivo %s: %w", destPath, err)
	}

	body := &contextReader{ctx: ctx, reader: resp.Body}
	progress := NewProgressReader(body, total, callback)

	written, err := io.Copy(f, progress)
	f.Close()

	if err != nil {
		if ctx.Err() != nil {
			return ErrCancelled
		}
		return fmt.Errorf("%w: %w", ErrDownloadAborted, err)
	}

	if total > 0 && written != total {
		return fmt.Errorf("%w: esperado %d bytes, recebido %d", ErrSizeMismatch, total, written)
	}

	if callback != nil {
		progress.reportFinal()
	}

	return nil
}

type contextReader struct {
	ctx    context.Context
	reader io.Reader
}

func (r *contextReader) Read(p []byte) (int, error) {
	select {
	case <-r.ctx.Done():
		return 0, r.ctx.Err()
	default:
	}
	return r.reader.Read(p)
}

func ResumeDownload(ctx context.Context, url, destPath string, existingSize int64, callback ProgressCallback) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("erro ao criar requisição: %w", err)
	}
	req.Header.Set("Range", fmt.Sprintf("bytes=%d-", existingSize))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if ctx.Err() != nil {
			return ErrCancelled
		}
		return fmt.Errorf("%w: %w", ErrNoConnection, err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusNotFound:
		return ErrNotFound
	case http.StatusOK:
		return ErrRangeNotSupported
	case http.StatusRequestedRangeNotSatisfiable:
		return ErrRangeNotSupported
	case http.StatusPartialContent:
	default:
		return fmt.Errorf("status HTTP inesperado: %d", resp.StatusCode)
	}

	f, err := os.OpenFile(destPath, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("erro ao abrir arquivo para resumo: %w", err)
	}
	defer f.Close()

	body := &contextReader{ctx: ctx, reader: resp.Body}
	progress := NewProgressReader(body, resp.ContentLength, callback)

	written, err := io.Copy(f, progress)
	f.Close()

	if err != nil {
		if ctx.Err() != nil {
			return ErrCancelled
		}
		return fmt.Errorf("%w: %w", ErrDownloadAborted, err)
	}

	if resp.ContentLength > 0 && written != resp.ContentLength {
		return fmt.Errorf("%w: esperado %d bytes do range, recebido %d", ErrSizeMismatch, resp.ContentLength, written)
	}

	return nil
}
