
SHELL := bash
.ONESHELL:
.SHELLFLAGS := -eu -o pipefail -c
.DELETE_ON_ERROR:
MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules
MAKEFLAGS += --output-sync=target
n_prc := $(shell sysctl -n hw.logicalcpu)
MAKEFLAGS := --jobs=$(n_prc)
cmd_bats := bats --jobs $(n_prc)

ifeq ($(origin .RECIPEPREFIX), undefined)
				$(error This Make does not support .RECIPEPREFIX. Please use GNU Make 4.0 or later)
endif
.RECIPEPREFIX = >

.PHONY: build watch exclude graph_dependencies

build: install;

.PHONY: install
install: build/tests_fast
> go install ./.

.PHONY: deploy
deploy: build/tests_slower
> go install ./.

files_go_generate := $(shell ag 'go:generate' -l)

build/go_generate: $(files_go_generate)
> go generate ./...
> touch "$@"

files_go := $(shell find src -type f)

build/zit: build/go_generate $(files_go)
> go build -o build/zit ./.

build/go_vet: $(files_go)
> go vet ./...
> touch "$@"

build/tests_unit: $(files_go) build/go_generate
> go test -timeout 5s ./...
> touch "$@"

build/tests_fast: build/go_vet build/tests_unit
> touch "$@"

files_tests_bats := $(shell find zz-tests_bats -type f)

build/tests_bats: build/zit $(files_tests_bats)
> $(cmd_bats) zz-tests_bats/*.bats
> touch "$@"

build/tests_slow: build/tests_fast build/tests_bats
> touch "$@"

build/tests_bats_migration: build/zit
> $(cmd_bats) zz-tests_bats/migration/*.bats
> touch "$@"

build/tests_slower: build/tests_fast build/tests_slow build/tests_bats_migration;
> touch "$@"

graph_dependencies:
> ./bin/graph_dependencies

watch:
> echo .

exclude:
> echo .DS_Store
> echo .git/
> echo zit/.git/
> echo zit/\.zit/
> echo build/
> echo zit/zit$$
