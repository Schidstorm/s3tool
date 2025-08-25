package s3lib

type ConnectorLoader interface {
	Load() ([]Connector, error)
}
