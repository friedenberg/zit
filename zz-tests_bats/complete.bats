#! /usr/bin/env bats

setup() {
	load "$(dirname "$BATS_TEST_FILE")/common.bash"

	version="v$(zit info store-version)"
	copy_from_version "$DIR" "$version"
}

teardown() {
	rm_from_version "$version"
}

function assert_complete_fails() (
	for cmd in "$@"; do
		run_zit "$cmd" -complete
		assert_failure
		assert_output "command \"$cmd\" does not support completion"
	done
)

function complete_fails { # @test
	cmds=(
		checkin
		checkin-blob
		checkin-json
		checkout
	)

	assert_complete_fails "${cmds[@]}"
}

function complete_show { # @test
	run_zit show -complete
	assert_success
	assert_output_unsorted --regexp - <<-'EOM'
		tag-1.*Etikett
		tag-2.*Etikett
		tag-3.*Etikett
		tag-4.*Etikett
		tag.*Etikett
	EOM
}

function complete_show_all { # @test
	run_zit show -complete :z,t,b,e
	assert_success
	assert_output_unsorted --regexp - <<-EOM
		.*Bestandsaufnahme
		.*Bestandsaufnahme
		.*Bestandsaufnahme
		.*Bestandsaufnahme
		!md.*Typ
		one/dos.*Zettel: !md wow ok again
		one/uno.*Zettel: !md wow the first
		tag.*Etikett
		tag.1.*Etikett
		tag.2.*Etikett
		tag.3.*Etikett
		tag.4.*Etikett
	EOM
}

function complete_show_zettels { # @test
	run_zit show -verbose -complete :z
	assert_success
	assert_output_unsorted --regexp - <<-EOM
		one/dos.*Zettel: !md wow ok again
		one/uno.*Zettel: !md wow the first
	EOM
}

function complete_show_types { # @test
	run_zit show -complete :t
	assert_success
	assert_output_unsorted --regexp - <<-EOM
		md.*Typ
	EOM
}

function complete_show_tags { # @test
	run_zit show -complete :e
	assert_success
	assert_output_unsorted --regexp - <<-EOM
		tag-3.*Etikett
		tag-4.*Etikett
	EOM
}

function complete_subcmd { # @test
	run_zit complete
	assert_success
	assert_output_unsorted --regexp - <<-'EOM'
		add
		cat-alfred
		cat-blob
		cat-blob-shas
		checkin
		checkin-blob
		checkin-json
		checkout
		clean
		clone
		complete.*complete a command-line
		deinit
		diff
		dormant-add
		dormant-edit
		dormant-remove
		edit
		edit-config
		exec
		export
		find-missing
		format-blob
		format-object
		format-organize
		fsck
		import
		info
		info-repo
		info-workspace
		init
		init-archive
		init-workspace
		last
		merge-tool
		new
		organize
		peek-zettel-ids
		pull
		pull-blob-store
		push
		read-blob
		reindex
		revert
		save
		serve
		show
		status
		test
		write-blob
	EOM
}

function complete_complete { # @test
	run_zit complete complete
	assert_success
	assert_output_unsorted --regexp - <<-'EOM'
		-bash-style.*
		-in-progress.*
	EOM
}

function complete_init_workspace { # @test
	run_zit complete init-workspace
	assert_success
	assert_output_unsorted --regexp - <<-'EOM'
		-query.*default query for `show`
		-tags.*tags added for new objects in `checkin`, `new`, `organize`
		-type.*type used for new objects in `new` and `organize`
	EOM

	run_zit complete init-workspace -tags
	assert_success
	assert_output_unsorted --regexp - <<-'EOM'
		tag-1.*Etikett
		tag-2.*Etikett
		tag-3.*Etikett
		tag-4.*Etikett
		tag.*Etikett
	EOM

	run_zit complete -in-progress="tag" init-workspace -tags tag
	assert_success
	assert_output_unsorted --regexp - <<-'EOM'
		tag-1.*Etikett
		tag-2.*Etikett
		tag-3.*Etikett
		tag-4.*Etikett
		tag.*Etikett
	EOM

	mkdir -p workspaces/test

	run_zit complete -in-progress="workspaces" init-workspace -tags tag workspaces
	assert_success
	assert_output_unsorted --regexp - <<-'EOM'
		-query.*default query for `show`
		-tags.*tags added for new objects in `checkin`, `new`, `organize`
		test/.*directory
		-type.*type used for new objects in `new` and `organize`
	EOM
}
