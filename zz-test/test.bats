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

function can_run_zit { # @test
	run zit
}

function provides_help_with_no_params { # @test
	run zit
	assert_output --partial 'No subcommand provided.'
}

function can_initialize_without_age { # @test
	wd="$(mktemp -d)"
	cd "$wd" || exit 1

	run zit init -disable-age -yin <(cat_yin) -yang <(cat_yang)
	[ -d .zit/ ]
	[ ! -f .zit/AgeIdentity ]
}

function can_initialize_with_age { # @test
	wd="$(mktemp -d)"
	cd "$wd" || exit 1

	run zit init -yin <(cat_yin) -yang <(cat_yang)
	[ -d .zit/ ]
	[ -f .zit/AgeIdentity ]
}

function can_new_zettel_file { # @test
	wd="$(mktemp -d)"
	cd "$wd" || exit 1

	run zit init -disable-age -yin <(cat_yin) -yang <(cat_yang)

	to_add="$(mktemp)"
	{
		echo "---"
		echo "# wow"
		echo "- ok"
		echo "---"
	} >>"$to_add"

	run zit new "$to_add"
	assert_output --partial '[one/uno '

	run zit show one/uno
	assert_output "$(cat "$to_add")"
}

function can_new_zettel { # @test
	wd="$(mktemp -d)"
	cd "$wd" || exit 1

	run zit init -disable-age -yin <(cat_yin) -yang <(cat_yang)

	expected="$(mktemp)"
	{
		echo "---"
		echo "# wow"
		echo "- ok"
    echo "! md"
		echo "---"
	} >"$expected"

	run zit new -bezeichnung wow -etiketten ok
	assert_output --partial '[one/uno '

	run zit show one/uno
	assert_output "$(cat "$expected")"
}

function can_checkout_and_checkin { # @test
	wd="$(mktemp -d)"
	cd "$wd" || exit 1

	run zit init -disable-age -yin <(cat_yin) -yang <(cat_yang)

	to_add="$(mktemp)"
	{
		echo "---"
		echo "# wow"
		echo "- ok"
		echo "---"
	} >>"$to_add"

	run zit new "$to_add"
	assert_output --partial '[one/uno '

	run zit checkout one/uno
	assert_output --partial '[one/uno '
	assert_output --partial '(checked out)'

	{
		echo "---"
		echo "# wow"
		echo "- ok"
		echo "---"
		echo
		echo "content"
	} >"one/uno.md"

	run zit checkin one/uno.md
	assert_output --partial '[one/uno '
	assert_output --partial '(updated)'
}

function can_checkout_via_etiketten { # @test
	skip
	wd="$(mktemp -d)"
	cd "$wd" || exit 1

	run zit init -disable-age -yin <(cat_yin) -yang <(cat_yang)

	to_add="$(mktemp)"
	{
		echo "---"
		echo "# wow"
		echo "- ok"
		echo "---"
	} >>"$to_add"

	run zit new "$to_add"
	assert_output --partial '[one/uno '

	run zit checkout -etiketten ok
	assert_output --partial '[one/uno '
	assert_output --partial '(checked out)'
}

function can_output_organize { # @test
	skip
	wd="$(mktemp -d)"
	cd "$wd" || exit 1

	run zit init -disable-age -yin <(cat_yin) -yang <(cat_yang)

	to_add="$(mktemp)"
	{
		echo "---"
		echo "# wow"
		echo "- ok"
		echo "---"
	} >>"$to_add"

	run zit new "$to_add"
	assert_output --partial '[one/uno '

	expected_organize="$(mktemp)"
	{
		echo
		echo "# ok"
		echo
		echo "- [one/uno] wow"
	} >>"$expected_organize"

	run zit organize ok <"$(tty)"
	assert_output "$(cat "$expected_organize")"

	{
		echo "# wow"
		echo
		echo
		echo "- [one/uno] wow"
	} >"$expected_organize"

	run zit organize ok <"$expected_organize"

	expected_zettel="$(mktemp)"
	{
		echo "---"
		echo "# wow"
		echo "- wow"
		echo "---"
	} >>"$expected_zettel"

	run zit show one/uno
	assert_output "$(cat "$expected_zettel")"
}

