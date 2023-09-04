uniview.toml(5)

; SPDX-FileCopyrightText: © Moritz Poldrack & AUTHORS
; SPDX-License-Identifier: CC-BY-SA-4.0

# LOCATION

When loading, the following files are searched in order.

- $XDG_CONFIG_HOME/
- ~/.univiewd.toml

# Media

## Directories
Set the directories to scan for files. These are also made available via the
Client-API.

Type: _string[]_++
Format: ["dir1","dir2",…,"dirn"]++
Default: ["~/Videos"]

# BUGS

Bugs can be reported at <https://todo.sr.ht/~mpldr/uniview>. Please consider
discussing them on our mailing list <~mpldr/uniview@lists.sr.ht> before filing
a ticket. (Archive: https://lists.sr.ht/~mpldr/uniview)

# SEE ALSO

*univiewd.toml*(5)

# AUTHOR

See the list named AUTHORS in the project root or at
<https://git.sr.ht/~mpldr/uniview/tree/master/item/AUTHORS> for a list of
contributors.
