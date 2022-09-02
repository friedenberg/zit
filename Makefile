
.PHONY: build watch exclude;

build:
	./bin/graph_dependencies src/bravo
	# go build -o build/zit ./.
	# go vet ./...
	# go test ./...
	# go install ./.
	# bats zz-test/*.bats

watch:
	echo dot
	echo ./bin/graph_dependencies

exclude:
	echo .DS_Store
	echo zit/.git/
	echo zit/\.zit/
	echo build/
	echo zit/zit$$
