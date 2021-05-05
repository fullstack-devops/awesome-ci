---
layout: default
title: Commands
nav_order: 2
has_children: true
permalink: /docs/commands
---

# Commands
{: .no_toc }

Overview of all available transfer parameters.

These commands are used to get information in a pipeline or local build.

Some calls can also create releases or more.

## Command structure

```bash
awesome-ci [option] [subcommand] [subcommand-option]
```

| Option          | Description                                             | requiered |
| --------------- | ------------------------------------------------------- |:---------:|
| `-version`      | Returns current used version form awesome-ci            | false     |

### Subcommands

You can find out more about the subcommands by clicking on the relevant one in the navigation.

| Subcommand                                                                         | Description                                             |
| ---------------------------------------------------------------------------------- | ------------------------------------------------------- |
| [createRelease](https://eksrvb.github.io/awesome-ci/commands/createRelease.html)   | creates a release at GitHub or GitLab                   |
| [getBuildInfos](https://eksrvb.github.io/awesome-ci/commands/getBuildInfos.html)   | prints out any git information and can manipulate these |
| [parseJSON](https://eksrvb.github.io/awesome-ci/commands/parseJSON.html)           | can parse simple JSON files                             |
| [parseYAML](https://eksrvb.github.io/awesome-ci/commands/parseYAML.html)           | can parse simple YAML files                             |
