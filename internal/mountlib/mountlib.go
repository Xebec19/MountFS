package mountlib

import (
	"fmt"
	"os"

	"bazil.org/fuse"
	ffs "bazil.org/fuse/fs"
	"github.com/Xebec19/reimagined-lamp/internal/fs"
)

func Mount(mountPath string, remote fs.Remote, verbose bool) error {

	if err := os.MkdirAll(mountPath, 0755); err != nil {
		return fmt.Errorf("failed to create mount point: %v", err)
	}

	c, err := fuse.Mount(mountPath, fuse.FSName("rclone"), fuse.Subtype("rclone"))

	if err != nil {
		return fmt.Errorf("failed to mount: %v", err)
	}
	defer c.Close()

	if verbose {
		fmt.Printf("Mounted %s at %s\n", "remote", mountPath)
		fmt.Println("Press Ctrl+C to unmount")
	}

	filesys := &fs.FS{Remote: remote}

	if err := ffs.Serve(c, filesys); err != nil {
		return fmt.Errorf("failed to serve filesystem: %v", err)
	}

	// <-c.Ready
	// if err := c.MountError; err != nil {
	// 	return fmt.Errorf("mount error: %v", err)
	// }
	// return nil
}
