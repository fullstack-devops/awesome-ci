[![Go Report Card](https://goreportcard.com/badge/github.com/layer87-labs/relctl)](https://goreportcard.com/report/github.com/layer87-labs/relctl)
[![GitHub release](https://img.shields.io/github/release/layer87-labs/relctl.svg)](https://github.com/layer87-labs/relctl/releases/latest)
[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/layer87-labs/relctl.svg)](https://github.com/layer87-labs/relctl)
[![License: Apache-2.0](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/layer87-labs/relctl/blob/main/LICENSE)

[![Publish Release](https://github.com/layer87-labs/relctl/actions/workflows/Release.yaml/badge.svg)](https://github.com/layer87-labs/relctl/actions/workflows/Release.yaml)
[![gh-pages](https://github.com/layer87-labs/relctl/actions/workflows/pages/pages-build-deployment/badge.svg)](https://github.com/layer87-labs/relctl/actions/workflows/pages/pages-build-deployment)

# relctl

**Description**: relctl is the smart connection between your pipeline for continuous integration and GitHub. The focus is on the release process, followed by the version management of [SemVer](https://semver.org/). The required version number is created with the correct naming of the branch prefix.

- **Technology stack**: This tool is written in golang
- **Status**: Stable.
- **Requests and Issues**: Please feel free to open an question or feature request in the Issue Board.
- **Supported environments**:
  - GitHub & GitHub Enterprise
  - GitHub actions
  - Jenkins Pipelines
- **Sweet Spot**: If you use GitHub or GitHub Enterprise and GitHub Actions, you can use relctl to its full potential!

## Getting Started

You can use this tool in your CI pipeline or locally on your command line. Just [download](https://github.com/layer87-labs/relctl/releases/latest/download/relctl) the most recently released version and get started.

## Usage

To integrate relctl into your pipeline, follow these steps:

1. Utilize the github action [layer87-labs/relctl-action](https://github.com/layer87-labs/relctl-action) to install relctl.
2. Configure your pipeline to use relctl.
3. Use the any command to interact with relctl.

You can find more information on how to integrate relctl into your pipeline in the [manual](https://layer87-labs.github.io/relctl/).

## Examples

You can find several examples of how to use relctl in the [examples section](https://layer87-labs.github.io/relctl/docs/examples) of the documentation.

## Frequently Asked Questions

You can find frequently asked questions in the [Questions and Answers](https://layer87-labs.github.io/relctl/docs/questions_and_answers) section of the documentation.

## Getting Help

If you have questions, concerns, or bug reports, please file an issue in this repository's Issue Tracker.

## Community

- [Contributing](https://github.com/layer87-labs/.github/blob/main/CONTRIBUTING.md)
- [Code of Conduct](https://github.com/layer87-labs/.github/blob/main/CODE_OF_CONDUCT.md)
- [Security Policy](https://github.com/layer87-labs/.github/blob/main/SECURITY.md)

## License

relctl is licensed under the Apache License, Version 2.0. You can find the license file [here](LICENSE).

## Credits

- [SemVer](https://semver.org/)
- [Cobra CLI](https://github.com/spf13/cobra)
