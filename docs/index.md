---
layout: default
title: Home
nav_order: 1
---

## Welcome to the awesome-ci documentation!

### Notice:

Every command that you can use is in the sidebar under commands. All options are listed there.

If you need an example for your pipeline you can find it in the sidebar under the tab examples.

### Supported naming rules and effects on the version

The patching of the version only takes effect if the merged branch begins with the following aliases, for example: `feature/my-awesome-feature`

> The tailing `/` behind the alias is **always** requiered!

| SemVer | supported aliases                           | version example |
| ------ | ------------------------------------------- | --------------- |
| MAJOR  | major                                       | 1.0.0 => 2.0.0  |
| MINOR  | minor, feature                              | 1.0.0 => 1.1.0  |
| PATCH  | patch, bugfix, fix                          | 1.0.0 => 1.0.1  |



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
| `GITHUB_TOKEN`            | Must provided in pipeling as `env.GITHUB_TOKEN` (see examples)  | true      |
| `GIT_DEFAULT_BRANCH_NAME` | overrides the default branch name (default: `main`)             | false     |

> To see your Jenkins environment variables go to: `${YOUR_JENKINS_HOST}/env-vars.html`