# uniview

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

- Synchronise upon connection
- Support playback queues
- TLS
- API docs
- CI Pipeline
- maybe a UI
- maybe a webinterface with stats
- bugsquashing

## Licence
<!--    ↑ this is for you, rock -->

This thing's licensed under the [AGPL](./LICENSES/AGPL-3.0-or-later.txt).

[<img alt="AGPL logo" src="https://upload.wikimedia.org/wikipedia/commons/0/06/AGPLv3_Logo.svg" height="40">](./LICENSES/AGPL-3.0-or-later.txt)
