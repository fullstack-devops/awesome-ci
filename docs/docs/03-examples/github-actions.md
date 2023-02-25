---
title: GitHub Actions
---

# GitHub Actions examples

### Build a Release

This is an example from th awesome-ci project you can find the original workflow [here](https://github.com/fullstack-devops/awesome-ci/blob/main/.github/workflows/Release.yaml).

```yaml title="release.yaml"
name: Publish Release

on:
  push:
    branches:
      - "main"

jobs:
  create_release:
    runs-on: ubuntu-latest
    outputs:
      release-id: ${{ steps.tag.outputs.ACI_RELEASE_ID }}
      version: ${{ steps.tag.outputs.ACI_NEXT_VERSION }}
    steps:
      - name: Check out the repo
        uses: actions/checkout@v3
      - name: Setup awesome-ci
        uses: fullstack-devops/awesome-ci-action@main

      - name: create release
        id: tag
        run: awesome-ci release create --merge-sha ${{ github.sha }}
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}

  build:
    runs-on: ubuntu-latest
    needs: create_release
    strategy:
      matrix:
        arch: ["amd64", "arm64"]
    steps:
      - name: Checkout code
        uses: actions/checkout@v3
      - name: Set up Go
        uses: actions/setup-go@v3
        with:
          go-version: 1.19

      - name: Build "${{ matrix.arch }}"
        run: make VERSION="${{ needs.create_release.outputs.version }}"
        env:
          GOOS: linux
          GOARCH: "${{ matrix.arch }}"

      - name: Cache build outputs
        uses: actions/cache@v3
        env:
          cache-name: cache-outputs-modules
        with:
          path: out/
          key: awesome-ci-${{ github.sha }}-${{ hashFiles('out/awesome-ci*') }}
          restore-keys: |
            awesome-ci-${{ github.sha }}

  publish_release:
    runs-on: ubuntu-latest
    needs: [create_release, build]
    steps:
      - name: Checkout code
        uses: actions/checkout@v3
      - name: Setup awesome-ci
        uses: fullstack-devops/awesome-ci-action@main

      - name: get cached build outputs
        uses: actions/cache@v3
        env:
          cache-name: cache-outputs-modules
        with:
          path: out/
          key: awesome-ci-${{ github.sha }}

      - name: Publish Release
        run: awesome-ci release publish --release-id "$ACI_RELEASE_ID" --asset "file=out/$ARTIFACT1" --asset "file=out/$ARTIFACT2"
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
          ACI_RELEASE_ID: ${{ needs.create_release.outputs.release-id }}
          ARTIFACT1: awesome-ci_${{ needs.create_release.outputs.version }}_amd64
          ARTIFACT2: awesome-ci_${{ needs.create_release.outputs.version }}_arm64
```

You need more examples? Please open an issue!
