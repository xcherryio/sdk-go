package xc

type WorkerOptions struct {
	ObjectEncoder ObjectEncoder
}

func GetDefaultWorkerOptions() WorkerOptions {
	return WorkerOptions{
		ObjectEncoder: GetDefaultObjectEncoder(),
	}
}
