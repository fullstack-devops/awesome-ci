---
layout: default
title: createRelease
parent: Commands
nav_order: 1
---


- [Overview](#overview)
- Options
  - [-version](#-version)

## Overview

```bash
awesome-ci createRelease [... subcommand-option]
```

| Subcommand option   | Description                                                                 |
| ------------------- | --------------------------------------------------------------------------- |
| `-version`          | overrides any version from git and patches the given string.                |
| `-patchLevel`       | overrides the patchLevel. make shure our following the alias definition.    |
| `-publishNpm`       | after creating a release publish the given sources to a npm registry. This also overrides the npm version an set these in your package.json |
| `-uploadArtifacts`  | uploads the given Artifacts to a release. Eg.: "out/awesome-ci,..."         |
| `-dry-run`          | doesn't create a release. Prints out what it would do and check permissions |


### -version

The `-version` option can overwrite the evaluated version.
This can be useful in connection with `-patchLevel` when creating a manual release.

```bash
awesome-ci createRelease -version 0.1.0
```

### -patchLevel

The `-patchLevel` option can overwrite the evaluated patchLevel.
This can be useful in connection with `-version` when creating a manual release.

```bash
awesome-ci createRelease -patchLevel feature
```

### -uploadArtifacts

The `-patchLevel` option can updload a single or multiple artifacts.

However, you must choose the format of the artefacts.

eg.: `-uploadArtifacts "file=path/to/file,file=path/to/second/file"`


... more documentation to be done ;)