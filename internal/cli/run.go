package cli

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"

	sfs "github.com/robherley/shallow-fetch-sha/internal/sfs"
)

var (
	opts    = &sfs.Options{}
	flags   = pflag.NewFlagSet("shallow-fetch-sha", pflag.ContinueOnError)
	verbose bool
	help    bool
)

const (
	description = `For a given git repository and commit sha, fetch and checkout a specific commit
to save time and networking traffic. The resulting directory will not have any
ref/object history beyond the specified commit sha.

The repository can be specified as either SSH or HTTPS, but the commit must be
the 40 digit hexadecimal SHA1 representation. Both SSH and Basic authentication
are supported, granted the proper repository URLs are specified. This program
does not honor git-config files or options.

Note: this is only compatible with Git servers >= 2.50, they must support and
enable the 'uploadpack.allowReachableSHA1InWant' configuration option.`
	usage = "shallow-fetch-sha <repo> <sha> [flags]"
)

func init() {
	flags.StringP("directory", "d", ".", "working directory for the repository")
	flags.StringP("username", "u", "", "username for basic authentication")
	flags.StringP("password", "p", "", "password for basic authentication")
	flags.StringP("key-path", "i", "", "pem encoded private key file for ssh authentication")
	flags.StringP("key-passphrase", "P", "", "private key passphrase for ssh authentication")
	flags.BoolP("rm-dotgit", "D", false, "remove the '.git' directory after pulling files")
	flags.BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	flags.BoolVarP(&help, "help", "h", false, "help for shallow-fetch-sha")

	flags.SortFlags = false
	flags.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s\nsee \"shallow-fetch-sha --help\" for more information\n", usage)
	}

	if err := flags.Parse(os.Args[1:]); err != nil {
		failWithUsage(err)
	}
}

func helpme() {
	fmt.Fprintln(os.Stderr, description)
	fmt.Fprintf(os.Stderr, "\nUsage:\n  %s\n", usage)
	fmt.Fprintf(os.Stderr, "\nFlags:\n%s", flags.FlagUsages())
	os.Exit(0)
}

func failWithUsage(err error) {
	log.Errorln(err)
	flags.Usage()
	os.Exit(1)
}

func Run() {
	if verbose {
		log.SetLevel(log.DebugLevel)
	}

	if help {
		helpme()
	}

	if err := opts.BindArgs(flags.Args()); err != nil {
		failWithUsage(err)
	}

	if err := opts.BindFlags(flags); err != nil {
		failWithUsage(err)
	}

	if err := opts.Validate(); err != nil {
		failWithUsage(err)
	}

	fs, err := sfs.NewFileSystem(opts.Directory, sfs.DiskMode)
	if err != nil {
		failWithUsage(err)
	}

	if err := sfs.ShallowFetchSHA(fs, opts); err != nil {
		log.Fatalln(err)
	}
}
