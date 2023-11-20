package xc

type WorkerOptions struct {
	ObjectEncoder ObjectEncoder
	DBConverter   DBConverter
}

func GetDefaultWorkerOptions() WorkerOptions {
	return WorkerOptions{
		ObjectEncoder: GetDefaultObjectEncoder(),
		DBConverter:   NewBasicDBConverter(),
	}
}
