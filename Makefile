ESBUILD=node_modules/.bin/esbuild
TSC=node_modules/.bin/tsc

.PHONY: build
build:
	$(ESBUILD) --bundle --outfile=dist/index.js src/index.ts

.PHONY: watch-ts
watch-ts:
	$(TSC) --watch

