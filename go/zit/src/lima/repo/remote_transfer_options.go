package repo

import "flag"

// TODO add HTTP header options for these flags
type RemoteTransferOptions struct {
	PrintCopies         bool
	IncludeObjects      bool
	IncludeBlobs        bool
	AllowMergeConflicts bool
}

func (options *RemoteTransferOptions) SetFlagSet(f *flag.FlagSet) {
	f.BoolVar(
		&options.IncludeObjects,
		"include-objects",
		true,
		"imports the object during transfer",
	)

	f.BoolVar(
		&options.IncludeBlobs,
		"include-blobs",
		true,
		"copy the blob when performing the object transfer",
	)

	f.BoolVar(
		&options.AllowMergeConflicts,
		"allow-merge-conflicts",
		false,
		"ignore merge conflicts and allow incompatible histories to coexist",
	)
}

func (options RemoteTransferOptions) WithPrintCopies(
	value bool,
) RemoteTransferOptions {
	options.PrintCopies = value
	return options
}
