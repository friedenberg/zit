#! /usr/bin/env bats

setup() {
	load "$BATS_CWD/zz-tests_bats/common.bash"

	# for shellcheck SC2154
	export output
}

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

	run_zit_init_disable_age
	run_zit format-organize -prefix-joints=true -refine=true <(cat_organize)
	assert_output "$(cat_organize)"
}
