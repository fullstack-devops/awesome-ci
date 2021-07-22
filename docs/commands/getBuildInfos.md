---
layout: default
title: getBuildInfos
parent: Commands
nav_order: 2
---

# Creating a release

```bash
awesome-ci createRelease [subcommand-option]
```

| Subcommand option | Description                                                                   |
| ----------------- | ----------------------------------------------------------------------------- |
| `-version`        | overrides any version from git and patches the given string.                  |
| `-patchLevel`     | overrides the patchLevel. make shure our following the alias definition.      |
| `-format`         | pastes the required output to the console. This can be extracted to variables |

### -version

The `-version` option can overwrite the evaluated version.
This can be useful in connection with `-patchLevel` when creating a manual release.

### -patchLevel

The `-patchLevel` option can overwrite the evaluated patchLevel.
This can be useful in connection with `-version` when creating a manual release.

### -format

The `-format` option can put out your needed information about your current git status.

Hint: use a seperatoa as you like, the below values would be replaced!

Possible infos are: `patchLevel`, `pr`, `version`, `nextVersion`

#### Examples

```bash
awesome-ci getBuildInfos -patchLevel feature -version 1.0.0

#### Info output:
Pull Request: 17
Current release version: 1.0.0
Patch level: feature
Possible new release version: 1.1.0
```

```bash
awesome-ci getBuildInfos -patchLevel feature -version 1.0.0 -format "pr,next_version"
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
