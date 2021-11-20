package sfs

import (
	"fmt"
	"path/filepath"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/cache"
	"github.com/go-git/go-git/v5/storage"
	"github.com/go-git/go-git/v5/storage/filesystem"
	"github.com/go-git/go-git/v5/storage/memory"
	log "github.com/sirupsen/logrus"
)

type FileSystem struct {
	storage  storage.Storer
	worktree billy.Filesystem
}

type FSMode string

const (
	DiskMode   FSMode = "disk"
	MemoryMode FSMode = "mem"
)

func NewFileSystem(dir string, mode FSMode) (*FileSystem, error) {
	log.WithFields(log.Fields{
		"mode": mode,
		"dir":  dir,
	}).Debugln("filesystem: initalizing working tree and storage")

	fs := &FileSystem{}

	switch mode {
	case DiskMode:
		absDir, err := filepath.Abs(dir)
		if err != nil {
			return nil, fmt.Errorf("invalid directory: %s", err)
		}

		log.WithFields(log.Fields{
			"absdir": absDir,
		}).Debugln("filesystem(disk): resolved absolute directory")

		wt := osfs.New(absDir)

		dotGit, err := wt.Chroot(git.GitDirName)
		if err != nil {
			return nil, err
		}

		fs.worktree = wt
		fs.storage = filesystem.NewStorage(dotGit, cache.NewObjectLRUDefault())

		return fs, nil
	case MemoryMode:
		fs.worktree = memfs.New()
		fs.storage = memory.NewStorage()

		return fs, nil
	default:
		return nil, fmt.Errorf("%q is an invalid storage mode", mode)
	}
}

func (fs *FileSystem) GetWorkTree() billy.Filesystem {
	return fs.worktree
}

func (fs *FileSystem) GetStorage() storage.Storer {
	return fs.storage
}
