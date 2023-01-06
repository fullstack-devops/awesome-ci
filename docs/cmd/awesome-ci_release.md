## awesome-ci release

release

### Synopsis

tbd

```
awesome-ci release [flags]
```

### Options

```
  -b, --body string             custom release message (markdow string or file)
      --dry-run                 make dry-run before writing version to Git by calling it
  -h, --help                    help for release
      --hotfix                  create a hotfix release
      --merge-sha string        set the merge sha
  -L, --patch-level string      predefine patch level of version to Update
      --prnumber int            overwrite the issue number
      --release-branch string   set release branch (default: git default)
      --version string          override version to Update
```

### Options inherited from parent commands

```
  -v, --verbose   verbose output
```

### SEE ALSO

* [awesome-ci](awesome-ci.md)	 - Awesome CI make your release tagging easy
* [awesome-ci release create](awesome-ci_release_create.md)	 - create a GitHub release
* [awesome-ci release publish](awesome-ci_release_publish.md)	 - publish a GitHub release

###### Auto generated by spf13/cobra on 4-Jan-2023