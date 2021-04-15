# Awesome CI

**Description**: This tool makes your workflow easier! With the help of [SemVer](https://semver.org/lang/) and naming conventions, a lot of time can be saved when creating a release.

- **Technology stack**: This tool is written in golang
- **Status**: Alpha, at release the [CHANGELOG](CHANGELOG.md) would be updated.
- **Requests**: Please feel free to open an question or feature request in the Issue Board.

## Table of Contents

- [Usage](#usage)
  - [Supported naming rules and effects on the version](#supported-naming-rules-and-effects-on-the-version)
  - [Examples for your CI Pipeline](#examples-for-your-ci-pipeline)
  - [Supported commands](#supported-commands)
  - [Requiered and optional environment variables](#requiered-and-optional-environment-variables)
- [Known issues](#known-issues)
- [Getting help](#getting-help)
- [Open source licensing info](#open-source-licensing-info)
- [Credits and references](#credits-and-references)

## Usage

to use this tool in your ci pipeline, [download](https://github.com/eksrvb/awesome-ci/releases/latest/download/awesome-ci) the most recently published release. How you can integrate this into your respective pipeline can be found in the following document.

> Hint: this tool automatically detects your environment. Supported are __Jenkins__, __GitHub Actions__ and __GitLab CI__

### Supported naming rules and effects on the version

The patching of the version only takes effect if the merged branch begins with the following aliases, for example: `feature/my-awesome-feature`

> The tailing `/` behind the alias is **always** requiered!

| SemVer | supported aliases                           | version example |
| ------ | ------------------------------------------- | --------------- |
| MAJOR  | major                                       | 1.0.0 => 2.0.0  |
| MINOR  | minor, feature                              | 1.0.0 => 1.1.0  |
| PATCH  | patch, bugfix, fix                          | 1.0.0 => 1.0.1  |

### Examples for your CI Pipeline

- [GitHub Actions](examples/GitHub_Actions.md)
- Jenkins Pipeline (coming soon)
- GitLab CI (coming soon)

### Supported commands

To Print all available by calling `awesome-ci -help` 

### Requiered and optional environment variables

List of all environmental variables used per CI tool.

**GitHub Actions**
| Environment variable      | Description                                                     | requiered |
| ------------------------- | --------------------------------------------------------------- |:---------:|
| `GITHUB_API_URL`          | Returns the API URL. (Already set in runner)                    | true      |
| `GITHUB_REPOSITORY`       | The owner and repository name. (Already set in runner)          | true      |
| `GITHUB_TOKEN`            | Must provided in workflow as `env:` (see examples)              | true      |
| `GIT_DEFAULT_BRANCH_NAME` | overrides the default branch name (default: `main`)             | false     |

**Jenkins Pipeline**
| Environment variable      | Description                                                     | requiered |
| ------------------------- | --------------------------------------------------------------- |:---------:|
| `JENKINS_URL`             | Returns the URL of your Jenkins instance. (Already set)         | true      |
| `GIT_URL`                 | Will only be set by using the GitHub Plugin.                    | true      |
| `GITHUB_TOKEN`            | Must provided in pipeling as `env.GITHUB_TOKEN` (see examples)  | true      |
| `GIT_DEFAULT_BRANCH_NAME` | overrides the default branch name (default: `main`)             | false     |
 > To see your Jenkins environment variables go to: `${YOUR_JENKINS_HOST}/env-vars.html`

## Known issues

There are no known Issues yet.

## Getting help

If you have questions, concerns, bug reports, etc, please file an issue in this repository's Issue Tracker.

## Getting involved

General instructions on _how_ to contribute: [CONTRIBUTING](CONTRIBUTING.md)

General instructions _during_ the contribution period: [CONTRIBUTION](CONTRIBUTION.md)


----

## Open source licensing info
1. [LICENSE](LICENSE)


----

## Credits and references

- [SemVer](https://semver.org/)
