---
title: GitHub Actions
---

# GitHub Actions examples

### Build a Release

This is an example from th relctl project you can find the original workflow [here](https://github.com/layer87-labs/relctl/blob/main/.github/workflows/Release.yaml).

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
      release-id: ${{ steps.tag.outputs.RELCTL_RELEASE_ID }}
      version: ${{ steps.tag.outputs.RELCTL_NEXT_VERSION }}
    steps:
      - name: Check out the repo
        uses: actions/checkout@v4
      - name: Setup relctl
        uses: layer87-labs/relctl-action@main

      - name: create release
        id: tag
        run: relctl release create --merge-sha ${{ github.sha }}
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
        uses: actions/checkout@v4
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
          key: relctl-${{ github.sha }}-${{ hashFiles('out/relctl*') }}
          restore-keys: |
            relctl-${{ github.sha }}

  publish_release:
    runs-on: ubuntu-latest
    needs: [create_release, build]
    steps:
      - name: Checkout code
        uses: actions/checkout@v4
      - name: Setup relctl
        uses: layer87-labs/relctl-action@main

      - name: get cached build outputs
        uses: actions/cache@v3
        env:
          cache-name: cache-outputs-modules
        with:
          path: out/
          key: relctl-${{ github.sha }}

      - name: Publish Release
        run: relctl release publish --release-id "$RELCTL_RELEASE_ID" --asset "file=out/$ARTIFACT1" --asset "file=out/$ARTIFACT2"
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
          RELCTL_RELEASE_ID: ${{ needs.create_release.outputs.release-id }}
          ARTIFACT1: relctl_${{ needs.create_release.outputs.version }}_amd64
          ARTIFACT2: relctl_${{ needs.create_release.outputs.version }}_arm64
```

You need more examples? Please open an issue!
