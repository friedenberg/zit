
.PHONY: build watch exclude bats_tests unit_tests go_vet graph_dependencies install;

build: install unit_tests go_vet graph_dependencies;
# build: install unit_tests go_vet;

go_build:
	go build -o build/zit ./.

go_vet: go_build
	go vet ./...

unit_tests:
	go test ./...

install: bats_tests
	go install ./.

bats_tests: go_build
	bats --jobs 8 zz-test/*.bats

graph_dependencies:
	./bin/graph_dependencies

watch:
	echo .

exclude:
	echo .DS_Store
	echo zit/.git/
	echo zit/\.zit/
	echo build/
	echo zit/zit$$
