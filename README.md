# uniview

[![GitHub tag (with filter)](https://img.shields.io/github/v/tag/mpldr/uniview?label=version)]()
![Licence: AGPL](https://img.shields.io/badge/-AGPL--3-green?logo=opensourceinitiative&label=License&cacheSeconds=31536000)
![Demo available under uv.mpldr.de](https://img.shields.io/badge/-uv.mpldr.de-blue?label=Demo&cacheSeconds=31536000)
![GitHub code size in bytes](https://img.shields.io/github/languages/code-size/mpldr/uniview)
![![REUSE status](https://api.reuse.software/badge/git.sr.ht/~mpldr/uniview)](https://api.reuse.software/info/git.sr.ht/~mpldr/uniview)
[![Liberapay receiving](https://img.shields.io/liberapay/receives/mpldr)](https://liberapay.com/mpldr)
[![GitHub Sponsors](https://img.shields.io/github/sponsors/mpldr?logo=github&color=lightgrey)](https://github.com/sponsors/mpldr)

<img alt="Uniview Logo" src="https://git.sr.ht/~mpldr/uniview/blob/master/contrib/icon.svg" height="64">

This program syncronises video playback across multiple mpv instances.

## Building it

Install Go and you should be good to go.

```bash
make
```

## A single binary to rule them all

Server and client are the same binary, so hardlinking them works fine. If you
name the binary `univiewd` and run it, it will open a server on `:1558` you can
connect to.

## Roadmap

- Support playback queues
- API docs
- CI Pipeline
- a UI for managing the queue
- better handling of web streams
- bugsquashing

## Licence
<!--    ↑ this is for you, rock -->

This thing's licensed under the [AGPL](./LICENSES/AGPL-3.0-or-later.txt).

[<img alt="AGPL logo" src="https://upload.wikimedia.org/wikipedia/commons/0/06/AGPLv3_Logo.svg" height="40">](./LICENSES/AGPL-3.0-or-later.txt)
