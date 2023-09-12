# SPDX-FileCopyrightText: © nobody
# SPDX-License-Identifier: CC0-1.0

WEBSRC!=find html/ -type f
WEBTARGETS:=$(subst html/,dist/,$(WEBSRC))

.PHONY: dist
dist: $(WEBTARGETS) dist/logo.svg

dist/%: html/%
	go run github.com/tdewolff/minify/v2/cmd/minify -o $@ $<

html/logo.svg:
	curl -o $@ https://git.sr.ht/~mpldr/uniview/blob/master/contrib/icon.svg
