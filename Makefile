
SHELL := bash
.ONESHELL:
.SHELLFLAGS := -eu -o pipefail -c
.DELETE_ON_ERROR:
MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules
MAKEFLAGS += --output-sync=target
MAKEFLAGS := --jobs=$(shell nproc)


ifeq ($(origin .RECIPEPREFIX), undefined)
  $(error This Make does not support .RECIPEPREFIX. Please use GNU Make 4.0 or later)
endif
.RECIPEPREFIX = >

.PHONY: build watch exclude graph_dependencies

build: install;

.PHONY: install
install: tests_fast
> go install ./.

.PHONY: deploy
deploy: tests_slower
> go install ./.

.PHONY: go_generate
go_generate:
> go generate ./...

.PHONY: go_build
go_build: go_generate
> go build -o build/zit ./.

.PHONY: go_vet
go_vet:
> go vet ./...

.PHONY: tests_unit
tests_unit:
> go test -timeout 5s ./...

.PHONY: tests_fast
tests_fast: go_vet tests_unit;

.PHONY: tests_bats
tests_bats: go_build
> bats --jobs 8 zz-tests_bats/*.bats

.PHONY: tests_slow
tests_slow: tests_fast tests_bats;

.PHONY: tests_bats_migration
tests_bats_migration: go_build
> bats --jobs 8 zz-tests_bats/migration/*.bats

.PHONY: tests_slower
tests_slower: tests_fast tests_slow tests_bats_migration;

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
