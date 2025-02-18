[![Go Report Card](https://goreportcard.com/badge/github.com/fullstack-devops/awesome-ci)](https://goreportcard.com/report/github.com/fullstack-devops/awesome-ci)
[![GitHub release](https://img.shields.io/github/release/fullstack-devops/awesome-ci.svg)](https://github.com/fullstack-devops/awesome-ci/releases/latest)
[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/fullstack-devops/awesome-ci.svg)](https://github.com/fullstack-devops/awesome-ci)
[![License: Apache-2.0](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/fullstack-devops/awesome-ci/blob/main/LICENSE)

[![Publish Release](https://github.com/fullstack-devops/awesome-ci/actions/workflows/Release.yaml/badge.svg)](https://github.com/fullstack-devops/awesome-ci/actions/workflows/Release.yaml)
[![gh-pages](https://github.com/fullstack-devops/awesome-ci/actions/workflows/pages/pages-build-deployment/badge.svg)](https://github.com/fullstack-devops/awesome-ci/actions/workflows/pages/pages-build-deployment)

# Awesome CI

**Description**: Awesome CI is the smart connection between your pipeline for continuous integration and GitHub. The focus is on the release process, followed by the version management of [SemVer](https://semver.org/). The required version number is created with the correct naming of the branch prefix.

- **Technology stack**: This tool is written in golang
- **Status**: Stable.
- **Requests and Issues**: Please feel free to open an question or feature request in the Issue Board.
- **Supported environments**:
  - GitHub & GitHub Enterprise
  - GitHub actions
  - Jenkins Pipelines
- **Sweet Spot**: If you use GitHub or GitHub Enterprise and GitHub Actions, you can use awesome-ci to its full potential!

## Getting Started

You can use this tool in your CI pipeline or locally on your command line. Just [download](https://github.com/fullstack-devops/awesome-ci/releases/latest/download/awesome-ci) the most recently released version and get started.

## Usage

To integrate Awesome CI into your pipeline, follow these steps:

1. Utilize the github action [fullstack-devops/awesome-ci-action](https://github.com/fullstack-devops/awesome-ci-action) to install Awesome CI.
2. Configure your pipeline to use Awesome CI.
3. Use the any command to interact with Awesome CI.

You can find more information on how to integrate Awesome CI into your pipeline in the [manual](https://fullstack-devops.github.io/awesome-ci/).

## Examples

You can find several examples of how to use Awesome CI in the [examples section](https://fullstack-devops.github.io/awesome-ci/docs/examples) of the documentation.

## Frequently Asked Questions

You can find frequently asked questions in the [Questions and Answers](https://fullstack-devops.github.io/awesome-ci/docs/questions_and_answers) section of the documentation.

## Getting Help

If you have questions, concerns, or bug reports, please file an issue in this repository's Issue Tracker.

## Contributing

If you want to contribute to Awesome CI, please read the [CONTRIBUTING](docs/CONTRIBUTING.md) guide.

## License

Awesome CI is licensed under the Apache License, Version 2.0. You can find the license file [here](LICENSE).

## Credits

- [SemVer](https://semver.org/)
- [Cobra CLI](https://github.com/spf13/cobra)
