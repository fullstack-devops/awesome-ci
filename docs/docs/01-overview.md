---
title: Overview
---

# Welcome to the Awesome CI!

This project is the smart connection between your pipeline for continuous integration and your version management like GitHub. The focus is on the release process, followed by the version management of SemVer. The required version number is created with the correct naming of the branch prefix.

You can use this tool in your CI pipeline or locally on your command line. Just download the most recently released version and get started. You can find out how to integrate this into your respective pipeline in the following document. There are also several examples in the examples section of the documentation. If an example is not included, please feel free to inquire about a related issue.

If more functionality is needed you can just open a problem in this project and of course bugs can be fixed in the same way by filing a bug report.

If you have any questions, you can find a form on the issue board. First, make sure your question is already in the Questions and Answers section before asking a question. You can find frequently asked questions directly in the "Questions and Answers" section.

:::info

Every command that you can use is in the sidebar at cli. All options are listed there or use the `awesome-ci help` command.

:::
:::tip

If you need an example for your pipeline you can find it in the sidebar under the tab examples.

:::

## Supported naming rules and effects on the version

The patching of the version only takes effect if the merged branch begins with the following aliases, for example: `feature/my-awesome-feature`

:::caution
The tailing `/` behind the alias is **always** requiered!
:::

| SemVer | supported aliases                      | version example |
| ------ | -------------------------------------- | --------------- |
| MAJOR  | `major`                                | 1.2.3 => 2.0.0  |
| MINOR  | `minor`, `feature`, `feat`             | 1.2.3 => 1.3.0  |
| PATCH  | `patch`, `fix`, `bugfix`, `dependabot` | 1.2.3 => 1.2.4  |

:::info
see also [override specialties](#override-specialties)
:::

![awesome-ci release process](/img/release-process.drawio.svg "awesome-ci release process")
![awesome-ci workflow](/img/aci-workflow.drawio.png "awesome-ci workflow")

:::tip
Awesoce CI automatically detects your environment. Supported are **Jenkins Pipelines** and **GitHub Actions**
:::

## Override specialties

To set some attributes during developement you can comment a pullrequest.

| command                       | description                                                   |
| ----------------------------- | ------------------------------------------------------------- |
| `aci_patch_level: major`      | create a major version bump                                   |
| `aci_version_override: 2.1.0` | set the version to 2.1.0 using only semver compatible syntax! |

## Requiered and optional environment variables

List of all environmental variables used per CES (code execution service).

### GitHub Actions

| Environment variable | Description                                        | Status        | Requiered |
| -------------------- | -------------------------------------------------- | ------------- | :-------: |
| `CI`                 | Is set by GitHub actions and returns `true`        | set by runner |   true    |
| `GITHUB_SERVER_URL`  | The GitHub-Server URL.                             | set by runner |   true    |
| `GITHUB_REPOSITORY`  | The owner and repository name.                     | set by runner |   true    |
| `GITHUB_TOKEN`       | Must provided in workflow as `env:` (see examples) | set by runner |   true    |

### Jenkins Pipeline

| Environment variable | Description                                                    | Status                       | requiered |
| -------------------- | -------------------------------------------------------------- | ---------------------------- | :-------: |
| `CI`                 | Is set by Jenkins Pipeline and returns `true`                  | set by jenkins               |   true    |
| `JENKINS_URL`        | Returns the URL of your Jenkins instance. (Already set)        | set by jenkins               |   true    |
| `GIT_URL`            | Will only be set by using the GitHub Plugin.                   | set by jenkins plugin github |   true    |
| `GITHUB_REPOSITORY`  | The owner and repository name.                                 | must be set manually         |   true    |
| `GITHUB_TOKEN`       | Must provided in pipeline as `env.GITHUB_TOKEN` (see examples) | must be set manually         |   true    |

:::tip
To see your Jenkins environment variables go to: `${YOUR_JENKINS_HOST}/env-vars.html`
:::
