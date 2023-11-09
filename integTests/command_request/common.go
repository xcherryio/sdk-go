package command_request

const (
	testQueueName1 = "test-queue-1"
	testQueueName2 = "test-queue-2"
	testQueueName3 = "test-queue-3"
)

var testMyMsq = MyMessage{
	Str: "test-message",
	Int: 123,
}

type MyMessage struct {
	Str string
	Int int
}
