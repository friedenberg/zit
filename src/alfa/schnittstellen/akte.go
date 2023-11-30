package schnittstellen

type (
	Akte[T any] interface {
		GattungGetter
		Equatable[T]
	}

	AktePtr[T any] interface {
		Akte[T]
		Resetable[T]
		Ptr[T]
	}
)

//      _    _    _       ___ ___
//     / \  | | _| |_ ___|_ _/ _ \
//    / _ \ | |/ / __/ _ \| | | | |
//   / ___ \|   <| ||  __/| | |_| |
//  /_/   \_\_|\_\\__\___|___\___/
//

type AkteGetter[
	V any,
] interface {
	GetAkte(ShaLike) (V, error)
}

type AktePutter[
	V any,
] interface {
	PutAkte(V)
}

type AkteGetterPutter[
	V any,
] interface {
	AkteGetter[V]
	AktePutter[V]
}

type AkteIOFactory interface {
	AkteReaderFactory
	AkteWriterFactory
}

type AkteReaderFactory interface {
	AkteReader(ShaGetter) (ShaReadCloser, error)
}

type AkteWriterFactory interface {
	AkteWriter() (ShaWriteCloser, error)
}

type (
	FuncAkteReader func(ShaGetter) (ShaReadCloser, error)
	FuncAkteWriter func() (ShaWriteCloser, error)
)

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
	sh ShaGetter,
) (ShaReadCloser, error) {
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

func (b bespokeAkteWriteFactory) AkteWriter() (ShaWriteCloser, error) {
	return b.FuncAkteWriter()
}
