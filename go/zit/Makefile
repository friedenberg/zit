
SHELL := bash
.ONESHELL:
.SHELLFLAGS := -eu -o pipefail -c
.DELETE_ON_ERROR:
MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules
MAKEFLAGS += --output-sync=target
n_prc := $(shell sysctl -n hw.logicalcpu)
MAKEFLAGS := --jobs=$(n_prc)
timeout := 10
cmd_bats := BATS_TEST_TIMEOUT=$(timeout) bats --tap

ifeq ($(origin .RECIPEPREFIX), undefined)
				$(error This Make does not support .RECIPEPREFIX. Please use GNU Make 4.0 or later)
endif
.RECIPEPREFIX = >

.PHONY: build watch exclude graph_dependencies

build: build/zit;

files_go_generate := $(shell ag 'go:generate' -l src/)

build/go_generate: $(files_go_generate)
> go generate ./...
> touch "$@"

files_go := $(shell find src -type f)

build/zit: build/go_generate $(files_go)
> go build -o build/zit ./.

build/go_vet: $(files_go)
> go vet ./... || true
> touch "$@"

dirs_go_unit := $(shell find src -mindepth 2 -iname '*_test.go' | sort -u)

build/tests_unit: $(files_go) build/go_generate
> go test -v ./...
> touch "$@"

build/tests_fast: build/go_vet build/tests_unit
> @touch "$@"

files_tests_bats := $(shell find zz-tests_bats -type f)

build/tests_bats: build/zit $(files_tests_bats) build/tests_gen_fixture
> $(cmd_bats) zz-tests_bats/*.bats
> touch "$@"

files_tests_gen_fixture := $(shell find zz-tests_bats/migration)

build/tests_gen_fixture: build/zit $(files_tests_gen_fixture)
> ./zz-tests_bats/migration/generate_fixture.bash "$$(mktemp -d)" >/dev/null 2>&1
> touch "$@"

build/tests_slow: build/tests_fast build/tests_bats
> touch "$@"

files_tests_bats_migration := $(shell find zz-tests_bats/migration)

build/tests_bats_migration: build/zit $(files_tests_bats_migration)
> $(cmd_bats) zz-tests_bats/migration/*.bats
> touch "$@"

build/tests_slower: build/tests_fast build/tests_slow build/tests_bats_migration
> touch "$@"

build/tests_slowest: build/tests_fast build/tests_slow build/tests_bats_migration
> touch "$@"

build/deploy: build/tests_slowest;

graph_dependencies:
> ./bin/graph_dependencies

build/gen_fixture: build/zit $(files_tests_gen_fixture) build/zit
> ./zz-tests_bats/migration/generate_fixture.bash
> touch "$@"
.PHONY: build/gen_fixture

watch:
> echo .

exclude:
> echo .DS_Store
> echo .git/
> echo zit/.git/
> echo zit/\.zit/
> echo build/
> echo zit/zit$$
