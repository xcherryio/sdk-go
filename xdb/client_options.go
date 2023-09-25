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
	defaultWorkerUrl = "http://localhost:" + DefaultWorkerPort
	defaultServerUrl = "http://localhost:" + DefaultServerPort
)

var localDefaultClientOptions = ClientOptions{
	ServerUrl:     defaultServerUrl,
	WorkerUrl:     defaultWorkerUrl,
	ObjectEncoder: GetDefaultObjectEncoder(),
}

func GetLocalDefaultClientOptions() *ClientOptions {
	return &localDefaultClientOptions
}
