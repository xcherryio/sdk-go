package xdb

type ClientOptions struct {
	Namespace           string
	ServerUrl           string
	WorkerUrl           string
	ObjectEncoder       ObjectEncoder
	DBConverter         DBConverter
	EnabledDebugLogging bool
	// TODO API timeout and retry policy
}

const (
	DefaultNamespace  = "default"
	DefaultWorkerPort = "8803"
	DefaultServerPort = "8801"

	DefaultWorkerUrl = "http://localhost:" + DefaultWorkerPort
	DefaultServerUrl = "http://localhost:" + DefaultServerPort
)

var localDefaultClientOptions = ClientOptions{
	Namespace:     DefaultNamespace,
	ServerUrl:     DefaultServerUrl,
	WorkerUrl:     DefaultWorkerUrl,
	ObjectEncoder: GetDefaultObjectEncoder(),
	DBConverter:   NewBasicDBConverter(),
}

func GetLocalDefaultClientOptions() *ClientOptions {
	return &localDefaultClientOptions
}
