#! /usr/bin/env bats

setup() {
	load 'test_helper/bats-support/load'
	load 'test_helper/bats-assert/load'
	load 'common.bash'
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

function provides_help_with_no_params { # @test
	run zit
	assert_output --partial 'No subcommand provided.'
}

function can_initialize_without_age { # @test
	wd="$(mktemp -d)"
	cd "$wd" || exit 1

	run_zit_init_disable_age

	[ -d .zit/ ]
	[ ! -f .zit/AgeIdentity ]
}

function can_initialize_with_age { # @test
	wd="$(mktemp -d)"
	cd "$wd" || exit 1

	run zit init -yin <(cat_yin) -yang <(cat_yang)
	assert_success

	[ -d .zit/ ]
	[ -f .zit/AgeIdentity ]
}

function can_new_zettel_file { # @test
	wd="$(mktemp -d)"
	cd "$wd" || exit 1

	run_zit_init_disable_age

	to_add="$(mktemp)"
	{
		echo "---"
		echo "# wow"
		echo "- ok"
		echo "! md"
		echo "---"
	} >"$to_add"

	run_zit new -edit=false "$to_add"
	assert_output '[one/uno@9a638e2b183562da6d3c634d5a3841d64bc337c9cf79f8fffa0d0194659bc564 !md "wow"]'

	run_zit show one/uno
	assert_output "$(cat "$to_add")"
}

function can_new_zettel { # @test
	wd="$(mktemp -d)"
	cd "$wd" || exit 1

	run_zit_init_disable_age

	expected="$(mktemp)"
	{
		echo "---"
		echo "# wow"
		echo "- ok"
		echo "! md"
		echo "---"
	} >"$expected"

	run_zit new -edit=false -bezeichnung wow -etiketten ok
	assert_output '[one/uno@9a638e2b183562da6d3c634d5a3841d64bc337c9cf79f8fffa0d0194659bc564 !md "wow"]'

	run_zit show one/uno
	assert_output "$(cat "$expected")"
}

function can_checkout_and_checkin { # @test
	wd="$(mktemp -d)"
	cd "$wd" || exit 1

	run_zit_init_disable_age

	to_add="$(mktemp)"
	{
		echo "---"
		echo "# wow"
		echo "- ok"
		echo "---"
	} >"$to_add"

	run_zit new -edit=false "$to_add"
	assert_output '[one/uno@9a638e2b183562da6d3c634d5a3841d64bc337c9cf79f8fffa0d0194659bc564 !md "wow"]'

	run_zit checkout one/uno
	assert_output '  (checked out) [one/uno.zettel@9a638e2b183562da6d3c634d5a3841d64bc337c9cf79f8fffa0d0194659bc564 !md "wow"]'

	{
		echo "---"
		echo "# wow"
		echo "- ok"
		echo "---"
		echo
		echo "content"
	} >"one/uno.zettel"

	run_zit checkin one/uno.zettel
	#TODO fix missing typ
	assert_output '[one/uno@14d2d788146303057462fbf3d181a3c8c3397ebc238c07970b206b5db6203a3a ! "wow"]'
}

function can_checkout_via_etiketten { # @test
	skip
	wd="$(mktemp -d)"
	cd "$wd" || exit 1

	run_zit_init_disable_age

	to_add="$(mktemp)"
	{
		echo "---"
		echo "# wow"
		echo "- ok"
		echo "---"
	} >"$to_add"

	run_zit new -edit=false "$to_add"
	assert_output --partial '[one/uno '

	run_zit checkout -etiketten ok
	assert_output --partial '[one/uno '
	assert_output --partial '(checked out)'
}

function can_output_organize { # @test
	skip
	wd="$(mktemp -d)"
	cd "$wd" || exit 1

	run_zit_init_disable_age

	to_add="$(mktemp)"
	{
		echo "---"
		echo "# wow"
		echo "- ok"
		echo "---"
	} >"$to_add"

	run_zit new -edit=false "$to_add"
	assert_output --partial '[one/uno '

	expected_organize="$(mktemp)"
	{
		echo
		echo "# ok"
		echo
		echo "- [one/uno] wow"
	} >"$expected_organize"

	run_zit organize ok <"$(tty)"
	assert_output "$(cat "$expected_organize")"

	{
		echo "# wow"
		echo
		echo
		echo "- [one/uno] wow"
	} >"$expected_organize"

	run_zit organize ok <"$expected_organize"

	expected_zettel="$(mktemp)"
	{
		echo "---"
		echo "# wow"
		echo "- wow"
		echo "---"
	} >"$expected_zettel"

	run_zit show one/uno
	assert_output "$(cat "$expected_zettel")"
}

function hides_hidden_etiketten_from_organize { # @test
	skip
	wd="$(mktemp -d)"
	cd "$wd" || exit 1

	run_zit_init_disable_age

	{
		echo "[tags.zz-archive]"
		echo "hide = true"
	} >.zit/Konfig

	to_add="$(mktemp)"
	{
		echo ---
		echo "# split hinweis for usability"
		echo - project-2021-zit
		echo - zz-archive-task-done
		echo ! md
		echo ---
	} >"$to_add"

	run_zit new -edit=false "$to_add"
	assert_output --partial '[one/uno '

	expected_organize="$(mktemp)"
	{
		echo
		echo "# project-2021-zit"
		echo
	} >"$expected_organize"

	run_zit organize project-2021-zit
	assert_output "$(cat "$expected_organize")"
}

function can_new_zettel_with_metadatei { # @test
	wd="$(mktemp -d)"
	cd "$wd" || exit 1

	run_zit_init_disable_age

	expected="$(mktemp)"
	{
		echo ---
		echo "# bez"
		echo - et1
		echo - et2
		echo ! md
		echo ---
	} >"$expected"

	run_zit new -edit=false -bezeichnung bez -etiketten et1,et2
	assert_output '[one/uno@22bfa88b3975bc7cad702c2c7f262d5a754d9ad7423b96b134c6bbc1fbd88aea !md "bez"]'
}

function can_duplicate_zettel_content { # @test
	# TODO
	skip

	# setup
	wd="$(mktemp -d)"
	cd "$wd" || exit 1

	run_zit_init_disable_age

	expected="$(mktemp)"
	{
		echo ---
		echo "# bez"
		echo - et1
		echo - et2
		echo ! md
		echo ---
		echo
		echo the body
	} >"$expected"

	run_zit new -edit=false "$expected"
	assert_output ''

	run_zit new -edit=false "$expected"
	assert_output ''

	# when
	run_zit show one/uno
	assert_output ''

	run_zit show two/dos
	assert_output ''
}

function indexes_are_implicitly_correct { # @test
	# setup
	wd="$(mktemp -d)"
	cd "$wd" || exit 1

	run_zit_init_disable_age

	expected="$(mktemp)"
	{
		echo ---
		echo "# bez"
		echo - et1
		echo - et2
		echo ! md
		echo ---
		echo
		echo the body
	} >"$expected"

	run_zit new -edit=false "$expected"
	assert_output '[one/uno@18df16846a2f8bbce5f03e1041baff978a049aabd169ab9adac387867fe1706c !md "bez"]'

	{
		echo et1
		echo et2
	} >"$expected"

	run_zit cat -gattung etikett
	assert_output "$(cat "$expected")"

	{
		echo one/uno
	} >"$expected"

	#TODO
	# run_zit cat -gattung hinweis
	# assert_output --partial "$(cat "$expected")"

	{
		echo ---
		echo "# bez"
		echo - et1
		echo ! md
		echo ---
		echo
		echo the body
	} >"$expected"

	mkdir -p one
	cp "$expected" "one/uno.zettel"
	run_zit checkin -delete "one/uno.zettel"
	assert_output --partial '[one/uno@50bedb194bbd829d5d5d11de711a58b8486954a481ae43b4d1a8c4bd7f1f1370 !md "bez"]'
	assert_output --partial '      (deleted) [one/uno.zettel]'

	{
		echo et1
	} >"$expected"

	run_zit cat -gattung etikett
	assert_output "$(cat "$expected")"

	{
		echo one/uno
	} >"$expected"

	#TODO
	# run_zit cat -gattung hinweis
	# assert_output --partial "$(cat "$expected")"
}

function checkouts_dont_overwrite { # @test
	# setup
	wd="$(mktemp -d)"
	cd "$wd" || exit 1

	run_zit init -disable-age -yin <(cat_yin) -yang <(cat_yang)
	assert_success

	expected="$(mktemp)"
	{
		echo ---
		echo "# bez"
		echo - et1
		echo - et2
		echo ! md
		echo ---
		echo
		echo the body
	} >"$expected"

	run_zit new -edit=false "$expected"
	assert_output '[one/uno@18df16846a2f8bbce5f03e1041baff978a049aabd169ab9adac387867fe1706c !md "bez"]'

	run_zit checkout one/uno
	assert_output '  (checked out) [one/uno.zettel@18df16846a2f8bbce5f03e1041baff978a049aabd169ab9adac387867fe1706c !md "bez"]'

	run cat one/uno.zettel
	assert_output "$(cat "$expected")"

	{
		echo ---
		echo "# bez"
		echo - et1
		echo - et2
		echo ! md
		echo ---
		echo
		echo the body 2
	} >"$expected"

	cat "$expected" >"one/uno.zettel"

	run_zit checkout one/uno
	assert_output '  (checked out) [one/uno.zettel@63b65ad24c58d43d363f8074a5513e5cf71337cc132f452095a779b933cfee15 !md "bez"]'

	run cat one/uno.zettel
	assert_output "$(cat "$expected")"
}