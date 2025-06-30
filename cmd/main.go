package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/Xebec19/reimagined-lamp/internal/backend"
	"github.com/Xebec19/reimagined-lamp/internal/fs"
	"github.com/Xebec19/reimagined-lamp/internal/mountlib"
)

func main() {

	var (
		remote  = flag.String("remote", "", "Remote name to mount")
		mount   = flag.String("mount", "", "Mount point")
		verbose = flag.Bool("verbose", false, "Enable verbose output")
	)

	flag.Parse()

	if *remote == "" || *mount == "" {
		fmt.Fprintf(os.Stderr, "Usage: %s -remote <remote> -mount <mount>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Example: %s -remote local:/home/user/data -mount /mnt/mydata\n", os.Args[0])
		os.Exit(1)
	}

	parts := strings.SplitN(*remote, ":", 2)

	if len(parts) != 2 {
		fmt.Fprintf(os.Stderr, "Invalid remote format. Use <type>:<path>\n")
		os.Exit(1)
	}

	remoteType, remotePath := parts[0], parts[1]

	var r fs.Remote
	switch remoteType {
	case "local":
		r = backend.NewLocalRemote(remotePath)
	default:
		fmt.Fprintf(os.Stderr, "Unsupported remote type: %s\n", remoteType)
	}

	if err := mountlib.Mount(*mount, r, *verbose); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to mount: %v\n", err)
		os.Exit(1)
	}
}
