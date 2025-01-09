package srv_types

type OptionsServer struct {
	Address  string `env:"ADDRESS"`
	Interval int    `env:"STORE_INTERVAL"`
	Filename string `env:"FILE_STORAGE_PATH"`
	Restore  bool   `env:"RESTORE"`
	DBDSN    string `env:"DATABASE_DSN"`
	Key      string `env:"KEY"`
}
