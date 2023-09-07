#!/bin/sh

# SPDX-FileCopyrightText: © Moritz Poldrack
# SPDX-License-Identifier: CC-BY-SA-4.0

set -e

mkdir windows
cd windows

inkscape -o .16.png -w 16 ../contrib/icon.svg
convert .16.png -thumbnail "16x16" -gravity center -background transparent -extent 16x16 .16s.png
inkscape -o .32.png -w 32 ../contrib/icon.svg
convert .32.png -thumbnail "32x32" -gravity center -background transparent -extent 32x32 .32s.png
inkscape -o .48.png -w 48 ../contrib/icon.svg
convert .48.png -thumbnail "48x48" -gravity center -background transparent -extent 48x48 .48s.png
inkscape -o .64.png -w 64 ../contrib/icon.svg
convert .64.png -thumbnail "64x64" -gravity center -background transparent -extent 64x64 .64s.png
inkscape -o .128.png -w 128 ../contrib/icon.svg
convert .128.png -thumbnail "128x128" -gravity center -background transparent -extent 128x128 .128s.png
inkscape -o .256.png -w 256 ../contrib/icon.svg
convert .256.png -thumbnail "256x256" -gravity center -background transparent -extent 256x256 .256s.png
convert .16s.png .32s.png .48s.png .64s.png .128s.png .256s.png icon.ico
go run github.com/josephspurrier/goversioninfo/cmd/goversioninfo \
	-copyright="Moritz Poldrack and AUTHORS" \
	-description="uniview client and server bundle" \
	-64 \
	-product-version="0.3.0" \
	-icon icon.ico \
	-skip-versioninfo \
	-o=resource.syso

GOOS=windows make -C .. univiewd.exe "EXTRA_GO_LDFLAGS= -H=windowsgui"

osslsigncode sign \
	-certs ../signcert.pem \
	-key ../signkey.pem \
	-pass "$(secret-tool lookup Title "Code Signing Cert")" \
	-n Uniview \
	-i https://moritz.sh \
	-h sha256 \
	-t http://timestamp.digicert.com \
	-in ../uniview.exe -out uniview.exe
osslsigncode sign \
	-certs ../signcert.pem \
	-key ../signkey.pem \
	-pass "$(secret-tool lookup Title "Code Signing Cert")" \
	-n Uniview \
	-i https://moritz.sh \
	-h sha256 \
	-t http://timestamp.digicert.com \
	-in ../univiewd.exe -out univiewd.exe

echo copying source…
mkdir src
(cd .. && git archive HEAD | tar xC windows/src/)

curl -o AGPL.rtf https://www.gnu.org/licenses/agpl-3.0.rtf
curl -o mpv.7z https://nav.dl.sourceforge.net/project/mpv-player-windows/release/mpv-0.36.0-x86_64.7z
7z x -ompv mpv.7z

echo creating installer…
makensis -V3 -WX -INPUTCHARSET UTF8 "-DPWD=$(pwd)" ../contrib/nsis.nsi

echo signing installer…
osslsigncode sign \
	-certs ../signcert.pem \
	-key ../signkey.pem \
	-pass "$(secret-tool lookup Title "Code Signing Cert")" \
	-n Uniview \
	-i https://moritz.sh \
	-h sha256 \
	-t http://timestamp.digicert.com \
	-in uniview-setup.exe -out ../uniview-setup.exe
