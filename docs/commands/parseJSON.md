---
layout: default
title: parseJSON
parent: Commands
nav_order: 3
---

# Parsing a JSON File

```bash
awesome-ci parseJSON [subcommand-option]
```

| Subcommand option | Description                    |
| ----------------- | ------------------------------ |
| `-file`           | your file location             |
| `-value`          | your value you want to extract |

## Examlpe

Example demo.json:

```json
{
  "value1": "hello",
  "value2": "world",
  "deepObject": {
    "value1": "hello",
    "value2": "world"
  }
}
```

Example command 1:

```shell
awesome-ci parseJSON -file demo.json -value .value2
world
```

Example command 2:

```shell
awesome-ci parseJSON -file demo.json -value .deepObject
map[value1:hello value2:world]
```

You can check more details of the json construct in one of the following versions.
