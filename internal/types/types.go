package types

type Query struct {
	PeerId  int64 `mapstructure:"peer_id"`
	Pattern string `mapstructure:"pattern"`
}

type QueryResult struct {
	PeerId int64
	Pattern string
	Found []byte
}