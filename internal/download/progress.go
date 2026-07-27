package download

import (
	"io"
	"sync"
	"time"
)

const (
	progressMinInterval = 200 * time.Millisecond
	speedSmoothFactor   = 0.3
)

type ProgressInfo struct {
	BytesDownloaded int64
	BytesTotal      int64
	Percent         float64
	SpeedMBps       float64
	ETA             time.Duration
}

type ProgressCallback func(info ProgressInfo)

type ProgressReader struct {
	reader      io.Reader
	total       int64
	downloaded  int64
	callback    ProgressCallback
	start       time.Time
	lastCall    time.Time
	lastBytes   int64
	avgSpeed    float64
	minInterval time.Duration
	mu          sync.Mutex
}

func NewProgressReader(reader io.Reader, total int64, callback ProgressCallback) *ProgressReader {
	now := time.Now()
	r := &ProgressReader{
		reader:      reader,
		total:       total,
		callback:    callback,
		start:       now,
		lastCall:    now,
		lastBytes:   0,
		avgSpeed:    0,
		minInterval: progressMinInterval,
	}

	if callback != nil && total > 0 {
		callback(ProgressInfo{
			BytesTotal: total,
			Percent:    0,
		})
	}

	return r
}

func (r *ProgressReader) Read(p []byte) (int, error) {
	n, err := r.reader.Read(p)
	if n > 0 {
		r.mu.Lock()
		r.downloaded += int64(n)
		downloaded := r.downloaded
		r.mu.Unlock()

		r.tryReport(downloaded)
	}
	return n, err
}

func (r *ProgressReader) tryReport(downloaded int64) {
	now := time.Now()
	elapsed := now.Sub(r.start)

	if elapsed < r.minInterval {
		return
	}

	if now.Sub(r.lastCall) < r.minInterval {
		return
	}

	dt := now.Sub(r.lastCall).Seconds()
	bytesSinceLast := downloaded - r.lastBytes

	var speed float64
	if dt > 0 {
		instSpeed := float64(bytesSinceLast) / dt
		if r.avgSpeed == 0 {
			r.avgSpeed = instSpeed
		} else {
			r.avgSpeed = r.avgSpeed*(1-speedSmoothFactor) + instSpeed*speedSmoothFactor
		}
		speed = r.avgSpeed
	}

	r.lastCall = now
	r.lastBytes = downloaded

	var percent float64
	if r.total > 0 {
		percent = (float64(downloaded) / float64(r.total)) * 100
	}

	var eta time.Duration
	if speed > 0 && r.total > 0 {
		remaining := float64(r.total - downloaded)
		eta = time.Duration(remaining/speed) * time.Second
	}

	r.callback(ProgressInfo{
		BytesDownloaded: downloaded,
		BytesTotal:      r.total,
		Percent:         percent,
		SpeedMBps:       speed / (1024 * 1024),
		ETA:             eta,
	})
}

func (r *ProgressReader) reportFinal() {
	r.mu.Lock()
	downloaded := r.downloaded
	r.mu.Unlock()

	if r.callback == nil {
		return
	}

	elapsed := time.Since(r.start)

	var speed float64
	if elapsed.Seconds() > 0 {
		speed = float64(downloaded) / elapsed.Seconds()
	}

	r.callback(ProgressInfo{
		BytesDownloaded: downloaded,
		BytesTotal:      r.total,
		Percent:         100,
		SpeedMBps:       speed / (1024 * 1024),
	})
}
