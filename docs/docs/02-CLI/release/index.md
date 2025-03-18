# release

Manage GitHub releases with ease

## Synopsis

The release command is used to manage GitHub releases. It provides subcommands to create and publish releases, allowing you to automate the release process and integrate it into CI/CD workflows. Use this command to streamline the release tagging and deployment of your software projects.

```
awesome-ci release [flags]
```

## Options

```
  -b, --body string             custom release message (markdown string or file)
      --dry-run                 make dry-run before writing version to Git by calling it
  -h, --help                    help for release
      --hotfix                  create a hotfix release
      --merge-sha string        set the merge sha
  -l, --patch-level string      predefine patch level of version to Update
      --prnumber int            overwrite the issue number
      --release-branch string   set release branch (default: git default)
      --release-prefix string   set a custom release prefix (default -> Release or Hotfix)
      --version string          override version to Update
```

## Options inherited from parent commands

```
  -v, --verbose   verbose output
```

## SEE ALSO

* **awesome-ci**	 - Awesome CI make your release tagging easy
* **awesome-ci release create**	 - Create a new GitHub release
* **awesome-ci release publish**	 - Publish a recently created GitHub release

##### Auto generated on 18-Mar-2025
