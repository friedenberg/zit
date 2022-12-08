
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

.PHONY: build watch exclude bats_tests unit_tests go_vet graph_dependencies install;

# build: install unit_tests go_vet graph_dependencies;
build: install unit_tests go_vet;

go_build:
> go build -o build/zit ./.

go_vet: go_build
> go vet ./...

unit_tests:
> go test -timeout 5s ./...

install: go_build unit_tests go_build bats_tests
> go install ./.

bats_tests: go_build
> if [[ ! -f build_options/skip_bats_tests ]]; then
>   bats --jobs 8 zz-test/*.bats
> fi

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
