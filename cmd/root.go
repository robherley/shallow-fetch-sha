package cmd

import (
	"errors"

	"github.com/robherley/shallow-fetch-sha/internal/fetch"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	options = &fetch.Options{}
	verbose bool

	rootCmd = &cobra.Command{
		Use:   "shallow-fetch-sha <repo> <sha> [flags]",
		Short: "Shallow fetch a specific git repository's commit to a directory.",
		Long: `For a given git repository and commit sha, fetch and checkout a specific commit
to save time and networking traffic. The resulting directory will not have any
ref/object history beyond the specified commit sha.

The repository can be specified as either SSH or HTTPS, but the commit must be
the 40 digit hexadecimal SHA1 representation.

Both SSH and Basic authentication are supported, granted the proper repository
URLs are specified. This program does not honor git-config files or options.

Note: this is only compatible with Git servers >= 2.50, they must support and
enable the 'uploadpack.allowReachableSHA1InWant' configuration option.
`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 2 {
				return errors.New("must specify both repo and sha arguments")
			}
			return nil
		},
		Run: func(c *cobra.Command, args []string) {
			if verbose {
				log.SetLevel(log.DebugLevel)
			}

			if err := options.BindArgs(args); err != nil {
				log.Fatalln(err)
			}

			if err := options.BindFlags(c.Flags()); err != nil {
				log.Fatalln(err)
			}

			if err := options.Validate(); err != nil {
				c.Usage()
				log.Fatalln(err)
			}

			if err := fetch.ShallowSHA(options); err != nil {
				log.Fatalln(err)
			}
		},
	}
)

func init() {
	rootCmd.Flags().StringP("directory", "d", ".", "working directory for the repository")
	rootCmd.Flags().StringP("username", "u", "", "username for basic authentication")
	rootCmd.Flags().StringP("password", "p", "", "password for basic authentication")
	rootCmd.Flags().StringP("key-path", "i", "", "pem encoded private key file for ssh authentication")
	rootCmd.Flags().StringP("key-passphrase", "P", "", "private key passphrase for ssh authentication")
	rootCmd.Flags().Bool("rm-dotgit", false, "remove the '.git' directory after pulling files")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.Flags().SortFlags = false
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
