package xc

type ClientOptions struct {
	Namespace           string
	ServerUrl           string
	WorkerUrl           string
	ObjectEncoder       ObjectEncoder
	EnabledDebugLogging bool
	// DefaultProcessTimeoutSecondsOverride is used when StartProcess is called and
	// 1. no timeout specified in ProcessOptions(default as zero)
	// 2. no timeout specified in ProcessStartOptions(default as nil)
	// currently mainly for testing purpose
	DefaultProcessTimeoutSecondsOverride int32
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
}

func GetLocalDefaultClientOptions() *ClientOptions {
	return &localDefaultClientOptions
}
