#!/bin/sh

find \
	./sample.d \
	-type f \
	-name '*.xlsx' |
	sed 's,./sample.d,/guest.d,' |
	wazero \
		run \
		-mount "${PWD}/sample.d:/guest.d:ro" \
		./xfiles2shnames.wasm |
	jq -c
