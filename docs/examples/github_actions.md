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


### Build and publish npm packages
```yaml
name: Publish

on:
  push:
    branches:    
    - 'master'

jobs:
  setup:
    runs-on: nodejs
    steps:
      - name: Checkout code
        uses: actions/checkout@v1
      - uses: eksrvb/awesome-ci@main
      - uses: actions/setup-node@v2
        with:
          node-version: '14.x'
          registry-url: 'https://registry.npmjs.org'
      - name: install npm packages
        run: npm install
      - name: package Applikation
        run: npm run build
      name: DryRun Release
        run: ./awesome-ci createRelease -publishNpm dist/my-project/
        env:
          GITHUB_TOKEN: $\{\{ secrets.GITHUB_TOKEN \}\}
```

### Build and upload a go project
```yaml
name: Publish

on:
  push:
    branches:    
    - 'master'

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout code
        uses: actions/checkout@v2
      - uses: eksrvb/awesome-ci@main
      - name: Set up Go
        uses: actions/setup-go@v2
        with:
          go-version: 1.15
      - name: set environment variables
        run: awesome-ci getBuildInfos
        env:
          GITHUB_TOKEN: $\{\{ secrets.GITHUB_TOKEN \}\}
      - name: Build and write version to binary
        run: go build -v -ldflags "-X main.version=$NEXT_VERSION"
      - name: Create Release and upload
        run: ./awesome-ci createRelease -uploadArtifacts file=my-compiled-binary
        env:
          GITHUB_TOKEN: $\{\{ secrets.GITHUB_TOKEN \}\}
```