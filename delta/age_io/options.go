package age_io

import "github.com/friedenberg/zit/charlie/age"

type CommonOptions struct {
	age.Age
	UseZip bool
}

type ReadOptions struct {
	CommonOptions
}

type WriteOptions struct {
	CommonOptions
}
