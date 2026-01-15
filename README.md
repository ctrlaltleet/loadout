# loadout

**loadout** is a no-nonsense toolset manager for pentesters and operators. It automates downloading + verifying tool binaries, and repos so you don't waste time hunting URLs or manual installs.

---

## Why loadout?

During engagements, you need tools *now*. `loadout` ensures you always have your toolbox ready, verified, and platform-appropriate, with no manual searching, no mistakes, no delays.

* Stop Googling tool URLs during engagements
* Define your toolkit in a simple YAML file
* Download verified binaries or clone git repos fast
* Works cross-platform: Linux, Windows, macOS
* Download multiple tools in parallel
* Select tools by name, tags, or just grab them all

---

## Quick Start

1. **Create your config YAML**:

```yaml
packages:
  nmap:
    version: "7.93"
    tags: [scanner]
    global_assets:
      some_asset_key:
        url: "git://github.com/emadshanab/Nmap-NSE-scripts-collection"
    platform_assets:
      darwin_arm64:
        url: "https://nmap.org/dist/nmap-7.98-setup.exe"
        hash: "sha256:none"
```

2. **List available tools for your platform:**

```bash
./loadout -config config.yaml -list
```

3. **Download selected tools:**

```bash
./loadout -config config.yaml -select nmap
```

4. **Default downloads go to `./downloads`. Change with `-output-dir`.**

---

## Usage

```
$ loadout
Usage of loadout:
  -concurrency int
        number of parallel downloads (default 4) (default 4)
  -config string
        path to config file (required)
  -list
        list packages for current platform
  -output-dir string
        output directory (default ./downloads) (default "./downloads")
  -platform string
        override platform (GOOS_GOARCH), e.g. linux_amd64
  -select string
        comma-separated package names and/or tags to select (or 'all')
```

---

## Supported Asset Types

* **HTTP/HTTPS downloads** — verified with SHA256, SHA512, or MD5 hashes
* **Git repos** — cloned from `git://` URLs

---

## Platform Support

Builds and runs on:

* Linux amd64 and arm64
* Windows amd64 and arm64
* macOS amd64 and arm64

---

## Build

```bash
make all
```

Binaries land in `build/`.
