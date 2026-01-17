# compose-volume-migrate

Migrate Docker Compose external volumes between hosts using tar.gz archives.

This tool assumes that volumes without external: true are Compose-managed ephemeral data (cache, temporary files) that don't need migration.

## Features

- Automatically detects external volumes from `docker-compose.yml`
- Exports/imports only volumes with `external: true`

## Installation

```bash
$ go install github.com/zinrai/compose-volume-migrate@latest
```

## Migration Workflow

**Source host:**

```bash
$ cd /path/to/project
$ compose-volume-migrate export
```

**Transfer:**

```bash
$ scp -r /path/to/project user@target:/path/to/
```

**Target host:**

```bash
$ cd /path/to/project
$ compose-volume-migrate import
$ docker compose up -d
```

## How It Works

### Export

1. Parses `docker-compose.yml` and filters `external: true` volumes
2. Checks for existing tar.gz files (fails if exist)
3. Checks if each volume is in use (fails if running container detected)
4. Exports each volume: `busybox:stable-glibc` + `tar`

### Import

1. Parses `docker-compose.yml` and filters `external: true` volumes
2. Checks for missing tar.gz files (fails if missing)
3. Creates volumes if they don't exist
4. Imports each volume: `busybox:stable-glibc` + `tar`

## License

This project is licensed under the [MIT License](./LICENSE).
