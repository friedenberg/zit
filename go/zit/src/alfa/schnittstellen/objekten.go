package schnittstellen

type FileExtensionGetter interface {
	GetFileExtensionForGattung(GattungGetter) string
	GetFileExtensionZettel() string
	GetFileExtensionOrganize() string
	GetFileExtensionTyp() string
	GetFileExtensionEtikett() string
	GetFileExtensionKasten() string
}

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
	ObjekteReader(ShaGetter) (ShaReadCloser, error)
}

type ObjekteWriterFactory interface {
	ObjekteWriter() (ShaWriteCloser, error)
}

type (
	FuncObjekteReader func(ShaGetter) (ShaReadCloser, error)
	FuncObjekteWriter func() (ShaWriteCloser, error)
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
	sh ShaGetter,
) (ShaReadCloser, error) {
	return b.FuncObjekteReader(sh)
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

func (b bespokeObjekteWriteFactory) ObjekteWriter() (ShaWriteCloser, error) {
	return b.FuncObjekteWriter()
}
