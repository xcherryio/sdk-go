package xdb

type LocalQueuePublishOptions struct {
	// DedupSeed is the seed to generate the DedupUUID
	// by uuid.NewMD5(uuid.NameSpaceOID, []byte(*DedupSeed))
	// only used if DedupUUID is nil
	DedupSeed *string
	// DedupUUID is the deduplication UUID
	DedupUUID *string
}

type LocalQueuePublishMessage struct {
	QueueName string
	Payload   interface{}
	// DedupSeed is the seed to generate the DedupUUID
	// by uuid.NewMD5(uuid.NameSpaceOID, []byte(*DedupSeed))
	// only used if DedupUUID is nil
	DedupSeed *string
	// DedupUUID is the deduplication UUID
	DedupUUID *string
}
