# Eloi

[![GoDoc](https://godoc.org/github.com/mbaraa/eloi?status.png)](https://godoc.org/github.com/mbaraa/eloi)
[![build](https://github.com/mbaraa/eloi/actions/workflows/build.yml/badge.svg)](https://github.com/mbaraa/eloi/actions/workflows/build.yml)
[![coverage](https://github.com/mbaraa/eloi/actions/workflows/coverage.yml/badge.svg)](https://github.com/mbaraa/eloi/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/mbaraa/eloi)](https://goreportcard.com/report/github.com/mbaraa/eloi)

A Gentoo ebuilds searcher and installer ([eix](https://wiki.gentoo.org/wiki/Eix) with extra steps). Searches through all Gentoo's overlays provided by [Layman](https://wiki.gentoo.org/wiki/Layman) and listed by Zugania's [website](http://gpo.zugaina.org/).

Eloi's [server](https://github.com/mbaraa/eloi-server) just scrapes over Zugania's website and provides a list of overlays to be used by the CLI client.

---

## Features:

- Find and install an ebuild package from any overlay
- Enable ebuild overlays
- more in the future...

## Dependencies:

- [gentoo system](https://gentoo.org) -- duh
- [gentoolkit](https://wiki.gentoo.org/wiki/Gentoolkit)
- [python3.10](https://wiki.gentoo.org/wiki/Python) -- install python3.10 

## Installation:

### Using Portage

1. Add my [overlay](https://github.com/mbaraa/mbaraa-overlay) 🥰
2. Install `app-portage/eloi` using emerge

### Using Go's installer

```bash
go install github.com/mbaraa/eloi@latest
```

## Usage:

#### Update local repos' cache

local repos are stored in `/var/cache/eloi`

```bash
eloi --download
```

#### Find an Ebuild

using `-S` or `--search`

```bash
eloi -S pulseaudio-equalizer
```

This will list all ebuild that has _pulseaudio-equalizer_ in their name, with some other details, like version, overlay name, use flags, license, ...

#### Enable an Overlay repository

This can be done using the `--enable` flag

```bash
eloi --enable underworld
```

or by installing a package from a repo that's not enabled in your system

```bash
eloi -qav pulseaudio-equalizer-ladspa
```

#### Sync portage repos

```bash
eloi --sync
```
