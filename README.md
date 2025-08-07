# Git Backups

[![License](https://img.shields.io/badge/License-MIT-yellow.svg)](https://github.com/AntonKosov/git-backups/blob/master/LICENSE.md)
[![Tests](https://github.com/AntonKosov/git-backups/actions/workflows/quality-of-code.yaml/badge.svg)](https://github.com/AntonKosov/git-backups/actions/workflows/quality-of-code.yaml)
[![Coverage Status](https://coveralls.io/repos/github/AntonKosov/git-backups/badge.svg?branch=master)](https://coveralls.io/github/AntonKosov/git-backups?branch=master)


🚧 Under Development...

`config.yaml`

```yaml
version: 1

repositories:
  # It is not implemented yet
  generic:
    - category: "category name"
      root_folder: "/home/user/git_backup/folder_name"
      targets:
        - url: "https://github.com/Username1/repo_name_1.git"
          folder: "repo_folder_name_1"
        - url: "https://github.com/Username2/repo_name_2.git"
          folder: "repo_folder_name_2"
  # It is not implemented yet
  github:
    - category: "category name 2"
      root_folder: "/home/user/git_backup/folder_name_2"
      token: "GH_XXX"
      # 'include' is optional, only repositories listed in this field will be included
      include: ["repo_name_1", "repo_name_2"]
      # 'exclude' is optional, repositories which are listed here will NOT be included (even if they are listed in the 'include' field)
      exclude: ["repo_name_3"]
```