# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## \[Unreleased\]

### Added

- Flag unmask action

## \[[v0.2](https://github.com/mbaraa/eloi/releases/tag/v0.2)\] 2023-03-06

### Added

- SQLite local database to store local repos and ebuilds
- Ebuild's description

### Changed

- restructured actions' execution path

### Fixed

- color output, added a wrapper that sets output's color, style, and background

## \[[v0.1.1](https://github.com/mbaraa/eloi/releases/tag/v0.1.1)\] - 2023-02-18

### Fixed

- finding a package with a different letters casing than the original package name
- package selection prompt appears when no packages are found
- duplicate packages search results
- synchronizing repos

## \[[v0.1](https://github.com/mbaraa/eloi/releases/tag/v0.1)\] - 2023-02-01

### Added

- added a changelog :)

### Changed

- changed the structure of the overlay model to represent some data more elegantly

### Fixed

- fixed overlay source type when adding a new overlay repo
