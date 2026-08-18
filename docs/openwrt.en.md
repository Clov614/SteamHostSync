# Automatic SteamHostSync updates on OpenWrt

The OpenWrt integration loads a dedicated hosts file through dnsmasq `addnhosts`. It never appends to, replaces, or takes ownership of `/etc/hosts`. The updater runs briefly as a scheduled task; there is no resident SteamHostSync process.

## Scope

The integration targets OpenWrt installations using UCI and dnsmasq. Hosts state and temporary output are stored under `/var/lib/steamhostsync`; the update lock is `/var/lock/steamhostsync.lock`. On OpenWrt, `/var` is normally memory-backed, avoiding periodic flash writes.

Two source modes are available:

- `remote` (default): download a generated `Hosts` or `Hosts_<platform>` file. This has the lowest resource cost and is CPU-architecture independent.
- `local` (optional): run SteamHostSync on the router to perform DoH resolution and TCP probing, then publish the result through the same validation and rollback path.

The Steam++-style model of mapping hosts to loopback and proxying HTTPS is out of scope. This integration does not install certificates, listen on ports 80/443, or modify firewall rules.

## Build and install the OpenWrt package

Copy `contrib/openwrt` into an OpenWrt SDK or source tree as `package/steamhostsync-openwrt`, then run:

```sh
make package/steamhostsync-openwrt/compile V=s
```

Upload and install the resulting package:

```sh
opkg install steamhostsync-openwrt_*.ipk
steamhostsyncctl setup
steamhostsyncctl update
steamhostsyncctl status
```

`setup` performs these idempotent operations:

1. Creates `/var/lib/steamhostsync/Hosts`.
2. Adds that exact path to the dnsmasq `addnhosts` list without removing existing entries.
3. Installs a twice-daily cron task.
4. Enables a boot-time one-shot update service.

Running `setup` again does not duplicate UCI or cron entries.

## Mode A: remote subscription

The default `/etc/config/steamhostsync` contains:

```text
config steamhostsync 'main'
        option source 'remote'
        option url 'https://cdn.jsdelivr.net/gh/Clov614/SteamHostSync@main/Hosts'
        option fallback_url 'https://cdn.statically.io/gh/Clov614/SteamHostSync@main/Hosts'
```

To subscribe only to Steam:

```sh
uci set steamhostsync.main.url='https://cdn.jsdelivr.net/gh/Clov614/SteamHostSync@main/Hosts_steam'
uci set steamhostsync.main.fallback_url='https://cdn.statically.io/gh/Clov614/SteamHostSync@main/Hosts_steam'
uci commit steamhostsync
steamhostsyncctl update
```

Downloaded content must:

- Use HTTPS with certificate verification.
- Stay below the configured `max_bytes` limit.
- Contain the SteamHostSync header and at least one valid record.
- Contain only valid public IPv4 and hosts domain entries.
- Use domains listed in `/etc/steamhostsync/allowed-domains` when strict mode is enabled.

By default, both the primary and fallback sources must download successfully and produce byte-identical content before publication. This reduces the risk of a single CDN cache or path being poisoned. A download or validation failure leaves the current file untouched and does not reload dnsmasq. The repository publisher and GitHub Actions remain the content trust root.

For custom domains, extend the allowlist. Setting `strict_domains='0'` is possible, but it allows the remote source to control arbitrary LAN DNS names and is not recommended.

## Mode B: local generation

Local mode requires a static SteamHostSync binary matching the router CPU. Releases provide:

- `SteamHostSync_linux_amd64.tar.gz`
- `SteamHostSync_linux_386.tar.gz`
- `SteamHostSync_linux_armv7.tar.gz`
- `SteamHostSync_linux_arm64.tar.gz`
- `SteamHostSync_linux_mips_softfloat.tar.gz`
- `SteamHostSync_linux_mipsle_softfloat.tar.gz`

For older OpenWrt MIPS devices, verify endianness: `mips` is big-endian and `mipsle` is little-endian. An incorrect artifact normally fails with `Exec format error` or `Illegal instruction`.

Install the extracted binary and enable local mode:

```sh
install -m 0755 SteamHostSync /usr/bin/steamhostsync
steamhostsyncctl mode local
steamhostsyncctl update
```

Local mode uses `/etc/steamhostsync/router-config.yaml`, which defaults to:

- concurrency `1`
- one probe attempt per IP
- a limited set of core GitHub, Steam, and Linux Steam domains
- no README generation
- temporary output under `/var/lib/steamhostsync`

When adding domains, update both the router YAML and the allowlist. The Go process runs only during an update and exits afterward.

Resource-constrained devices should keep using mode A. MIPS/MIPSLE artifacts are cross-compiled, but memory use and instruction compatibility must still be verified on the target hardware.

## Update, rollback, and disable

```sh
# Update immediately
steamhostsyncctl update

# Show the mode, state path, and dnsmasq integration status
steamhostsyncctl status

# Restore the previous valid version from the current boot
steamhostsyncctl rollback

# Remove this package's cron and dnsmasq entries, disable updates, and disable init
steamhostsyncctl disable
```

An atomic lock prevents boot, cron, and manual updates from running concurrently. New content is downloaded or generated into a temporary path and fully validated before publication. A failed dnsmasq reload automatically restores the previous version.

`Hosts` and `Hosts.previous` are memory-backed and disappear after reboot. The init service fetches them again. If WAN is unavailable at boot, ordinary DNS keeps working, but SteamHostSync entries are unavailable until a later successful update.

## Troubleshooting

### LAN clients do not use the records

Ensure clients use the router as their DNS server. Browser or operating-system DoH/DoT can bypass dnsmasq and therefore bypass this integration.

### The router uses AdGuardHome

Automatic package integration targets dnsmasq only and does not modify AdGuardHome. Add the project's `Hosts` URL directly to AdGuardHome as a hosts-format custom filtering source instead.

### HTTPS downloads fail

Ensure the router clock is correct and `ca-bundle` is installed. The updater will not disable certificate verification to work around TLS errors.

### Hosts do not fully solve Steam connectivity

Hosts entries only affect name resolution. Region restrictions, TLS behavior, proxy requirements, client networking, and upstream service failures are outside their scope.
