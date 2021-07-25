---
layout: default
title: Home
nav_order: 1
---

## Welcome to the Awesome CI documentation!

This project is the smart connection between your pipeline for continuous integration and your version management like GitLab or GitHub. The focus is on the release process, followed by the version management of SemVer. The required version number is created with the correct naming of the branch prefix.

You can use this tool in your CI pipeline or locally on your command line. Just download the most recently released version and get started. You can find out how to integrate this into your respective pipeline in the following document. There are also several examples in the examples section of the documentation. If an example is not included, please feel free to inquire about a related issue.

If more functionality is needed you can just open a problem in this project and of course bugs can be fixed in the same way by filing a bug report.

If you have any questions, you can find a form on the issue board. First, make sure your question is already in the Questions and Answers section before asking a question. You can find frequently asked questions directly in the "Questions and Answers" section.

### Notice:

Every command that you can use is in the sidebar under commands. All options are listed there.

If you need an example for your pipeline you can find it in the sidebar under the tab examples.

### Supported naming rules and effects on the version

The patching of the version only takes effect if the merged branch begins with the following aliases, for example: `feature/my-awesome-feature`

> The tailing `/` behind the alias is **always** requiered!

| SemVer | supported aliases                           | version example |
| ------ | ------------------------------------------- | --------------- |
| MAJOR  | `comment PR`                                | 1.2.3 => 2.0.0  |
| MINOR  | feature                                     | 1.2.3 => 1.3.0  |
| PATCH  | bugfix                                      | 1.2.3 => 1.2.4  |

![awesome-ci workflow](aci-workflow.drawio.png "awesome-ci workflow")

> Hint: this tool automatically detects your environment. Supported are __Jenkins__, __GitHub Actions__ and __GitLab CI__


### Requiered and optional environment variables

List of all environmental variables used per CI tool.

#### GitHub Actions

| Environment variable      | Description                                                     | requiered |
| ------------------------- | --------------------------------------------------------------- |:---------:|
| `GITHUB_API_URL`          | Returns the API URL. (Already set in runner)                    | true      |
| `GITHUB_REPOSITORY`       | The owner and repository name. (Already set in runner)          | true      |
| `GITHUB_TOKEN`            | Must provided in workflow as `env:` (see examples)              | true      |
| `GIT_DEFAULT_BRANCH_NAME` | overrides the default branch name (default: `main`)             | false     |

#### Jenkins Pipeline

| Environment variable      | Description                                                     | requiered |
| ------------------------- | --------------------------------------------------------------- |:---------:|
| `JENKINS_URL`             | Returns the URL of your Jenkins instance. (Already set)         | true      |
| `GIT_URL`                 | Will only be set by using the GitHub Plugin.                    | true      |
| `GITHUB_TOKEN`            | Must provided in pipeline as `env.GITHUB_TOKEN` (see examples)  | true      |
| `GIT_DEFAULT_BRANCH_NAME` | overrides the default branch name (default: `main`)             | false     |

> To see your Jenkins environment variables go to: `${YOUR_JENKINS_HOST}/env-vars.html`