# publish

Publish a recently created GitHub release

## Synopsis

Publish a recently created GitHub release

Publish will publish a previously created release to GitHub. The release body
can be written in markdown and will be rendered as such on GitHub. Additionally,
any number of assets can be uploaded, including but not limited to zip and tgz
files. The assets will be uploaded to GitHub and linked to the release.

The assets can be specified as a list of local file paths.

The zip and tgz mechanisms can be used by specifying the path to a directory
containing the files to be uploaded. The directory will be zipped or tarred and
gzipped and uploaded to GitHub as a single asset. The name of the asset will be
the name of the directory with the appropriate extension appended. For example,
if the directory is named "myfiles", the asset will be named "myfiles.zip" or "myfiles.tgz".

Example of using the release publish command

To publish a release with a specific release ID and assets, use the following command:

    awesome-ci release publish --release-id 12345 \
        --asset "file=out/awesome-ci_v1.0.0_amd64" \
        --asset "file=out/awesome-ci_v1.0.0_arm64"
	awesome-ci release publish --release-id 12345 --asset "file=out/awesome-ci_v1.0.0_amd64" --asset "file=out/awesome-ci_v1.0.0_arm64"

If the release ID is not provided, the command will look for the 'ACI_RELEASE_ID' environment variable:

    export ACI_RELEASE_ID=12345
    awesome-ci release publish \
        --asset "file=out/awesome-ci_v1.0.0_amd64" \
        --asset "file=out/awesome-ci_v1.0.0_arm64"
	export ACI_RELEASE_ID=12345
	awesome-ci release publish --asset "file=out/awesome-ci_v1.0.0_amd64" --asset "file=out/awesome-ci_v1.0.0_arm64"

You can also publish a release with a directory as a zip asset:

    awesome-ci release publish --release-id 12345 \
        --asset "zip=out/myfiles"

The assets should be specified as local file paths to be uploaded to the GitHub release.

```
awesome-ci release publish [flags]
```

## Options

```
  -a, --asset stringArray   add an asset to the release, can be specified multiple times.
  -h, --help                help for publish
      --release-id int      publish an early defined release (also looking for env ACI_RELEASE_ID)
```

## Options inherited from parent commands

```
  -b, --body string             custom release message (markdown string or file)
      --dry-run                 make dry-run before writing version to Git by calling it
      --hotfix                  create a hotfix release
      --merge-sha string        set the merge sha
  -l, --patch-level string      predefine patch level of version to Update
      --prnumber int            overwrite the issue number
      --release-branch string   set release branch (default: git default)
      --release-prefix string   set a custom release prefix (default -> Release or Hotfix)
  -v, --verbose                 verbose output
      --version string          override version to Update
```

## SEE ALSO

* **awesome-ci release**	 - Manage GitHub releases with ease

##### Auto generated on 18-Mar-2025
