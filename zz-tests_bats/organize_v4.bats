#! /usr/bin/env bats

setup() {
	load "$(dirname "$BATS_TEST_FILE")/common.bash"

	# for shellcheck SC2154
	export output
}

teardown() {
	rm_from_version
}

cmd_def_organize=(
	-prefix-joints=true
	-refine=true
	-new-organize=false
)

cat_organize() (
	cat - <<-EOM

		- [ ach/vil    ] blah

		               # project

		              ##        -2021-zit

		             ###                 -22q1-uws-140

		            ####                              -moving

		- [  io/poliwr ] update billing addresses

		            ####                              -mvp-main_room

		- [prot/nidora ] Brainstorm where to place toolbox.md

		             ###                 -commands

		- [ mer/golb   ] use error types to generate specific exit status codes
		- [ tec/slowp  ] update output of commands to use new store

		             ###                 -etiketten_and_organize

		- [  pe/mo     ] add etikett rule type for removing etiketts based on conditions
		- [yttr/gole   ] use default etiketten with add

		             ###                 -init

		- [  ph/hitmonc] Add bats test for initing more than once.md
		- [ rub/rap    ] add .exrc to init
	EOM
)

function outputs_organize_one_etikett { # @test
	cd "$BATS_TEST_TMPDIR" || exit 1
	run_zit_init_disable_age

	run_zit format-organize "${cmd_def_organize[@]}" <(cat_organize)
	assert_output "$(cat_organize)"
}
