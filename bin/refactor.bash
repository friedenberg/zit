#! /bin/bash -e

old="$1"
new="$2"
package="$3"

# file_diff="$(mktemp)"

(
	go_refactor_args=(
		"-r"
		"$package.$old -> $package.$new"
	)

	gofmt "${go_refactor_args[@]}" src/
)

(
	go_refactor_args=(
		"-r"
		"$old -> $package"
	)

	gofmt "${go_refactor_args[@]}" "src/$package"
)
