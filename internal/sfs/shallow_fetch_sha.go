package sfs

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	gitcfg "github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	log "github.com/sirupsen/logrus"
)

const (
	remoteName = git.DefaultRemoteName
	depth      = 1
)

func ShallowFetchSHA(fs *FileSystem, opts *Options) error {
	if opts == nil {
		return errors.New("must initialize options")
	}

	if fs == nil {
		return errors.New("must initialize filesystem")
	}

	log.WithFields(log.Fields{
		"sha": opts.SHA,
		"dir": fs.worktree.Root(),
	}).Info("shallow fetching repository")

	log.Debugln("initalizing repository on filesystem")
	repo, err := git.Init(fs.storage, fs.worktree)
	if err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"remote": remoteName,
		"url":    opts.Repo,
	}).Debugln("creating remote")
	_, err = repo.CreateRemote(&gitcfg.RemoteConfig{
		Name: remoteName,
		URLs: []string{opts.Repo},
	})
	if err != nil {
		return err
	}

	refspec := gitcfg.RefSpec(fmt.Sprintf(gitcfg.DefaultFetchRefSpec, opts.SHA))

	log.WithFields(log.Fields{
		"remote":  remoteName,
		"url":     opts.Repo,
		"refspec": refspec,
	}).Debugln("fetching ref")
	err = repo.Fetch(&git.FetchOptions{
		RemoteName: remoteName,
		Depth:      depth,
		RefSpecs: []gitcfg.RefSpec{
			refspec,
		},
		Progress: os.Stdout,
	})
	if err != nil {
		return err
	}

	log.Debugln("retrieving worktree")
	worktree, err := repo.Worktree()
	if err != nil {
		return err
	}

	if worktree == nil {
		return errors.New("unknown working tree")
	}

	log.WithFields(log.Fields{
		"hash": opts.SHA,
	}).Debugln("checking out hash")
	err = worktree.Checkout(&git.CheckoutOptions{
		Hash: plumbing.NewHash(opts.SHA),
	})
	if err != nil {
		return err
	}

	if opts.RemoveDotGit {
		log.Debugln("removing \".git\" directory")
		dotGitPath := filepath.Join(fs.worktree.Root(), git.GitDirName)
		if err := os.RemoveAll(dotGitPath); err != nil {
			return fmt.Errorf("unable to remove .git path: %s", err)
		}
	}

	return nil
}
