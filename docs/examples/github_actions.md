---
layout: default
title: GitHub Actions
parent: Examples
nav_order: 1
---

# GitHub Actions examples
{: .no_toc }

<details open markdown="block">
  <summary>
    Table of contents
  </summary>
  {: .text-delta }
1. TOC
{:toc}
</details>


### Build a Release

This is an example from th awesome-ci project you can find the original workflow [here](https://github.com/fullstack-devops/awesome-ci/blob/main/.github/workflows/Release.yaml)

```yaml
name: Publish Release

on:
  push:
    branches:
      - "main"
    paths-ignore:
      - "README.md"
      - 'docs/**'
      - '.github/ISSUE_TEMPLATE/**'
      - '.github/PULL_REQUEST_TEMPLATE.md'


jobs:
  generate_infos:
    runs-on: ubuntu-latest
    outputs:
      releaseid: $\{\{ steps.tag.outputs.releaseid \}\}
      version: $\{\{ steps.tag.outputs.version \}\}
      pr: $\{\{ steps.tag.outputs.pr \}\}
    steps:
      - name: Check out the repo
        uses: actions/checkout@v2
      - name: Setup awesome-ci
        uses: fullstack-devops/awesome-ci-action@main

      - name: collect infos and create release
        run: |
          awesome-ci pr info
          awesome-ci release create # making a draft
        env:
          GITHUB_TOKEN: $\{\{ secrets.GITHUB_TOKEN \}\}

      - name: collect infos
        id: tag
        shell: bash
        run: |
          echo "::set-output name=version::$ACI_VERSION"
          echo "::set-output name=pr::$ACI_PR"
          echo "::set-output name=releaseid::$ACI_RELEASE_ID"

  build:
    runs-on: ubuntu-latest
    needs: generate_infos
    strategy:
      matrix:
        arch: ["amd64", "arm64"]
    steps:
      - name: Checkout code
        uses: actions/checkout@v2
      - name: Set up Go
        uses: actions/setup-go@v2
        with:
          go-version: 1.17

      - name: Build "$\{\{ matrix.arch \}\}"
        run: go build -v -ldflags "-X main.version=$\{\{ needs.generate_infos.outputs.version \}\}" -o out/awesome-ci_$\{\{ needs.generate_infos.outputs.version \}\}_$\{\{ matrix.arch \}\}
        env:
          GOOS: linux
          GOARCH: "$\{\{ matrix.arch \}\}"

      - name: Cache build outputs
        uses: actions/cache@v2
        env:
          cache-name: cache-outputs-modules
        with:
          path: out/
          key: awesome-ci-$\{\{ github.sha \}\}-$\{\{ hashFiles('out/awesome-ci*') \}\}
          restore-keys: |
            awesome-ci-$\{\{ github.sha \}\}

  publish_release:
    runs-on: ubuntu-latest
    needs: [generate_infos, build]
    steps:
      - name: Checkout code
        uses: actions/checkout@v2
      - name: Setup awesome-ci
        uses: fullstack-devops/awesome-ci-action@main

      - name: get cached build outputs
        uses: actions/cache@v2
        env:
          cache-name: cache-outputs-modules
        with:
          path: out/
          key: awesome-ci-$\{\{ github.sha \}\}
      
      - name: Publish Release
        run: awesome-ci release publish -releaseid "$ACI_RELEASE_ID" -upload "file=out/$ARTIFACT1,file=out/$ARTIFACT2"
        env:
          GITHUB_TOKEN: $\{\{ secrets.GITHUB_TOKEN \}\}
          ACI_RELEASE_ID: $\{\{ needs.generate_infos.outputs.releaseid \}\}
          ARTIFACT1: awesome-ci_$\{\{ needs.generate_infos.outputs.version \}\}_amd64
          ARTIFACT2: awesome-ci_$\{\{ needs.generate_infos.outputs.version \}\}_arm64
```

You need more examples? Please open an issue!