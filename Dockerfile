FROM golang:alpine AS build
RUN apk add make
COPY . /src
WORKDIR /src
RUN CGO_ENABLED=0 EXTRA_GO_LDFLAGS="-s -w" make univiewd

FROM scratch
COPY --from=build /src/univiewd /
ENTRYPOINT /univiewd
