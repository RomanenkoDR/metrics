package agn_types

type OptionsAgent struct {
	ServerAddress  string `env:"ADDRESS"`
	PollInterval   int    `env:"POLL_INTERVAL"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	RateLimit      int    `env:"RATE_LIMIT"`
	Key            string `env:"KEY"`
	KeyByte        []byte
	Encrypt        bool
}
