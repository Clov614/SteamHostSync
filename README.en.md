# SteamHostSync

[![Stars](https://img.shields.io/github/stars/Clov614/SteamHostSync?style=flat-square)](https://github.com/Clov614/SteamHostSync/stargazers)
[![Go](https://img.shields.io/badge/go-1.22%2B-00ADD8?style=flat-square&logo=go)](https://go.dev/)
[![License](https://img.shields.io/github/license/Clov614/SteamHostSync?style=flat-square)](https://github.com/Clov614/SteamHostSync/blob/main/LICENSE)

SteamHostSync is a small Go utility and ready-to-use hosts source for syncing Steam, GitHub, Docker, GOG, Ubisoft, and related domains.

It is aimed at users who want to keep an organized, updatable hosts configuration with SwitchHosts or by hand.

> **Languages / 语言**: [English](README.en.md) · [简体中文](README.zh-CN.md)

## Table of Contents

- [Why This Repo Exists](#why-this-repo-exists)
- [What It Generates](#what-it-generates)
- [Quick Start](#quick-start)
- [Manual Setup](#manual-setup)
- [Refresh DNS](#refresh-dns)
- [Customize config.yaml](#customize-configyaml)
- [Fork Customization (making your fork self-updating)](#fork-customization-making-your-fork-self-updating)
- [Build From Source](#build-from-source)
- [Repository Layout](#repository-layout)

## Why This Repo Exists

This project periodically resolves the IPs for common platform domains and publishes ready-to-use hosts files, so you can reach these services more reliably:

- Steam
- GitHub
- Docker Hub
- GOG Galaxy
- Ubisoft download services

## What It Generates

The repo publishes multiple host lists so you can subscribe to only what you need:

- `Hosts` for the full bundle
- `Hosts_steam` for Steam only
- `Hosts_steam_linux` for Linux Steam only (apt repo + client update domains)
- `Hosts_github` for GitHub only
- `Hosts_docker` for Docker-related domains
- `Hosts_gog` for GOG services
- `Hosts_ubisoft` for Ubisoft download services

## Quick Start

### Option 1: Use SwitchHosts (recommended)

[SwitchHosts](https://github.com/oldj/SwitchHosts) is the easiest way to keep the hosts file updated automatically.

Available subscription URLs:

**Primary (jsDelivr)**

1. ALL: `https://cdn.jsdelivr.net/gh/Clov614/SteamHostSync@main/Hosts`
2. Steam: `https://cdn.jsdelivr.net/gh/Clov614/SteamHostSync@main/Hosts_steam`
3. Linux Steam: `https://cdn.jsdelivr.net/gh/Clov614/SteamHostSync@main/Hosts_steam_linux`
4. GitHub: `https://cdn.jsdelivr.net/gh/Clov614/SteamHostSync@main/Hosts_github`

**Fallback (Statically)**

5. ALL: `https://cdn.statically.io/gh/Clov614/SteamHostSync@main/Hosts`
6. Steam: `https://cdn.statically.io/gh/Clov614/SteamHostSync@main/Hosts_steam`
7. Linux Steam: `https://cdn.statically.io/gh/Clov614/SteamHostSync@main/Hosts_steam_linux`
8. GitHub: `https://cdn.statically.io/gh/Clov614/SteamHostSync@main/Hosts_github`

If the primary CDN is unreachable, switch to the matching fallback URL.

### Optional: Auto-start SwitchHosts on Windows

1. Press `Win + R`
2. Run `shell:startup`
3. Copy your SwitchHosts shortcut into that folder

![](/img/1.png)
![](/img/2.png)

### Configure automatic updates in SwitchHosts

![](/img/3.png)

## Manual Setup

### Hosts file locations

- Windows: `C:\Windows\System32\drivers\etc\hosts`
- Linux: `/etc/hosts`
- macOS: `/etc/hosts`

### Apply manually

Append the contents of one of the generated files, such as `Hosts` or `Hosts_steam`, to the end of your system hosts file.

Example:

```bash
# Linux / macOS
cat Hosts_steam | sudo tee -a /etc/hosts
```

Or open the file directly and paste the generated entries at the bottom.

## Refresh DNS

Most systems pick up the change quickly. If not, refresh DNS:

1. Windows: `ipconfig /flushdns`
2. Linux: `sudo nscd restart`
3. macOS: `sudo killall -HUP mDNSResponder`

## Customize config.yaml

Running the executable once will generate `config.yaml`, which defines the platforms and domains to resolve, the DNS-over-HTTPS upstreams, concurrency and timeouts.

A simplified example:

```yaml
version: 1
concurrency: 8
timeout:
  resolve: 5s
  probe: 2s
probe:
  port: 443
  attempts: 3
dns_servers:
  - https://dns.alidns.com/resolve
  - https://doh.pub/dns-query
  - https://dns.google/resolve
platforms:
  - name: github
    domains:
      - alive.github.com
      - github.com
      - raw.githubusercontent.com
  - name: steam
    domains:
      - steamcommunity.com
      - store.steampowered.com
```

Notes:

- Each `platforms` entry produces a `Hosts_<name>` file; all platforms combined produce `Hosts`.
- Platform names are lowercased into safe filenames (illegal characters are replaced with `_`); lowercase names without spaces are recommended.
- Domains that fail to resolve are kept as `# domain` comment lines and do not affect the other entries.
- If `config.yaml` does not exist, the embedded default configuration is written automatically on the first run.

## Fork Customization (making your fork self-updating)

> The following is for users who fork this repo to serve as their own hosts source, not for developer collaboration.

After forking, the built-in CI (`.github/workflows/Update.yml`) does not run automatically. Configure it once as below, and it will refresh your hosts files every 12 hours from then on:

### 1. Enable Actions

The `schedule` trigger on a fork is **disabled by default**, so enable it manually:

1. Open your fork and click the **Actions** tab.
2. The first time you will be prompted to enable Actions — click **I understand my workflows, go ahead and enable them**.
3. If the `Update` workflow shows as disabled in the sidebar, open its "…" menu and choose **Enable workflow** (this also runs it once and resets the timer).

> ⚠️ **It will be disabled again after 60 days of inactivity**: GitHub automatically disables scheduled workflows in public repositories that have had no activity for 60 consecutive days. As long as your fork receives pushes (or you trigger runs manually), it stays alive; if it does get disabled, just repeat the steps above.

### 2. Grant the GITHUB_TOKEN write access (critical)

`Update.yml` `git push`es the generated hosts files back to your `main` branch, which requires the repo's default `GITHUB_TOKEN` to have write access. New repositories default to read-only, otherwise CI fails with `Permission denied` on push.

1. Open **Settings → Actions → General**.
2. Under **Workflow permissions**, select **Read and write permissions** and save.

### 3. Customize the platforms and domains (optional)

Edit the root `config.yaml` (`platforms` / `dns_servers` / `concurrency` / `timeout` / `probe`); field descriptions are in the previous section. Each CI run commits this file as well, so your customization is preserved in git history.

> Note: do not manually commit the CI-generated `Hosts*` and `README.md` files — the workflow maintains them. If you also build locally, remember to mirror the same config into `internal/config/default_config.yaml` (see `CLAUDE.md`).

### 4. Update your subscription URLs (if you use SwitchHosts)

Replace `Clov614` with your own GitHub username in every CDN link in the README:

```
https://cdn.jsdelivr.net/gh/<your-username>/SteamHostSync@main/Hosts
```

> Note: jsDelivr / Statically only serve content from public repositories, so keep your fork public; otherwise the subscription URLs will not be accessible.

### 5. (Optional) Publish a release

Without a `v*` tag, `release.yml` (goreleaser) never triggers and no configuration is needed. To publish your own release, just push a `v*` tag; the workflow publishes automatically using `GITHUB_TOKEN`, no extra secrets required.

### Summary of triggers

- Scheduled: every 12 hours (requires step 1)
- `push` to `main`
- Manual: Actions tab → `Update` → **Run workflow**

## Build From Source

```bash
git clone https://github.com/Clov614/SteamHostSync.git
cd SteamHostSync
go run .
```

This updates the generated host files in the repository root and refreshes `README.md` from `README_TEMP.md` when present.

## Repository Layout

```text
main.go                  Entry point (flags + orchestration)
internal/                Layered packages: config / resolve / probe / render / fileio / app
config.yaml              Runtime configuration (auto-generated on first run)
Hosts*                   Generated host output files
img/                     README screenshots
README_TEMP.md           Template used when regenerating README (host content placeholder)
```
