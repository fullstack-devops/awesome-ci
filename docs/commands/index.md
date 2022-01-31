---
layout: default
title: Commands
nav_order: 2
has_children: true
has_toc: false
permalink: /docs/commands
---

# Commands

Overview of all available transfer parameters.

These commands are used to get information in a pipeline or local build.

Some calls can also create releases or more.

## Command structure

```bash
awesome-ci <subcommand> [subcommand-option]
```

| Option          | Description                                             | requiered |
| --------------- | ------------------------------------------------------- |:---------:|
| `-version`      | Returns current used version form awesome-ci            | false     |

### Subcommands

You can find out more about the subcommands by clicking on the relevant one in the navigation.

| Subcommand                              | Description                               |
| --------------------------------------- | ------------------------------------------|
| [release](/commands/createRelease.html) | creates a release at GitHub or            |
| [pr](/commands/getBuildInfos.html)      | prints out any git information and can manipulate |
| [parseJSON](/commands/parseJSON.html)   | can parse simple JSON                     |
| [parseYAML](/commands/parseYAML.html)   | can parse simple YAML files               |


## Release Body updates

Now available: Any string or any markdown file can now be attached in the release section. In addition, the release assets are attached as text to each release body with the associated sha256. More about that in the picture below and in the release section.

![Release Body with Asstes](../pictures/release-assets-readme.png "Release Body with Asstes")