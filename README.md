# Git Backups

[![License](https://img.shields.io/badge/License-MIT-yellow.svg)](https://github.com/AntonKosov/git-backups/blob/master/LICENSE.md)
[![Tests](https://github.com/AntonKosov/git-backups/actions/workflows/quality-of-code.yaml/badge.svg)](https://github.com/AntonKosov/git-backups/actions/workflows/quality-of-code.yaml)
[![Coverage Status](https://coveralls.io/repos/github/AntonKosov/git-backups/badge.svg?branch=master)](https://coveralls.io/github/AntonKosov/git-backups?branch=master)

A Docker-based tool for creating local backups of Git repositories with support for multiple platforms and authentication methods.

## Features

- **Multi-platform support**: Generic Git repositories and GitHub-specific profiles
- **Flexible authentication**: SSH keys and GitHub personal access tokens
- **Batch operations**: Backup multiple repositories with a single configuration
- **Docker deployment**: Easy setup and consistent runtime environment

## Supported Profiles

### Generic Profile
For any Git repository accessible via HTTPS or SSH.

### GitHub Profile  
Specialized support for GitHub repositories with features like:
- Automatic repository discovery based on user affiliation
- Personal access token authentication
- Repository filtering (include/exclude lists)

## Prerequisites

For GitHub profiles, you'll need a personal access token with `repo` scope. The token will be used for:
- Reading repository lists from your account
- Cloning/fetching private repositories (if SSH key not provided)

**Note**: SSH authentication is recommended over token-based authentication to avoid exposing tokens in remote URLs.

## Quick Start

1. Create a configuration file (`config.yaml`)
2. Run the Docker container with mounted volumes (listed below)

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
      private_ssh_key: "/app/ssh_key"
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
      private_ssh_key: "/app/ssh_key"
      # Optional: Only backup specific repositories
      # include: ["repo_name_1", "repo_name_2"]
      # Optional: Exclude specific repositories (overrides include)
      # exclude: ["repo_name_3"]
```

### Docker Usage

```shell
docker run --rm \
        -v config.yaml:/app/config.yaml:ro \
        -v ~/.ssh/id_rsa:/app/ssh_key:ro \
        -v backup:/app/backup:rw \
        ghcr.io/antonkosov/git-backups:latest
```

## Volume Mounts

| Host Path | Container Path | Mode | Description |
|-----------|----------------|------|-------------|
| `./config.yaml` | `/app/config.yaml` | `ro` | Configuration file |
| `~/.ssh/id_rsa` | `/app/ssh_key` | `ro` | Private SSH key (optional, path configurable) |
| `./backups` | `/app/backup` | `rw` | Backup destination directory (path configurable) |

**Note**: SSH key and backup paths can be customized in your configuration. You can mount multiple SSH keys to different container paths and specify multiple backup directories as needed.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details.
