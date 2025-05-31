# go-git-scan

This small application can scan all GitHub repositories of a given username, and fetch the names and emails of all collaborators who made commits.

In the future we plan to support:
- GitHub (self-hosted)
- GitLab.com
- GitLab (self-hosted)
- Codeberg
- Sourcehut

# Running go-git-scan

You will need to install Go.

Then, you can simply run:

```
$ cd src
$ go run main.go <platform> <username>
```

For example:

```
$ go run main.go github L1Cafe
```