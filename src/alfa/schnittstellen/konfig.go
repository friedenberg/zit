package schnittstellen

type Konfig interface {
	UsePredictableHinweisen() bool
	UseNewHinweisIndex() bool
}
