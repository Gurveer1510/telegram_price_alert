package types

type Alert struct {
	ID               int
	Instrument_token int
	Instrument_name  string
	ChatId           int64
	Exchange         string
	Trigger_price    float64
	Condition        string
}
