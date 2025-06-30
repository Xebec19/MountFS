package fs

import "time"

type FileInfo struct {
	Name     string
	Size     int64
	ModeTime time.Time
	IsDir    bool
}
