package fs

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
)

type FileInfo struct {
	Name     string
	Size     int64
	ModeTime time.Time
	IsDir    bool
}

type Remote interface {
	List(ctx context.Context, path string) ([]FileInfo, error)
	Get(ctx context.Context, path string) (io.ReadCloser, error)
	Put(ctx context.Context, path string, data io.Reader) error
	Delete(ctx context.Context, path string) error
	Stat(ctx context.Context, path string) (FileInfo, error)
}

type FS struct {
	Remote Remote
}

type Dir struct {
	fs   *FS
	path string
}

type File struct {
	fs   *FS
	path string
	size int64
}

func (f *FS) Root() (fs.Node, error) {
	return &Dir{
		fs:   f,
		path: "",
	}, nil
}

func (d *Dir) Attr(ctx context.Context, a *fuse.Attr) error {
	a.Mode = os.ModeDir | 0755
	return nil
}

func (d *Dir) Lookup(ctx context.Context, name string) (fs.Node, error) {
	path := filepath.Join(d.path, name)

	info, err := d.fs.Remote.Stat(ctx, path)
	if err != nil {
		return nil, fuse.Errno(syscall.ENOENT)
	}

	if info.IsDir {
		return &Dir{fs: d.fs, path: path}, nil
	}

	return &File{
		fs:   d.fs,
		path: path,
		size: info.Size,
	}, nil
}

func (d *Dir) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	files, err := d.fs.Remote.List(ctx, d.path)
	if err != nil {
		return nil, err
	}

	var entries []fuse.Dirent
	for _, file := range files {
		typ := fuse.DT_File
		if file.IsDir {
			typ = fuse.DT_Dir
		}

		entries = append(entries, fuse.Dirent{
			Name: file.Name,
			Type: typ,
		})
	}
	return entries, nil
}

func (f *File) Attr(ctx context.Context, a *fuse.Attr) error {
	info, err := f.fs.Remote.Stat(ctx, f.path)
	if err != nil {
		a.Size = uint64(f.size)
		a.Mode = 0644
		return nil
	}

	a.Size = uint64(info.Size)
	a.Mode = 0644
	a.Mtime = info.ModeTime
	return nil
}

func (f *File) Open(ctx context.Context, req *fuse.OpenRequest, resp *fuse.OpenResponse) (fs.Handle, error) {
	return &FileHandle{
		file: f,
	}, nil
}

type FileHandle struct {
	file    *File
	reader  io.ReadCloser
	writing bool
	data    []byte
}

func (fh *FileHandle) Read(ctx context.Context, req *fuse.ReadRequest, resp *fuse.ReadResponse) error {
	if fh.reader == nil {
		var err error
		fh.reader, err = fh.file.fs.Remote.Get(ctx, fh.file.path)
		if err != nil {
			return err
		}
	}

	buf := make([]byte, req.Size)
	n, err := io.ReadFull(fh.reader, buf)
	if err != nil && err != io.ErrUnexpectedEOF && err != io.EOF {
		return err
	}

	resp.Data = buf[:n]
	return nil
}

func (fh *FileHandle) Write(ctx context.Context, req *fuse.WriteRequest, resp *fuse.WriteResponse) error {
	if !fh.writing {
		return fuse.EPERM
	}

	end := req.Offset + int64(len(req.Data))
	if end > int64(len(fh.data)) {
		newData := make([]byte, end)
		copy(newData, fh.data)
		fh.data = newData
	}

	copy(fh.data[req.Offset:], req.Data)
	resp.Size = len(req.Data)
	return nil
}

// Flush writes the data to remote
func (fh *FileHandle) Flush(ctx context.Context, req *fuse.FlushRequest) error {
	if fh.writing && len(fh.data) > 0 {
		return fh.file.fs.Remote.Put(ctx, fh.file.path, strings.NewReader(string(fh.data)))
	}
	return nil
}

func (fh *FileHandle) Release(ctx context.Context, req *fuse.ReleaseRequest) error {
	if fh.reader != nil {
		fh.reader.Close()
	}

	return nil
}
