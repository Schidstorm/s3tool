package s3lib

type MemoryLoader struct {
}

func (l *MemoryLoader) Load() ([]Connector, error) {
	return []Connector{
		&MemoryConnector{
			name: "Memory",
		},
	}, nil
}
