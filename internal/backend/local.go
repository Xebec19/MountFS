package backend

import (
	"context"
	"io"
	"os"

	"path/filepath"
	"sync"

	"github.com/Xebec19/reimagined-lamp/internal/fs"
)

type LocalRemote struct {
	root string
	mu   sync.RWMutex
}

func NewLocalRemote(root string) *LocalRemote {
	return &LocalRemote{root: root}
}

func (r *LocalRemote) List(ctx context.Context, path string) ([]fs.FileInfo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	fullPath := filepath.Join(r.root, path)
	entries, err := os.ReadDir(fullPath)
	if err != nil {
		return nil, err
	}

	var files []fs.FileInfo
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}

		files = append(files, fs.FileInfo{
			Name:     entry.Name(),
			Size:     info.Size(),
			ModeTime: info.ModTime(),
			IsDir:    entry.IsDir(),
		})
	}

	return files, nil
}

func (r *LocalRemote) Get(ctx context.Context, path string) (io.ReadCloser, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	fullPath := filepath.Join(r.root, path)
	return os.Open(fullPath)
}

func (r *LocalRemote) Put(ctx context.Context, path string, data io.Reader) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	fullPath := filepath.Join(r.root, path)
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return err
	}

	file, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, data)
	return err
}

func (r *LocalRemote) Delete(ctx context.Context, path string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	fullPath := filepath.Join(r.root, path)
	return os.Remove(fullPath)
}

func (r *LocalRemote) Stat(ctx context.Context, path string) (fs.FileInfo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	fullPath := filepath.Join(r.root, path)
	info, err := os.Stat(fullPath)
	if err != nil {
		return fs.FileInfo{}, err
	}

	return fs.FileInfo{
		Name:     info.Name(),
		Size:     info.Size(),
		ModeTime: info.ModTime(),
		IsDir:    info.IsDir(),
	}, nil
}
