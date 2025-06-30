package fs

import (
	"context"
	"io"
	"time"
)

type FileInfo struct {
	Name     string
	Size     int64
	ModeTime time.Time
	IsDir    bool
}

type Remote struct {
	List   func(ctx context.Context, path string) ([]FileInfo, error)
	Get    func(ctx context.Context, path string) (io.ReadCloser, error)
	Put    func(ctx context.Context, path string, data io.Reader) error
	Delete func(ctx context.Context, path string) error
	Stat   func(ctx context.Context, path string) (FileInfo, error)
}
