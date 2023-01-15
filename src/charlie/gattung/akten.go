package gattung

import "github.com/friedenberg/zit/src/bravo/sha"

//      _    _    _       ___ ___
//     / \  | | _| |_ ___|_ _/ _ \
//    / _ \ | |/ / __/ _ \| | | | |
//   / ___ \|   <| ||  __/| | |_| |
//  /_/   \_\_|\_\\__\___|___\___/
//

type AkteIOFactory interface {
	AkteReaderFactory
	AkteWriterFactory
}

type AkteReaderFactory interface {
	AkteReader(ShaLike) (sha.ReadCloser, error)
}

type AkteWriterFactory interface {
	AkteWriter() (sha.WriteCloser, error)
}

type ObjekteAkteReaderFactory interface {
	ObjekteReaderFactory
	AkteReaderFactory
}

type ObjekteAkteWriterFactory interface {
	ObjekteWriterFactory
	AkteWriterFactory
}

type FuncAkteReader func(ShaLike) (sha.ReadCloser, error)
type FuncAkteWriter func() (sha.WriteCloser, error)

type bespokeAkteReadWriterFactory struct {
	AkteReaderFactory
	AkteWriterFactory
}

func MakeBespokeAkteReadWriterFactory(
	r AkteReaderFactory,
	w AkteWriterFactory,
) AkteIOFactory {
	return bespokeAkteReadWriterFactory{
		AkteReaderFactory: r,
		AkteWriterFactory: w,
	}
}

type bespokeAkteReadFactory struct {
	FuncAkteReader
}

func MakeBespokeAkteReadFactory(
	r FuncAkteReader,
) AkteReaderFactory {
	return bespokeAkteReadFactory{
		FuncAkteReader: r,
	}
}

func (b bespokeAkteReadFactory) AkteReader(
	sh ShaLike,
) (sha.ReadCloser, error) {
	return b.FuncAkteReader(sh)
}

type bespokeAkteWriteFactory struct {
	FuncAkteWriter
}

func MakeBespokeAkteWriteFactory(
	r FuncAkteWriter,
) AkteWriterFactory {
	return bespokeAkteWriteFactory{
		FuncAkteWriter: r,
	}
}

func (b bespokeAkteWriteFactory) AkteWriter() (sha.WriteCloser, error) {
	return b.FuncAkteWriter()
}
