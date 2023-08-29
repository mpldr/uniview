---
toc: false
title: "Overview"
---

# `uniview://` URLs

The `uniview://` URL is an essential part to allow users to connect to any
instance. The basic structure is:

```
uniview://[host]/[room]?[insecure]#[password]
```

- `host` is the address of the server.
- `room` is the name of the room. It is required to URL-escape characters and
  sent as the path. The leading `/` is removed, and all characters are
  permitted in a room name (including line-breaks). UIs and servers may impose
  arbitrary limits on the characters permitted in a room name, but should
  always assume unwanted characters may have been sent.
- `insecure` is set as a query parameter, when the server certificate should
  not be validated.
- `password` is the password of the room. If no password is set for the room,
  it should be ignored.

## Why `?insecure` instead of a separate protocol

`http` has established the s-suffix for secure transport, but adding a second
schema would increase the burden on the user when selecting the default
application, therefore the query parameter has been used.

---
[![Creative Commons BY-SA](https://i.creativecommons.org/l/by-sa/4.0/80x15.png)](http://creativecommons.org/licenses/by-sa/4.0/)
