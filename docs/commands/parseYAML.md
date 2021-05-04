---
layout: default
title: parseYAML
parent: Commands
nav_order: 4
---

# Parsing a Yaml File

```bash
awesome-ci parseYAML [subcommand-option]
```

| Subcommand option | Description                    |
| ----------------- | ------------------------------ |
| `-file`           | your file location             |
| `-value`          | your value you want to extract |

## Examlpe

Example demo.yaml:

```yaml
value1: hello
value2: world
deepObject:
  value1: hello
  value2: world
```

Example command 1:

```shell
awesome-ci parseYAML -file demo.yaml -value .value2
world
```

Example command 2:

```shell
awesome-ci parseYAML -file demo.yaml -value .deepObject
map[value1:hello value2:world]
```

You can check more details of the yaml construct in one of the following versions.
