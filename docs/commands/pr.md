---
layout: default
title: pr
parent: Commands
nav_order: 2
---

- [Overview](#overview)
- [Subcommands](#subcommands)
  - [-number](#-number)
  - [-format](#-format)
- [Examples](#examples)
- [Special to github actions](#special-to-github-actions)

## Overview

```bash
awesome-ci pr <subcommand> [subcommand-option]
```

## Subcommands

| Subcommand          | Description                                                                 |
| ------------------- | --------------------------------------------------------------------------- |
| `info`              | creates an release, but doesn't publish it                                  |


| Subcommand option | Description                                                                   |
| ----------------- | ----------------------------------------------------------------------------- |
| `-number`         | overwrite the issue number                                                    |
| `-format`         | pastes the required output to the console. This can be extracted to variables |

### -number

By default, awesome-ci recognizes the number of the PullRequest being built. To speed up this process and make it more stable, the `-number` can optionally be specified if known. This brings additional stability to the workflow.

### -format

The `-format` option can put out your needed information about your current git status.

  Hint: use a seperatoa as you like, the below values would be replaced!

Possible infos are: `patchLevel`, `pr`, `version`, `nextVersion`

#### Examples

```bash
#### Info output:
Pull Request: 17
Current release version: 1.0.0
Patch level: feature
Possible new release version: 1.1.0
```

```bash
awesome-ci getBuildInfos -format "pr,next_version"
# Output:
17,1.1.0
```

### Special to github actions

With a github action, all available information is always set as environment variables. Once set the Variables ale awailable in all steps and runners.

```bash
#### Setting Env variables:
ACI_PR=17
ACI_ORGA=eksrvb
ACI_REPO=playground
ACI_BRANCH=bugfix/test-pr
ACI_PATCH_LEVEL=bugfix
ACI_VERSION=0.4.4
ACI_NEXT_VERSION=0.4.5

#### Info output:
Pull Request: 17
Current release version: 0.4.4
Patch level: bugfix
Possible new release version: 0.4.5
```
