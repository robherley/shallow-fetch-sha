package sfs_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/joho/godotenv"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type repoFixture struct {
	SSH           string
	HTTPS         string
	Commit        string
	ExpectedFiles []string
}

const (
	testEnvFile = "../../.env"
)

var (
	publicRepo = repoFixture{
		SSH:    "git@github.com:robherley/fixture-public-repo.git",
		HTTPS:  "https://github.com/robherley/fixture-public-repo.git",
		Commit: "6c7ddd59f0feec261a1788ced53c3e06f8cddda6",
		ExpectedFiles: []string{
			"README.md",
			"index.js",
		},
	}

	privateRepo = repoFixture{
		SSH:    "git@github.com:robherley/fixture-private-repo.git",
		HTTPS:  "https://github.com/robherley/fixture-private-repo.git",
		Commit: "ce7f2ec17e070ddcb5c4a9a19e32267518ebab7e",
		ExpectedFiles: []string{
			"README.md",
			"index.js",
		},
	}

	dirsToCleanUp = make([]string, 0)

	// populated from env vars for testing w/ private fixtures
	// keys are scoped as deploy keys to the fixture repos
	// token is for a bot only scoped to the private repo
	sshKeyNoPassPath   string
	sshKeyWithPassPath string
	sshPassphrase      string
	botToken           string
)

func TestShallowFetchSHA(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ShallowFetchSHA Suite")
}

var _ = BeforeSuite(func() {
	// when testing locally we'll have a .env file
	// in ci, these will already be env vars
	_ = godotenv.Load(testEnvFile)

	requiredEnvVars := []string{
		"BOT_TOKEN",
		"SSH_PASSPHRASE",
		"SSH_PEM_NO_PASS",
		"SSH_PEM_WITH_PASS",
	}
	for _, envVar := range requiredEnvVars {
		if os.Getenv(envVar) == "" {
			panic("missing env var: " + envVar)
		}
	}

	botToken = os.Getenv("BOT_TOKEN")
	sshPassphrase = os.Getenv("SSH_PASSPHRASE")
	sshKeyNoPassPath = makePEMFile("id_no_pass.pem", os.Getenv("SSH_PEM_NO_PASS"))
	sshKeyWithPassPath = makePEMFile("id_no_pass.pem", os.Getenv("SSH_PEM_NO_PASS"))

	fmt.Println(botToken)
	fmt.Println(sshPassphrase)
	fmt.Println(sshKeyNoPassPath)
	fmt.Println(sshKeyWithPassPath)
})

var _ = AfterSuite(func() {
	// clean up best effort
	for _, dir := range dirsToCleanUp {
		// make sure we don't do anything stupid
		if strings.HasPrefix(dir, os.TempDir()) {
			_ = os.RemoveAll(dir)
		}
	}
})

func makeTemp() string {
	tmpdir, err := os.MkdirTemp("", "sfs-*")
	if err != nil {
		panic(err)
	}
	dirsToCleanUp = append(dirsToCleanUp, tmpdir)
	return tmpdir
}

func makePEMFile(filename, contents string) string {
	tmpdir := makeTemp()
	fp := filepath.Join(tmpdir, filename)
	bs := []byte(contents)

	if err := os.WriteFile(fp, bs, 0600); err != nil {
		panic(err)
	}

	return fp
}
