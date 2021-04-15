---
layout: default
title: Commands
nav_order: 2
---

# Commands

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

| Subcommand          | Description                                             |
| ------------------- | ------------------------------------------------------- |
| `createRelease`     | creates a release at GitHub or GitLab                   |
| `createRelease`     | prints out any git information and can manipulate these |
