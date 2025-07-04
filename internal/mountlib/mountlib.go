package mountlib

import (
	"fmt"
	"os"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
)

func Mount(mountPath string, remote fs.FS, verbose bool) error {

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

	if err := fs.Serve(c, remote); err != nil {
		return fmt.Errorf("failed to serve filesystem: %v", err)
	}

	return nil
}
