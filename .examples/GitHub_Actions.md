# Using awesome-ci in GitHub Actions

## Table of Contents
- [Examples](#examples)
  - [Build and publish npm packages](#build-and-publish-npm-packages)
  - [Build and upload a go project](#build-and-publish-npm-packages)

## Examples


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
      - uses: actions/setup-node@v2
        with:
          node-version: '14.x'
          registry-url: 'https://registry.npmjs.org'
      - name: Set up awesome-ci
        run: |
          wget https://github.com/eksrvb/awesome-ci/releases/latest/download/awesome-ci
          chmod +x awesome-ci
      - name: install npm packages
        run: npm install
      - name: package Applikation
        run: npm run build
      name: DryRun Release
        run: ./awesome-ci createRelease -publishNpm dist/my-project/
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
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
      - name: Set up Go
        uses: actions/setup-go@v2
        with:
          go-version: 1.15
      - name: Set up awesome-ci
        run: |
          wget https://github.com/eksrvb/awesome-ci/releases/latest/download/awesome-ci
          chmod +x awesome-ci
      - name: get next Version for project build
        id: version_step
        run: echo "::set-output name=new_version::$(./awesome-ci getBuildInfos -format version)"
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
      - name: Build and write version to binary
        run: go build -v -ldflags "-X main.version=${{ steps.version_step.outputs.new_version }}"
      - name: Create Release and upload
        run: ./awesome-ci createRelease -uploadArtifacts my-compiled-bynary
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```