# SPDX-FileCopyrightText: © nobody
# SPDX-License-Identifier: CC0-1.0

WEBSRC!=find html/ -type f
WEBTARGETS:=$(subst html/,dist/,$(WEBSRC))

.PHONY: dist
dist: $(WEBTARGETS) dist/logo.svg dist/bootstrap.min.css

dist/%: html/%
	go run github.com/tdewolff/minify/v2/cmd/minify -o $@ $<

html/logo.svg:
	curl -o $@ https://git.sr.ht/~mpldr/uniview/blob/master/contrib/icon.svg

html/bootstrap.min.css:
	curl -o $@ https://cdn.jsdelivr.net/npm/bootstrap@5.3.1/dist/css/bootstrap.min.css
