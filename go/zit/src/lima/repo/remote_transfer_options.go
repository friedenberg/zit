package repo

import "flag"

type RemoteTransferOptions struct {
	PrintCopies         bool
	IncludeBlobs        bool
	AllowMergeConflicts bool
}

func (options *RemoteTransferOptions) SetFlagSet(f *flag.FlagSet) {
	f.BoolVar(&options.IncludeBlobs, "include-blobs", true, "copy the blob when performing the object transfer")
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
