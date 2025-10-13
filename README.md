# Git Backups

[![License](https://img.shields.io/badge/License-MIT-yellow.svg)](https://github.com/AntonKosov/git-backups/blob/master/LICENSE.md)
[![Tests](https://github.com/AntonKosov/git-backups/actions/workflows/quality-of-code.yaml/badge.svg)](https://github.com/AntonKosov/git-backups/actions/workflows/quality-of-code.yaml)
[![Coverage Status](https://coveralls.io/repos/github/AntonKosov/git-backups/badge.svg?branch=master)](https://coveralls.io/github/AntonKosov/git-backups?branch=master)

A Docker-based tool for creating local backups of Git repositories with support for multiple platforms and authentication methods.

## Features

* **Multi-platform support**: Generic Git repositories and GitHub-specific profiles
* **Flexible authentication**: SSH Agent and SSH keys
* **Batch operations**: Backup multiple repositories with a single configuration
* **Docker deployment**: Easy setup and consistent runtime environment

## Supported Profiles

### Generic Profile

For any Git repository accessible via HTTPS or SSH.

### GitHub Profile  

Specialized support for GitHub repositories with features like:
* Automatic repository discovery based on user affiliation
* Personal access token authentication
* Repository filtering (include/exclude lists)

## Prerequisites

* For GitHub profiles, you'll need a personal access token with `repo` scope. The token will be used for reading repository lists from your account.
* For private and GitHub repositories, SSH keys or SSH agent forwarding is required.

## Quick Start

1. Create a configuration file (`config.yaml`)
1. Make sure the default `~/.ssh/known_hosts` file has all needed hosts added. Missing hosts can be added with `ssh-keyscan github.com >> ~/.ssh/known_hosts` command.
1. If the SSH forwarding is used, verify that the agent is up and running (`ssh-add -l`). If it's not running, add a private SSH key permanently or temporarily to the current session (`ssh-add ~/.ssh/<private key>`).
1. Run the Docker container with mounted volumes (listed below).

### Configuration Example

```yaml
version: 1

profiles:
  # Generic repositories - supports multiple profiles
  generic:
    - profile: "GitLab"
      # Backup destination directory
      root_folder: "/app/backup/gitlab"
      # Optional: Private SSH key for authentication
      # private_ssh_key: "/app/ssh_key"
      # Repository list with custom folder names
      targets:
        - url: "git@gitlab.com:Username1/repo_name_1.git"
          folder: "repo_name_1"
        - url: "git@gitlab.com:Username2/repo_name_2.git"
          folder: "repo_name_2"

  # GitHub repositories - supports multiple profiles  
  github:
    - profile: "GitHub Personal"
      root_folder: "/app/backup/github"
      # Repository ownership filter (at least one required)
      affiliation: "owner,collaborator,organization_member"
      # GitHub personal access token with "repo" scope
      token: "ghp_XXX"
      # Optional: Private SSH key for git operations
      # private_ssh_key: "/app/ssh_key"
      # Optional: Only backup specific repositories
      # include: ["repo_name_1", "repo_name_2"]
      # Optional: Exclude specific repositories (overrides include)
      # exclude: ["repo_name_3"]
```

### Docker Usage

```shell
docker run --rm \
    --mount type=bind,src=~/.ssh/known_hosts,dst=/home/appuser/.ssh/known_hosts,readonly \
    --mount type=bind,src=./config.yaml,dst=/app/config.yaml,readonly \
    --mount type=bind,src=./backup,dst=/app/backup \
    --volume type=bind,src=$SSH_AUTH_SOCK,dst=/tmp/ssh-auth.sock,readonly \
    --env SSH_AUTH_SOCK=/tmp/ssh-auth.sock \
    ghcr.io/antonkosov/git-backups:latest
```

## Volume Mounts

| Container Path | Mode | Description |
|----------------|------|-------------|
| `/home/appuser/.ssh/known_hosts` | `ro` | Known hosts file | 
| `/app/config.yaml` | `ro` | Configuration file |
| Any path | `ro` | Private SSH key (optional, path configurable) |
| Any path | `rw` | Backup destination directory (path configurable) |

**Notes**:
* SSH keys and backup paths can be customized in the configuration file.
* Multiple SSH keys can be mounted to different container paths.
* SSH keys should have the owner which matches the container user (`chown 100:101 <ssh key>`).

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details.
