package config

var Cfg *Config

type Config struct {
	App Application `toml:"application"`
	Log LogCfg      `toml:"log"`
	Ok  OkConfig    `toml:"ok"`
}

type Application struct {
	Env       string `toml:"env"`
	Namespace string `toml:"namespace"`
	Exchange  string `toml:"exchange"`
}

type LogCfg struct {
	Format      string `toml:"format"`
	Level       string `toml:"level"`
	TimeFormat  string `toml:"timeFormat"`
	LogFilePath string `toml:"logFilePath"`
}

type OkConfig struct {
	PubWs       string  `toml:"pub_ws"`
	PriWs       string  `toml:"pri_ws"`
	PubTradeWs  string  `toml:"pub_trade_ws"`
	Key         string  `toml:"key"`
	Secret      string  `toml:"secret"`
	Passphase   string  `toml:"passphase"`
	PriceOffset float64 `toml:"price_offset"`
}
