#! /usr/bin/env bats

setup() {
	load "$BATS_CWD/zz-tests_bats/common.bash"

	# for shellcheck SC2154
	export output
}

function add { # @test
	wd="$(mktemp -d)"
	cd "$wd" || exit 1

	run zit init -disable-age -yin <(cat_yin) -yang <(cat_yang)
	assert_success

	f=to_add.md
	{
		echo test file
	} >"$f"

	run_zit add \
		-dedupe \
		-delete \
		-etiketten zz-inbox-2022-11-14 \
		"$f"

	assert_success
	assert_output --partial '[one/uno@45e1f5fbfe972f697f4f4f4b77a21f6395e4cf3a1f0ca16d34a675e447ab3778 !md "to_add.md"]'
	assert_output --partial '[one/uno@11327fbe60cabd2a9eabf4a37d541cf04b539f913945897efe9bab1e30784781 !md "to_add.md"]'
	assert_output --partial '[to_add.md] (deleted)'
}

function add_1 { # @test
	wd="$(mktemp -d)"
	cd "$wd" || exit 1

	run zit init -disable-age -yin <(cat_yin) -yang <(cat_yang)
	assert_success

	f=to_add.md
	{
		echo test file
	} >"$f"

	run_zit add \
		-dedupe \
		-delete \
		-etiketten zz-inbox-2022-11-14 \
		"$f"

	assert_success
	assert_output --partial '[one/uno@45e1f5fbfe972f697f4f4f4b77a21f6395e4cf3a1f0ca16d34a675e447ab3778 !md "to_add.md"]'
	assert_output --partial '[one/uno@11327fbe60cabd2a9eabf4a37d541cf04b539f913945897efe9bab1e30784781 !md "to_add.md"]'
	assert_output --partial '[to_add.md] (deleted)'
}

function add_dedupe_1 { # @test
	skip
	wd="$(mktemp -d)"
	cd "$wd" || exit 1

	run zit init -disable-age -yin <(cat_yin) -yang <(cat_yang)
	assert_success

	f=to_add.md
	{
		echo test file
	} >"$f"

	run run_zit add \
		-dedupe \
		-etiketten zz-inbox-2022-11-14 \
		"$f"

	assert_success
	assert_output --partial '          (new) [o/u@b !md "to_add.md"]'
	assert_output --partial '      (updated) [o/u@d !md "to_add.md"]'
	# assert_output --partial '[to_add.md] (deleted)'

	run zit checkout o/u
	#TODO-P2 fix race condition
	assert_success
	assert_output '  (checked out) [one/uno.zettel@d !md "to_add.md"]'

	{
		echo '---'
		echo '# new title'
		echo '- new-tag'
		echo '! md'
		echo '---'
		echo ''
		echo 'test file'
	} >one/uno.zettel

	run zit checkin -delete one/uno.zettel
	assert_success
	assert_output --partial '      (updated) [o/u@a !md "new title"]'
	assert_output --partial '      (deleted) [one/uno.zettel]'

	run zit add \
		-predictable-hinweisen \
		-dedupe \
		-delete \
		-etiketten new-etikett-2 \
		"$f"

	{
		echo '---'
		echo '# new title'
		echo '- new-etikett-2'
		echo '- new-tag'
		echo '! md'
		echo '---'
		echo ''
		echo 'test file'
	} >expected

	run zit show o/u
	#TODO-P2 fix race condition
	assert_success
	assert_output "$(cat expected)"
}
