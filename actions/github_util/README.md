# Github Util

Check if the PR meets the specifications

* The PR must have either the pr/feature or pr/bugfix label

* For pr/bugfix PRs, they must be associated with a bug issue in the format: `fix https://github.com/vesoft-inc/nebula-ng/issues/8001`

```yaml
name: check pr

on:
  workflow_dispatch:
  pull_request:
    types: [synchronize, reopened, labeled, edited]
    branches:
      - master

jobs:
  check:
    name: check pr
    steps:
      - name: check
        uses: vesoft-inc/.github/actions/github_util@master
        with:
          token: ${{ secrets.GH_PAT }}
```
