univiewd.toml(5)

; SPDX-FileCopyrightText: © Moritz Poldrack & AUTHORS
; SPDX-License-Identifier: CC-BY-SA-4.0

# LOCATION

When loading, the following directories are searched for a _univiewd.toml_ in
order.

- ./
- /etc/
- /config/
- ~/

# GENERAL

## Bind
Set the interface and port to bind to.

Type: _string_++
Format: <_interface_>:<_port_>++
Default: _127.1.2.4_:_1558_

## Host
Set the address or FQDN of the server. Used in connection links and
certificates.

Type: _string_++
Default: ""

# BUGS

Bugs can be reported at <https://todo.sr.ht/~mpldr/uniview>. Please consider
discussing them on our mailing list <~mpldr/uniview@lists.sr.ht> before filing
a ticket. (Archive: https://lists.sr.ht/~mpldr/uniview)

# SEE ALSO

*uniview.toml*(5)

# AUTHOR

See the list named AUTHORS in the project root or at
<https://git.sr.ht/~mpldr/uniview/tree/master/item/AUTHORS> for a list of
contributors.
