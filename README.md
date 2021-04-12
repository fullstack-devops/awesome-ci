# Awesome CI

**Description**: This tool makes your workflow easier! With the help of [SemVer](https://semver.org/lang/) and naming conventions, a lot of time can be saved when creating a release.

- **Technology stack**: This tool is written in golang
- **Status**: Alpha, at release the [CHANGELOG](CHANGELOG.md) would be updated.
- **Requests**: Please feel free to open an question or feature request in the Issue Board.

## Usage

[Download](https://github.com/eksrvb/awesome-ci/releases/latest/download/awesome-ci) the latest [Release](https://github.com/eksrvb/awesome-ci/releases) or include the Docker container in your Multi stage Build.

### In your CI Pipeline

- [Jenkins Pipeline](.examples/GitLab_CI.md)
- [GitHub Actions](.examples/GitHub_Actions.md)
- [GitLab CI](.examples/GitLab_CI.md)

### Supported commands

To Print all available by calling `awesome-ci -help` 

### Requiered and optional environment variables

List of all environmental variables used per CI tool.

**GitHub Actions**
| Environment variable      | Description                                                     | requiered |
| ------------------------- | --------------------------------------------------------------- |:---------:|
| `GITHUB_API_URL`          | Returns the API URL.                                            | true      |
| `GITHUB_REPOSITORY`       | The owner and repository name.                                  | true      |
| `GITHUB_TOKEN`            | Must provided in workflow as `env:` (see examples)              | true      |
| `GIT_DEFAULT_BRANCH_NAME` | overrides the default branch name (default: `main`)             | false     |

## Installation

The installation is needed only for contirbuting.

Go to your GOPATH under `src/` and run: `go get https://gitlab.com/eksrvb/awesome-ci-semver`

## How to test the software

If the software includes automated tests, detail how to run those tests.

## Known issues

There are no known Issues yet.

## Getting help

If you have questions, concerns, bug reports, etc, please file an issue in this repository's Issue Tracker.

## Getting involved

This section should detail why people should get involved and describe key areas you are
currently focusing on; e.g., trying to get feedback on features, fixing certain bugs, building
important pieces, etc.

General instructions on _how_ to contribute should be stated with a link to [CONTRIBUTING](CONTRIBUTING.md).


----

## Open source licensing info
1. [TERMS](TERMS.md)
2. [LICENSE](LICENSE)
3. [CFPB Source Code Policy](https://github.com/cfpb/source-code-policy/)


----

## Credits and references

- [SemVer](https://semver.org/lang/de/)
