package schnittstellen

type ObjekteAkteReaderFactory interface {
	ObjekteReaderFactory
	AkteReaderFactory
}

type ObjekteAkteWriterFactory interface {
	ObjekteWriterFactory
	AkteWriterFactory
}

type ObjekteAkteFactory interface {
	ObjekteAkteReaderFactory
	ObjekteAkteWriterFactory
}
