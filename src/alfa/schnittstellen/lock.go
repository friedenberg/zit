package schnittstellen

type LockSmith interface {
	IsAcquired() bool
}

type LockSmithGetter interface {
	GetLockSmith() LockSmith
}
