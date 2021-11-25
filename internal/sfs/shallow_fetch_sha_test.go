package sfs_test

import (
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/robherley/shallow-fetch-sha/internal/sfs"
)

func checkFiles(dir string, expectedFiles []string) bool {
	foundFiles, err := os.ReadDir(dir)
	if err != nil {
		return false
	}

	for _, expected := range expectedFiles {
		seen := false
		for _, found := range foundFiles {
			if expected == found.Name() {
				seen = true
				break
			}
		}
		if !seen {
			return false
		}
	}

	return true
}

var _ = Describe("ShallowFetchSha", func() {
	It("should fetch a public repo via https", func() {
		tmpDir := makeTemp()
		options := sfs.Options{
			Repo:      publicRepo.HTTPS,
			SHA:       publicRepo.Commit,
			Directory: tmpDir,
			Silent:    true,
		}

		err := sfs.ShallowFetchSHA(&options)
		Expect(err).To(BeNil())

		seenAllFiles := checkFiles(tmpDir, publicRepo.ExpectedFiles)
		Expect(seenAllFiles).To(BeTrue())
	})

	It("should fetch a public repo via ssh", func() {
		tmpDir := makeTemp()
		options := sfs.Options{
			Repo:      publicRepo.SSH,
			SHA:       publicRepo.Commit,
			Directory: tmpDir,
			Silent:    true,
		}

		err := sfs.ShallowFetchSHA(&options)
		Expect(err).To(BeNil())

		seenAllFiles := checkFiles(tmpDir, publicRepo.ExpectedFiles)
		Expect(seenAllFiles).To(BeTrue())
	})

	It("should fetch a private repo via ssh", func() {
		tmpDir := makeTemp()
		options := sfs.Options{
			Repo:      privateRepo.SSH,
			SHA:       privateRepo.Commit,
			Directory: tmpDir,
			SSHAuth: &sfs.SSHAuthOptions{
				PEMPath: sshKeyNoPassPath,
			},
			Silent: true,
		}

		err := sfs.ShallowFetchSHA(&options)
		Expect(err).To(BeNil())

		seenAllFiles := checkFiles(tmpDir, privateRepo.ExpectedFiles)
		Expect(seenAllFiles).To(BeTrue())
	})

	It("should fetch a private repo via ssh with passphrase", func() {
		tmpDir := makeTemp()
		options := sfs.Options{
			Repo:      privateRepo.SSH,
			SHA:       privateRepo.Commit,
			Directory: tmpDir,
			SSHAuth: &sfs.SSHAuthOptions{
				PEMPath:    sshKeyWithPassPath,
				Passphrase: sshPassphrase,
			},
			Silent: true,
		}

		err := sfs.ShallowFetchSHA(&options)
		Expect(err).To(BeNil())

		seenAllFiles := checkFiles(tmpDir, privateRepo.ExpectedFiles)
		Expect(seenAllFiles).To(BeTrue())
	})

	It("should fetch a private repo via https", func() {
		tmpDir := makeTemp()
		options := sfs.Options{
			Repo:      privateRepo.HTTPS,
			SHA:       privateRepo.Commit,
			Directory: tmpDir,
			BasicAuth: &sfs.BasicAuthOptions{
				Username: "x-access-token",
				Password: botToken,
			},
			Silent: true,
		}

		err := sfs.ShallowFetchSHA(&options)
		Expect(err).To(BeNil())

		seenAllFiles := checkFiles(tmpDir, privateRepo.ExpectedFiles)
		Expect(seenAllFiles).To(BeTrue())
	})

	It("should remove dot git if specified", func() {
		tmpDir := makeTemp()
		options := sfs.Options{
			Repo:         publicRepo.SSH,
			SHA:          publicRepo.Commit,
			Directory:    tmpDir,
			RemoveDotGit: true,
			Silent:       true,
		}

		err := sfs.ShallowFetchSHA(&options)
		Expect(err).To(BeNil())

		_, err = os.Stat(filepath.Join(tmpDir, git.GitDirName))
		Expect(os.IsNotExist(err)).To(BeTrue())
	})
})
