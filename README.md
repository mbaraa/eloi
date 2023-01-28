# Eloi

A Gentoo ebuilds searcher and installer ([eix](https://wiki.gentoo.org/wiki/Eix) with extra steps). Searches through all Gentoo's overlays provided by [Layman](https://wiki.gentoo.org/wiki/Layman) and listed by Zugania's [website](http://gpo.zugaina.org/).

Eloi's [server](https://github.com/mbaraa/eloi-server) just scrapes over Zugania's website and provides a list of overlays to be used by the CLI client.

---

## Features:

- Search for an ebuild package
- Enable ebuild overlays
- more in the futire...

## Usage:

#### Update local repos

local repos are stored in `/var/cache/eloi`

```bash
eloi --sync
```

#### Find an Ebuild

```bash
eloi -R pulseaudio-equalizer
```

This will list all ebuild that has _pulseaudio-equalizer_ in their name, with some other details, like version, overlay name, use flags, licens, ...

#### Enable an Overlay repository

This can be done using the `--enable` flag

```bash
eloi --enable underworld
```

or by installing a package from a repo that's not enabled in your system

```bash
eloi -qav pulseaudio-equalizer-ladspa
```
