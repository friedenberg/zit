
.PHONY: build watch exclude;

build:
	go build -o build/zit ./.
	go vet ./...
	go test ./...
	go install ./.
	bats zz-test/*.bats
	./bin/graph_dependencies

watch:
	echo .

exclude:
	echo dot
	echo dot.svg
	echo .DS_Store
	echo zit/.git/
	echo zit/\.zit/
	echo build/
	echo zit/zit$$
