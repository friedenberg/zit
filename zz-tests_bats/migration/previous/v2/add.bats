#! /usr/bin/env bats

setup() {
	load "$(dirname "$BATS_TEST_FILE")/common.bash"
}

teardown() {
	chflags_and_rm
}

function add { # @test
	run_zit_init_disable_age

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
		[-zz@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[-zz-inbox@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[-zz-inbox-2022@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[-zz-inbox-2022-11@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[-zz-inbox-2022-11-14@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[one/uno@55f8718109829bf506b09d8af615b9f107a266e19f7a311039d1035f180b22d4 !md "to_add"]
		[one/uno@55f8718109829bf506b09d8af615b9f107a266e19f7a311039d1035f180b22d4 !md "to_add"]
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
	run_zit_init_disable_age

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
		[-zz@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[-zz-inbox@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[-zz-inbox-2022@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[-zz-inbox-2022-11@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[-zz-inbox-2022-11-14@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[one/uno@55f8718109829bf506b09d8af615b9f107a266e19f7a311039d1035f180b22d4 !md "to_add"]
		[one/uno@55f8718109829bf506b09d8af615b9f107a266e19f7a311039d1035f180b22d4 !md "to_add"]
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
		[-zz-inbox-2022-11-14@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[-zz-inbox-2022-11@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[-zz-inbox-2022@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[-zz-inbox@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[-zz@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[one/dos@6b8e3c36cb01aa01c65ccd86e04935695e9da0580e3dcc8c0c3bce146c274c2c !md "to_add2"]
		[one/dos@6b8e3c36cb01aa01c65ccd86e04935695e9da0580e3dcc8c0c3bce146c274c2c !md "to_add2"]
		[one/uno@55f8718109829bf506b09d8af615b9f107a266e19f7a311039d1035f180b22d4 !md "to_add"]
		[one/uno@55f8718109829bf506b09d8af615b9f107a266e19f7a311039d1035f180b22d4 !md "to_add"]
		[to_add.md] (deleted)
		[to_add2.md] (deleted)
	EOM
}

#function add_dedupe_1 { ## @test
#	wd="$(mktemp -d)"
#	cd "$wd" || exit 1

#	run_zit init -disable-age -yin <(cat_yin) -yang <(cat_yang)
#	assert_success

#	f=to_add.md
#	{
#		echo test file
#	} >"$f"

#	run_zit add \
#		-dedupe \
#		-etiketten zz-inbox-2022-11-14 \
#		"$f"

#	assert_success
#	assert_output - <<-EOM
#		[-zz@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249]
#		[-zz-inbox@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249]
#		[-zz-inbox-2022@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249]
#		[-zz-inbox-2022-11@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249]
#		[-zz-inbox-2022-11-14@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249]
#		[one/uno@8f8aa93ce3cb3da0e5eddb2c9556fe37980d0aaf58f2760de451a93ce337b0c2 !md "to_add"]
#		[one/uno@8f8aa93ce3cb3da0e5eddb2c9556fe37980d0aaf58f2760de451a93ce337b0c2 !md "to_add"]
#	EOM

#	run_zit checkout o/u
#	#TODO-P2 fix race condition
#	assert_success
#	assert_output - <<-EOM
#		      checked out [one/uno.zettel@8f8aa93ce3cb3da0e5eddb2c9556fe37980d0aaf58f2760de451a93ce337b0c2 !md "to_add"]
#	EOM

#	{
#		echo '---'
#		echo '# new title'
#		echo '- new-tag'
#		echo '! md'
#		echo '---'
#		echo ''
#		echo 'test file'
#	} >one/uno.zettel

#	run_zit checkin -delete one/uno.zettel
#	assert_success
#	assert_output - <<-EOM
#		[-new@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249]
#		[-new-tag@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249]
#		[one/uno@d4853a453015235e41b9513f7e70d91b1a28212f9bd342daf5024b84f35d209f !md "new title"]
#		          deleted [one/uno.zettel]
#		          deleted [one]
#	EOM

#	run zit add \
#		-predictable-hinweisen \
#		-dedupe \
#		-delete \
#		-etiketten new-etikett-2 \
#		"$f"

#	run zit show o/u
#	#TODO-P2 fix race condition
#	assert_success
#	assert_output - <<-EOM
#		---
#		# new title
#		- new-etikett-2
#		- new-tag
#		! md
#		---

#		test file
#	EOM
#}