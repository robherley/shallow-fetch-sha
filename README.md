# shallow-fetch-sha (sfs) 🏗️

For a given git repository and commit, fetch and checkout *just* that commit without any history. This can be extremely useful in CI/CD systems that need to ship code from repositories with large ref/object history and only need files at a specific commit.

Effectively the same as the following on a given directory:

```console
you@local:~$ git init
you@local:~$ git remote add origin $REPO
you@local:~$ git fetch --depth 1 origin $SHA
you@local:~$ git checkout $SHA
```
<sup>Credit to [sschuberth](https://stackoverflow.com/a/43136160) from StackOverflow</sup>

This utility is shipped as a standalone binary as well as a container. It is built using [go-git](https://github.com/go-git/go-git), a pure Go implementation of git.

**Note:** this is only compatible with Git servers >= 2.50, since they must support (and enable) `uploadpack.allowReachableSHA1InWant`.

## Usage

### CLI (standalone binary)

```console
you@local:~$ shallow-fetch-sha --help
For a given git repository and commit sha, fetch and checkout a specific commit
to save time and networking traffic. The resulting directory will not have any
ref/object history beyond the specified commit sha.

The repository can be specified as either SSH or HTTPS, but the commit must be
the 40 digit hexadecimal SHA1 representation. Both SSH and Basic authentication
are supported, granted the proper repository URLs are specified. This program
does not honor git-config files or options.

Note: this is only compatible with Git servers >= 2.50, they must support and
enable the 'uploadpack.allowReachableSHA1InWant' configuration option.

Usage:
  shallow-fetch-sha <repo> <sha> [flags]

Flags:
  -d, --directory string        working directory for the repository (default ".")
  -u, --username string         username for basic authentication
  -p, --password string         password for basic authentication
  -i, --key-path string         pem encoded private key file for ssh authentication
  -P, --key-passphrase string   private key passphrase for ssh authentication
  -D, --rm-dotgit               remove the '.git' directory after pulling files
  -s, --silent                  silent output (takes precedence over verbose)
  -v, --verbose                 verbose output
  -h, --help                    help for shallow-fetch-sha
```

### Container

The entrypoint is the `shallow-fetch-sha` binary, and the default working directory is `/usr/src/repo`. The user is the default non-priviledged `guest (uid=405)` user within the [alpine](https://hub.docker.com/_/alpine/) image.

Basic usage:

```console
you@local:~$ podman run -it registry.tbd/robherley/shallow-fetch-sha:$TAG <repo> <sha> [flags]
```

Fetching a repo/commit and saving to a local directory:

```console
you@local:~$ podman run -it --rm -v $(pwd)/repo:/usr/src/repo registry.tbd/robherley/shallow-fetch-sha:$TAG https://github.com/robherley/reb.gg fd40042f1a21da61b4abddebbe94f21dc700ffb0
INFO[0000] shallow fetching repository                   dir=/usr/src/repo sha=fd40042f1a21da61b4abddebbe94f21dc700ffb0
Enumerating objects: 4, done.
Counting objects: 100% (4/4), done.
Compressing objects: 100% (3/3), done.
Total 4 (delta 0), reused 3 (delta 0), pack-reused 0
```
