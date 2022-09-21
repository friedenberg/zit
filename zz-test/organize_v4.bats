#! /usr/bin/env bats

setup() {
	load 'test_helper/bats-support/load'
	load 'test_helper/bats-assert/load'
	# ... the remaining setup is unchanged

	# get the containing directory of this file
	# use $BATS_TEST_FILENAME instead of ${BASH_SOURCE[0]} or $0,
	# as those will point to the bats executable's location or the preprocessed file respectively
	DIR="$(cd "$(dirname "$BATS_TEST_FILENAME")" >/dev/null 2>&1 && pwd)"
	# make executables in src/ visible to PATH
	PATH="$DIR/../:$PATH"
	PATH="$DIR/../build/:$PATH"

	# for shellcheck SC2154
	export output
}

cat_yin() (
	echo "one"
	echo "two"
	echo "three"
)

cat_yang() (
	echo "uno"
	echo "dos"
	echo "tres"
)

cat_organize() (
	echo ''
	echo '- [ ach/vil    ] blah'
	echo ''
	echo '               # project'
	echo ''
	echo '              ##        -2021-zit'
	echo ''
	echo '             ###                 -22q1-uws-140'
	echo ''
	echo '            ####                              -moving'
	echo ''
	echo '- [  io/poliwr ] update billing addresses'
	echo ''
	echo '            ####                              -mvp-main_room'
	echo ''
	echo '- [prot/nidora ] Brainstorm where to place toolbox.md'
	echo ''
	echo '             ###                 -commands'
	echo ''
	echo '- [ tec/slowp  ] update output of commands to use new store'
	echo '- [ mer/golb   ] use error types to generate specific exit status codes'
	echo ''
	echo '             ###                 -etiketten_and_organize'
	echo ''
	echo '- [  pe/mo     ] add etikett rule type for removing etiketts based on conditions'
	echo '- [yttr/gole   ] use default etiketten with add'
	echo ''
	echo '             ###                 -init'
	echo ''
	echo '- [  ph/hitmonc] Add bats test for initing more than once.md'
	echo '- [ rub/rap    ] add .exrc to init'
	echo ''
)

function outputs_organize_one_etikett { # @test
	# skip
	wd="$(mktemp -d)"
	cd "$wd" || exit 1

	run zit init -disable-age -yin <(cat_yin) -yang <(cat_yang)
	run zit format-organize -prefix-joints=true -refine=true <(cat_organize)
	assert_output "$(cat_organize)"
}
