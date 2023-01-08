package fiddleware

// We might need to use etcd to coordinate this across multiple machines.
// TODO: look into using logSettings as sync.Map{} that allows us to turn parts of the service on/off without restarting
type Settings struct {
	Key           string // the path of the handler
	EnableLogging bool   // if we are enabled logging/recording
	LogLevel      string
	DumpResponse  bool
	DumpRequest   bool
	Record        bool
}
