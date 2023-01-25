package schnittstellen

type LockSmith interface {
	IsAcquired() bool
}
