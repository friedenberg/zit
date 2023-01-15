package gattung

import "github.com/friedenberg/zit/src/bravo/sha"

//    ___  _     _      _    _       ___ ___
//   / _ \| |__ (_) ___| | _| |_ ___|_ _/ _ \
//  | | | | '_ \| |/ _ \ |/ / __/ _ \| | | | |
//  | |_| | |_) | |  __/   <| ||  __/| | |_| |
//   \___/|_.__// |\___|_|\_\\__\___|___\___/
//            |__/

type ObjekteIOFactory interface {
	ObjekteReaderFactory
	ObjekteWriterFactory
}

type ObjekteReaderFactory interface {
	ObjekteReader(GattungLike, ShaLike) (sha.ReadCloser, error)
}

type ObjekteWriterFactory interface {
	ObjekteWriter(GattungLike) (sha.WriteCloser, error)
}

type FuncObjekteReader func(GattungLike, ShaLike) (sha.ReadCloser, error)
type FuncObjekteWriter func(GattungLike) (sha.WriteCloser, error)

type bespokeObjekteReadWriterFactory struct {
	ObjekteReaderFactory
	ObjekteWriterFactory
}

func MakeBespokeObjekteReadWriterFactory(
	r ObjekteReaderFactory,
	w ObjekteWriterFactory,
) ObjekteIOFactory {
	return bespokeObjekteReadWriterFactory{
		ObjekteReaderFactory: r,
		ObjekteWriterFactory: w,
	}
}

type bespokeObjekteReadFactory struct {
	FuncObjekteReader
}

func MakeBespokeObjekteReadFactory(
	r FuncObjekteReader,
) ObjekteReaderFactory {
	return bespokeObjekteReadFactory{
		FuncObjekteReader: r,
	}
}

func (b bespokeObjekteReadFactory) ObjekteReader(
	g GattungLike,
	sh ShaLike,
) (sha.ReadCloser, error) {
	return b.FuncObjekteReader(g, sh)
}

type bespokeObjekteWriteFactory struct {
	FuncObjekteWriter
}

func MakeBespokeObjekteWriteFactory(
	r FuncObjekteWriter,
) ObjekteWriterFactory {
	return bespokeObjekteWriteFactory{
		FuncObjekteWriter: r,
	}
}

func (b bespokeObjekteWriteFactory) ObjekteWriter(
	g GattungLike,
) (sha.WriteCloser, error) {
	return b.FuncObjekteWriter(g)
}
