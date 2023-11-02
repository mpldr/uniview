# SPDX-FileCopyrightText: © nobody
# SPDX-License-Identifier: CC0-1.0

FROM archlinux AS build
RUN pacman -Syu make go git protobuf which --noconfirm
COPY . /src
WORKDIR /src
RUN CGO_ENABLED=0 EXTRA_GO_LDFLAGS="-s -w" make

FROM scratch
COPY --from=build /src/univiewd /
CMD ["/univiewd"]
