package sfs_test

import (
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/robherley/shallow-fetch-sha/internal/cli"
	"github.com/robherley/shallow-fetch-sha/internal/sfs"

	flag "github.com/spf13/pflag"
)

var _ = Describe("Options", func() {
	var (
		options       sfs.Options
		sshAuthOpts   sfs.SSHAuthOptions
		basicAuthOpts sfs.BasicAuthOptions
		dummyDir      = makeTemp()
	)

	BeforeEach(func() {
		options = sfs.Options{
			Repo:         publicRepo.SSH,
			SHA:          publicRepo.Commit,
			Directory:    dummyDir,
			RemoveDotGit: false,
			BasicAuth:    nil,
			SSHAuth:      nil,
		}

		sshAuthOpts = sfs.SSHAuthOptions{
			PEMPath:    sshKeyWithPassPath,
			Passphrase: sshPassphrase,
		}

		basicAuthOpts = sfs.BasicAuthOptions{
			Username: "token",
			Password: "notpassword",
		}
	})

	Describe("Validate", func() {
		It("should succeed with no auth", func() {
			Expect(options.Validate()).To(BeNil())
		})

		It("should succeed with ssh auth", func() {
			options.SSHAuth = &sshAuthOpts
			Expect(options.Validate()).To(BeNil())
		})

		It("should succeed with basic auth", func() {
			options.BasicAuth = &basicAuthOpts
			Expect(options.Validate()).To(BeNil())
		})

		It("should fail with both auth options", func() {
			options.BasicAuth = &basicAuthOpts
			options.SSHAuth = &sshAuthOpts
			Expect(options.Validate()).To(Not(BeNil()))
		})

		It("should fail for invalid repo", func() {
			options.Repo = ""
			Expect(options.Validate()).To(Not(BeNil()))
		})

		It("should fail for invalid sha", func() {
			options.SHA = "deadbeef"
			Expect(options.Validate()).To(Not(BeNil()))

			options.SHA = "zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz"
			Expect(options.Validate()).To(Not(BeNil()))
		})

		It("should fail for invalid basic auth", func() {
			options.BasicAuth = &sfs.BasicAuthOptions{
				Username: "",
				Password: basicAuthOpts.Password,
			}
			Expect(options.Validate()).To(Not(BeNil()))

			options.BasicAuth = &sfs.BasicAuthOptions{
				Username: basicAuthOpts.Username,
				Password: "",
			}
			Expect(options.Validate()).To(Not(BeNil()))
		})

		It("should fail for invalid ssh auth", func() {
			options.SSHAuth = &sfs.SSHAuthOptions{
				PEMPath:    "",
				Passphrase: sshAuthOpts.Passphrase,
			}
			Expect(options.Validate()).To(Not(BeNil()))
		})
	})

	Describe("BindArgs", func() {
		It("should fail for != 2 args", func() {
			badArgs := []string{"foo", "bar", "baz"}
			Expect(options.BindArgs(badArgs)).To(Not(BeNil()))

			badArgs = []string{"foo"}
			Expect(options.BindArgs(badArgs)).To(Not(BeNil()))

			badArgs = []string{}
			Expect(options.BindArgs(badArgs)).To(Not(BeNil()))
		})

		It("should succeed for == 2 args", func() {
			goodArgs := []string{publicRepo.SSH, publicRepo.Commit}
			Expect(options.BindArgs(goodArgs)).To(BeNil())

			Expect(options.Repo).To(Equal(goodArgs[0]))
			Expect(options.SHA).To(Equal(goodArgs[1]))
		})
	})

	Describe("BindFlags", func() {
		var (
			dummyFlags *flag.FlagSet
		)

		BeforeEach(func() {
			dummyFlags = flag.NewFlagSet("dummyflags", flag.ContinueOnError)
			cli.AddFlags(dummyFlags)
		})

		It("should bind directory flag", func() {
			directory := "./foo/bar"
			_ = dummyFlags.Set("directory", directory)

			Expect(options.BindFlags(dummyFlags)).To(BeNil())
			Expect(options.Directory).To(Equal(directory))
		})

		It("should bind username flag", func() {
			username := "bob"
			_ = dummyFlags.Set("username", username)

			Expect(options.BindFlags(dummyFlags)).To(BeNil())
			Expect(options.BasicAuth.Username).To(Equal(username))
		})

		It("should bind password flag", func() {
			password := "notpassword"
			_ = dummyFlags.Set("password", password)

			Expect(options.BindFlags(dummyFlags)).To(BeNil())
			Expect(options.BasicAuth.Password).To(Equal(password))
		})

		It("should bind key-path flag", func() {
			keypath := "/my/key.pem"
			_ = dummyFlags.Set("key-path", keypath)

			Expect(options.BindFlags(dummyFlags)).To(BeNil())
			Expect(options.SSHAuth.PEMPath).To(Equal(keypath))
		})

		It("should bind key-passphrase flag", func() {
			passphrase := "foo-bar-baz"
			_ = dummyFlags.Set("key-passphrase", passphrase)

			Expect(options.BindFlags(dummyFlags)).To(BeNil())
			Expect(options.SSHAuth.Passphrase).To(Equal(passphrase))
		})

		It("should bind rm-dotgit flag", func() {
			_ = dummyFlags.Set("rm-dotgit", "true")

			Expect(options.BindFlags(dummyFlags)).To(BeNil())
			Expect(options.RemoveDotGit).To(Equal(true))

			_ = dummyFlags.Set("rm-dotgit", "false")

			Expect(options.BindFlags(dummyFlags)).To(BeNil())
			Expect(options.RemoveDotGit).To(Equal(false))
		})
	})

	Describe("Auth", func() {
		It("should return nil with no auth", func() {
			auth, err := options.Auth()

			Expect(auth).To(BeNil())
			Expect(err).To(BeNil())
		})

		It("should return ssh auth method", func() {
			options.SSHAuth = &sshAuthOpts
			auth, err := options.Auth()

			Expect(auth).To(Not(BeNil()))
			Expect(err).To(BeNil())
			Expect(auth.Name()).To(Equal(ssh.PublicKeysName))

			sshAuth, ok := auth.(*ssh.PublicKeys)
			Expect(ok).To(BeTrue())

			Expect(sshAuth.User).To(Equal("git"))
		})

		It("should return ssh auth method w/ custom user", func() {
			options.SSHAuth = &sshAuthOpts
			options.Repo = "not" + publicRepo.SSH
			auth, err := options.Auth()

			Expect(auth).To(Not(BeNil()))
			Expect(err).To(BeNil())
			Expect(auth.Name()).To(Equal("ssh-public-keys"))

			sshAuth, ok := auth.(*ssh.PublicKeys)
			Expect(ok).To(BeTrue())

			Expect(sshAuth.User).To(Equal("notgit"))
		})

		It("should return basic auth method", func() {
			options.BasicAuth = &basicAuthOpts
			options.Repo = publicRepo.HTTPS
			auth, err := options.Auth()

			Expect(auth).To(Not(BeNil()))
			Expect(err).To(BeNil())
			Expect(auth.Name()).To(Equal("http-basic-auth"))

			basicAuth, ok := auth.(*http.BasicAuth)
			Expect(ok).To(BeTrue())

			Expect(basicAuth.Username).To(Equal(basicAuthOpts.Username))
			Expect(basicAuth.Password).To(Equal(basicAuthOpts.Password))
		})
	})
})
