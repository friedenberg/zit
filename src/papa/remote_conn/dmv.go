package remote_conn

var chConnectionLimit chan struct{}

func init() {
	chConnectionLimit = make(chan struct{}, 100)
}

func AcquireConnLicense() {
	chConnectionLimit <- struct{}{}
}

func ReleaseConnLicense() {
	<-chConnectionLimit
}
