
SHELL := bash
.ONESHELL:
.SHELLFLAGS := -eu -o pipefail -c
.DELETE_ON_ERROR:
MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules

ifeq ($(origin .RECIPEPREFIX), undefined)
  $(error This Make does not support .RECIPEPREFIX. Please use GNU Make 4.0 or later)
endif
.RECIPEPREFIX = >

.PHONY: build watch exclude graph_dependencies install;

build: install;

install: tests_fast
> go install ./.

deploy: tests_slower
> go install ./.
.PHONY: deploy;

go_generate:
> go generate ./...
.PHONY: go_generate;

go_build: go_generate
> go build -o build/zit ./.
.PHONY: go_build;

go_vet: go_build
> go vet ./...
.PHONY: go_vet;

tests_unit:
> go test -timeout 5s ./...

tests_fast: go_vet tests_unit;
.PHONY: tests_fast;

tests_bats: go_build
> if [[ ! -f build_options/skip_bats_tests ]]; then
>   bats --jobs 8 zz-test/*.bats
> fi
.PHONY: tests_bats;

tests_slow: tests_fast tests_bats;
.PHONY: tests_slow;

tests_bats_migration: go_build;
.PHONY: tests_bats_migration;

tests_slower: tests_fast tests_slow tests_bats_migration;
.PHONY: tests_slower;

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