function hides_hidden_etiketten_from_organize { # @test
	skip
	wd="$(mktemp -d)"
	cd "$wd" || exit 1

	run zit init -disable-age -yin <(cat_yin) -yang <(cat_yang)

	{
		echo "[tags.zz-archive]"
		echo "hide = true"
	} >>.zit/Konfig

	to_add="$(mktemp)"
	{
		echo ---
		echo "# split hinweis for usability"
		echo - project-2021-zit
		echo - zz-archive-task-done
		echo ! md
		echo ---
	} >>"$to_add"

	run zit new "$to_add"
	assert_output --partial '[one/uno '

	expected_organize="$(mktemp)"
	{
		echo
		echo "# project-2021-zit"
		echo
	} >>"$expected_organize"

	run zit organize project-2021-zit
	assert_output "$(cat "$expected_organize")"
}

function can_new_zettel_with_metadatei { # @test
	wd="$(mktemp -d)"
	cd "$wd" || exit 1

	run zit init -disable-age -yin <(cat_yin) -yang <(cat_yang)

	expected="$(mktemp)"
	{
		echo ---
		echo "# bez"
		echo - et1
		echo - et2
		echo ! md
		echo ---
	} >>"$expected"

	run zit new -bezeichnung bez -etiketten et1,et2
	assert_output --partial '[one/uno '
	assert_output --partial '(checked out)'

	run cat one/uno.md
	assert_output "$(cat "$expected")"
}

function can_update_akte { # @test
	# setup
	wd="$(mktemp -d)"
	cd "$wd" || exit 1

	run zit init -disable-age -yin <(cat_yin) -yang <(cat_yang)

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
	} >>"$expected"

	run zit new "$expected"
	assert_output --partial '[one/uno '

	run zit show one/uno
	assert_output "$(cat "$expected")"

	# when
	new_akte="$(mktemp)"
	{
		echo the body but new
	} >>"$new_akte"

	run zit checkin-akte -verbose -new-etiketten et3 one/uno "$new_akte"
	assert_output --partial '[one/uno '
	assert_output --partial '(akte updated)'

	# then
	{
		echo ---
		echo "# bez"
		echo - et3
		echo ! md
		echo ---
		echo
		echo the body but new
	} >"$expected"

	run zit show one/uno
	assert_output "$(cat "$expected")"
}

function can_duplicate_zettel_content { # @test
	skip                                   #TODO:

	# setup
	wd="$(mktemp -d)"
	cd "$wd" || exit 1

	run zit init -disable-age -yin <(cat_yin) -yang <(cat_yang)

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
	} >>"$expected"

	run zit new "$expected"
	assert_output --partial '[one/uno '

	run zit new "$expected"
	assert_output --partial '[two/dos '

	# when
	run zit show one/uno
	assert_output --partial "$(cat "$expected")"
	run zit show two/dos
	assert_output --partial "$(cat "$expected")"
}

function indexes_are_implicitly_correct { # @test
	# setup
	wd="$(mktemp -d)"
	cd "$wd" || exit 1

	run zit init -disable-age -yin <(cat_yin) -yang <(cat_yang)

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
	} >>"$expected"

	run zit new "$expected"
	assert_output --partial '[one/uno '

	{
		echo et1
		echo et2
	} >"$expected"

	run zit cat -type etikett
	assert_output "$(cat "$expected")"

	{
		echo one/uno
	} >"$expected"

	#TODO
	# run zit cat -type hinweis
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
	cp "$expected" "one/uno.md"
	run zit checkin -delete "one/uno.md"

	{
		echo et1
	} >"$expected"

	run zit cat -type etikett
	assert_output "$(cat "$expected")"

	{
		echo one/uno
	} >"$expected"

	#TODO
	# run zit cat -type hinweis
	# assert_output --partial "$(cat "$expected")"
}

function checkouts_dont_overwrite { # @test
	# setup
	wd="$(mktemp -d)"
	cd "$wd" || exit 1

	run zit init -disable-age -yin <(cat_yin) -yang <(cat_yang)

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
	} >>"$expected"

	run zit new "$expected"
	assert_output --partial '[one/uno '

	run zit checkout one/uno

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

	cat "$expected" >"one/uno.md"

	run zit checkout one/uno
	assert_output --partial '[one/uno] (external has changes)'

	run cat one/uno.md
	assert_output "$(cat "$expected")"
}
