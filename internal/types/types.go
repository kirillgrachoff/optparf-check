package types

type Query struct {
	PeerId  int64  `mapstructure:"peer_id"`
	Pattern string `mapstructure:"pattern"`
}

type Item struct {
	Id      string
	Article string
	Name    string
	Price   float64
	URL     string
}

type QueryResult struct {
	PeerId  int64
	Pattern string
	Items   []Item
	Err     error
}
