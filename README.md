[![License](https://img.shields.io/github/license/Clov614/SteamHostSync)](./LICENSE)
[![Go](https://img.shields.io/badge/Go-project-00ADD8?logo=go&logoColor=white)](#usage)
[![Platform](https://img.shields.io/badge/Platforms-Windows%20%7C%20macOS%20%7C%20Linux-555)](#manual-setup)

# SteamHostSync

SteamHostSync is a small Go utility and hosts-file bundle for users who want a faster way to maintain host mappings for Steam, GitHub, Docker, GOG Galaxy, and related domains.

It can generate updated host files from `Source.yaml`, and the repository also publishes ready-to-use `Hosts*` files that can be consumed directly.

## Table of Contents
- [What it does](#what-it-does)
- [Usage](#usage)
- [Manual setup](#manual-setup)
- [Customize host sources](#customize-host-sources)
- [Repository structure](#repository-structure)
- [Notes](#notes)

## What it does

- resolves IPs for predefined platform domains
- writes combined and per-platform host files
- supports direct consumption through hosted raw files
- works well with tools such as SwitchHosts for automatic refreshes

## Usage

### Option 1, use with SwitchHosts

Recommended tool: [SwitchHosts](https://github.com/oldj/SwitchHosts)

Available subscription URLs:

**jsDelivr**
1. `https://cdn.jsdelivr.net/gh/Clov614/SteamHostSync@main/Hosts`
2. `https://cdn.jsdelivr.net/gh/Clov614/SteamHostSync@main/Hosts_steam`
3. `https://cdn.jsdelivr.net/gh/Clov614/SteamHostSync@main/Hosts_github`

**Statically backup**
4. `https://cdn.statically.io/gh/Clov614/SteamHostSync@main/Hosts`
5. `https://cdn.statically.io/gh/Clov614/SteamHostSync@main/Hosts_steam`
6. `https://cdn.statically.io/gh/Clov614/SteamHostSync@main/Hosts_github`

If one CDN is slow or unavailable, switch to the matching backup URL.

### Option 2, build and run locally

```bash
git clone https://github.com/Clov614/SteamHostSync.git
cd SteamHostSync
go run .
```

This reads `Source.yaml`, resolves the configured domains, and writes updated `Hosts` files.

## Manual setup

### Hosts file locations
- Windows: `C:\Windows\System32\drivers\etc\hosts`
- Linux: `/etc/hosts`
- macOS: `/etc/hosts`

### Apply the generated hosts entries
1. Copy the relevant entries into your local hosts file.
2. Save the file with administrator privileges.
3. Refresh DNS if changes do not take effect immediately.

### Refresh DNS
- Windows: `ipconfig /flushdns`
- Linux: `sudo nscd restart`
- macOS: `sudo killall -HUP mDNSResponder`

## Customize host sources

Edit `Source.yaml` to add or remove platform groups and domains.

Each platform entry uses this shape:

```yaml
platforms:
  - - github
    - alive.github.com
    - github.com
```

The first value is the group name. The remaining values are the domains to resolve.

## Repository structure

```text
main.go        Entry point for generating host files
fileIO/        Config and file helpers
Source.yaml    Domain source configuration
Hosts*         Generated output files
img/           Screenshots used in the repo
```

## Notes

- This project changes network resolution through the local hosts file, so review entries before applying them.
- Some domains may change IPs over time, which is why regenerating or refreshing the files matters.

## License

See [`LICENSE`](./LICENSE).
