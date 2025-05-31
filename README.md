# go-git-scan

This small application can scan all GitHub public repositories of a given username, and fetch the names and emails of all collaborators who made commits.

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

# Disclaimer

This software is able to retrieve public data from Git servers. This software does not attempt to circumvent any security mechanisms like rate-limiting, Web Application Firewalls, and more. It does not search authors on any other records, and the software does not try to hide its identity to the Git servers it connects to. While the potential for harmful usage is recognised, the author does not authorise or condone any criminal usage derived from this software.

## License

This project is licensed under the GNU General Public License v3.0 - see the [LICENSE](LICENSE) file for details.
