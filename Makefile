
.PHONY: build watch exclude;

build:
	go build -o build/zit ./.
	go test ./...
	go install ./.
	# bats zz-test/test.bats

watch:
	echo .

exclude:
	echo .DS_Store
	echo zit/.git/
	echo zit/\.zit/
	echo build/
	echo zit/zit$$
