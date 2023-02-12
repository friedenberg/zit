package schnittstellen

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
	ObjekteReader(GattungGetter, ShaGetter) (ShaReadCloser, error)
}

type ObjekteWriterFactory interface {
	ObjekteWriter(GattungGetter) (ShaWriteCloser, error)
}

type (
	FuncObjekteReader func(GattungGetter, ShaGetter) (ShaReadCloser, error)
	FuncObjekteWriter func(GattungGetter) (ShaWriteCloser, error)
)

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
	g GattungGetter,
	sh ShaGetter,
) (ShaReadCloser, error) {
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
	g GattungGetter,
) (ShaWriteCloser, error) {
	return b.FuncObjekteWriter(g)
}
