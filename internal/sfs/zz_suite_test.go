package sfs_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type repoFixture struct {
	URL           string
	Commit        string
	ExpectedFiles []string
}

var (
	publicRepo = repoFixture{
		URL:    "https://github.com/robherley/fixture-public-repo",
		Commit: "6c7ddd59f0feec261a1788ced53c3e06f8cddda6",
		ExpectedFiles: []string{
			"README.md",
			"index.js",
		},
	}

	privateRepo = repoFixture{
		URL:    "https://github.com/robherley/fixture-private-repo",
		Commit: "ce7f2ec17e070ddcb5c4a9a19e32267518ebab7e",
		ExpectedFiles: []string{
			"README.md",
			"index.js",
		},
	}
)

func TestSfs(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Sfs Suite")
}
