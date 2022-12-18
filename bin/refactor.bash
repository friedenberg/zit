#! /bin/bash -xe

old="$1"
new="$2"
package="$3"
super_package="$4"

# file_diff="$(mktemp)"

(
	go_refactor_args=(
    "-w"
		"-r"
		"$package.$old -> $package.$new"
	)

	gofmt "${go_refactor_args[@]}" src/
)

(
	go_refactor_args=(
    "-w"
		"-r"
		"$old -> $package"
	)

	gofmt "${go_refactor_args[@]}" "src/$super_package/$package"
)
