package xdb

type ClientOptions struct {
	ServerUrl     string
	WorkerUrl     string
	ObjectEncoder ObjectEncoder
	// TODO API timeout and retry policy
}

const DefaultWorkerPort = "8803"
const DefaultServerPort = "8801"
const (
	DefaultWorkerUrl = "http://localhost:" + DefaultWorkerPort
	DefaultServerUrl = "http://localhost:" + DefaultServerPort
)

var localDefaultClientOptions = ClientOptions{
	ServerUrl:     DefaultServerUrl,
	WorkerUrl:     DefaultWorkerUrl,
	ObjectEncoder: GetDefaultObjectEncoder(),
}

func GetLocalDefaultClientOptions() *ClientOptions {
	return &localDefaultClientOptions
}
