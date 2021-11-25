package cli

import (
	"fmt"
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"

	sfs "github.com/robherley/shallow-fetch-sha/internal/sfs"
)

var (
	opts    = &sfs.Options{}
	flags   = pflag.NewFlagSet("shallow-fetch-sha", pflag.ContinueOnError)
	silent  bool
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
	usage = "sfs <repo> <sha> [flags]"
)

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

func AddFlags(flagset *pflag.FlagSet) {
	flagset.StringP("directory", "d", ".", "working directory for the repository")
	flagset.StringP("username", "u", "", "username for basic authentication")
	flagset.StringP("password", "p", "", "password for basic authentication")
	flagset.StringP("key-path", "i", "", "pem encoded private key file for ssh authentication")
	flagset.StringP("key-passphrase", "P", "", "private key passphrase for ssh authentication")
	flagset.BoolP("rm-dotgit", "D", false, "remove the '.git' directory after pulling files")
	flagset.BoolVarP(&silent, "silent", "s", false, "silent output (takes precedence over verbose)")
	flagset.BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	flagset.BoolVarP(&help, "help", "h", false, "help for shallow-fetch-sha")
}

func Run() {
	AddFlags(flags)

	flags.SortFlags = false
	flags.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s\nsee \"shallow-fetch-sha --help\" for more information\n", usage)
	}

	if err := flags.Parse(os.Args[1:]); err != nil {
		failWithUsage(err)
	}

	if silent {
		log.SetOutput(ioutil.Discard)
		opts.Silent = true
	}

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

	if err := sfs.ShallowFetchSHA(opts); err != nil {
		log.Fatalln(err)
	}
}
