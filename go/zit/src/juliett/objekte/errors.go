package objekte

type VerlorenAndGefundenError interface {
	error
	AddToLostAndFound(string) (string, error)
}
