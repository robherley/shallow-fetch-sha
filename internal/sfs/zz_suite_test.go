package sfs_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/bradleyfalzon/ghinstallation/v2"
	log "github.com/sirupsen/logrus"

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
		Commit: "1bd1c0c32ff7d4b4db95a3591a5c018b86708c8b",
		ExpectedFiles: []string{
			"README.md",
			"index.js",
			"public.txt",
		},
	}

	privateRepo = repoFixture{
		SSH:    "git@github.com:robherley/fixture-private-repo.git",
		HTTPS:  "https://github.com/robherley/fixture-private-repo.git",
		Commit: "0196c49057e45387838aa3a1d0601e8c7a4317d0",
		ExpectedFiles: []string{
			"README.md",
			"index.js",
			"private.txt",
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
		"SSH_PASSPHRASE",
		"SSH_PEM_NO_PASS",
		"SSH_PEM_WITH_PASS",
		"BOT_PEM",
		"BOT_INSTALLATION_ID",
		"BOT_ID",
	}
	for _, envVar := range requiredEnvVars {
		if os.Getenv(envVar) == "" {
			panic("missing env var: " + envVar)
		}
	}

	sshPassphrase = os.Getenv("SSH_PASSPHRASE")
	sshKeyNoPassPath = makePEMFile("id_no_pass.pem", os.Getenv("SSH_PEM_NO_PASS"))
	sshKeyWithPassPath = makePEMFile("id_with_pass.pem", os.Getenv("SSH_PEM_WITH_PASS"))

	// authenticate as a bot, get an access token
	botKeyPath := makePEMFile("bot.pem", os.Getenv("BOT_PEM"))
	botID, err := strconv.Atoi(os.Getenv("BOT_ID"))
	plsno(err)
	installationID, err := strconv.Atoi(os.Getenv("BOT_INSTALLATION_ID"))
	plsno(err)
	itr, err := ghinstallation.NewKeyFromFile(http.DefaultTransport, int64(botID), int64(installationID), botKeyPath)
	plsno(err)

	botToken, err = itr.Token(context.Background())
	plsno(err)

	// reduce noise
	log.SetOutput(ioutil.Discard)
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

func plsno(err error) {
	if err != nil {
		panic(err)
	}
}

func makeTemp() string {
	tmpdir, err := os.MkdirTemp("", "sfs-*")
	plsno(err)
	dirsToCleanUp = append(dirsToCleanUp, tmpdir)
	return tmpdir
}

func makePEMFile(filename, contents string) string {
	tmpdir := makeTemp()
	fp := filepath.Join(tmpdir, filename)
	bs := []byte(contents)

	plsno(os.WriteFile(fp, bs, 0600))
	fmt.Println("added key to", fp)

	return fp
}
