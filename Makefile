# evids
# See LICENSE for copyright and license details.
.POSIX:

PREFIX ?= /usr
GO ?= go
GOFLAGS ?= -buildvcs=false
RM ?= rm -f

all: evids

evids:
	$(GO) build $(GOFLAGS)

clean:
	$(RM) evids

install: all
	mkdir -p $(DESTDIR)$(PREFIX)/sbin
	cp -f evids $(DESTDIR)$(PREFIX)/sbin
	chmod 755 $(DESTDIR)$(PREFIX)/sbin/evids

uninstall:
	$(RM) $(DESTDIR)$(PREFIX)/sbin/evids

.DEFAULT_GOAL := all

.PHONY: all evids clean install uninstall
