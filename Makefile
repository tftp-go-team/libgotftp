#!/usr/bin/make -f
#
# Subject: Makefile for libgotftp project.
#

.PHONY: test
test:
	$(MAKE) -C src test

.PHONY: doc
doc:		
