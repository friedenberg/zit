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
	assert_output - <<-EOM
		[-zz@5dbb297b5bde513be49fde397499eb89af8f5295f5137d75b52b015802b73ae0]
		[-zz-inbox@5dbb297b5bde513be49fde397499eb89af8f5295f5137d75b52b015802b73ae0]
		[-zz-inbox-2022@5dbb297b5bde513be49fde397499eb89af8f5295f5137d75b52b015802b73ae0]
		[-zz-inbox-2022-11@5dbb297b5bde513be49fde397499eb89af8f5295f5137d75b52b015802b73ae0]
		[-zz-inbox-2022-11-14@5dbb297b5bde513be49fde397499eb89af8f5295f5137d75b52b015802b73ae0]
		[one/uno@8f8aa93ce3cb3da0e5eddb2c9556fe37980d0aaf58f2760de451a93ce337b0c2 !md "to_add"]
		[one/uno@8f8aa93ce3cb3da0e5eddb2c9556fe37980d0aaf58f2760de451a93ce337b0c2 !md "to_add"]
		[to_add.md] (deleted)
	EOM

	run_zit show one/uno
	assert_success
	assert_output - <<-EOM
		---
		# to_add
		- zz-inbox-2022-11-14
		! md
		---

		test file
	EOM
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
	assert_output - <<-EOM
		[-zz@5dbb297b5bde513be49fde397499eb89af8f5295f5137d75b52b015802b73ae0]
		[-zz-inbox@5dbb297b5bde513be49fde397499eb89af8f5295f5137d75b52b015802b73ae0]
		[-zz-inbox-2022@5dbb297b5bde513be49fde397499eb89af8f5295f5137d75b52b015802b73ae0]
		[-zz-inbox-2022-11@5dbb297b5bde513be49fde397499eb89af8f5295f5137d75b52b015802b73ae0]
		[-zz-inbox-2022-11-14@5dbb297b5bde513be49fde397499eb89af8f5295f5137d75b52b015802b73ae0]
		[one/uno@8f8aa93ce3cb3da0e5eddb2c9556fe37980d0aaf58f2760de451a93ce337b0c2 !md "to_add"]
		[one/uno@8f8aa93ce3cb3da0e5eddb2c9556fe37980d0aaf58f2760de451a93ce337b0c2 !md "to_add"]
		[to_add.md] (deleted)
	EOM
}

function add_2 { # @test
	wd="$(mktemp -d)"
	cd "$wd" || exit 1

	run_zit_init_disable_age
	assert_success

	f=to_add.md
	{
		echo test file
	} >"$f"

	f2=to_add2.md
	{
		echo test file 2
	} >"$f2"

	run_zit add \
		-dedupe \
		-delete \
		-etiketten zz-inbox-2022-11-14 \
		"$f" "$f2"

	assert_success
	assert_output_unsorted - <<-EOM
		[-zz-inbox-2022-11-14@5dbb297b5bde513be49fde397499eb89af8f5295f5137d75b52b015802b73ae0]
		[-zz-inbox-2022-11-14@5dbb297b5bde513be49fde397499eb89af8f5295f5137d75b52b015802b73ae0]
		[-zz-inbox-2022-11@5dbb297b5bde513be49fde397499eb89af8f5295f5137d75b52b015802b73ae0]
		[-zz-inbox-2022-11@5dbb297b5bde513be49fde397499eb89af8f5295f5137d75b52b015802b73ae0]
		[-zz-inbox-2022@5dbb297b5bde513be49fde397499eb89af8f5295f5137d75b52b015802b73ae0]
		[-zz-inbox-2022@5dbb297b5bde513be49fde397499eb89af8f5295f5137d75b52b015802b73ae0]
		[-zz-inbox@5dbb297b5bde513be49fde397499eb89af8f5295f5137d75b52b015802b73ae0]
		[-zz-inbox@5dbb297b5bde513be49fde397499eb89af8f5295f5137d75b52b015802b73ae0]
		[-zz@5dbb297b5bde513be49fde397499eb89af8f5295f5137d75b52b015802b73ae0]
		[-zz@5dbb297b5bde513be49fde397499eb89af8f5295f5137d75b52b015802b73ae0]
		[one/uno@8f8aa93ce3cb3da0e5eddb2c9556fe37980d0aaf58f2760de451a93ce337b0c2 !md "to_add"]
		[one/uno@8f8aa93ce3cb3da0e5eddb2c9556fe37980d0aaf58f2760de451a93ce337b0c2 !md "to_add"]
		[one/dos@02425f5295479fc80efd565abe728696072de2422958209ef32ffb39427d80a1 !md "to_add2"]
		[one/dos@02425f5295479fc80efd565abe728696072de2422958209ef32ffb39427d80a1 !md "to_add2"]
		[to_add.md] (deleted)
		[to_add2.md] (deleted)
	EOM
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
