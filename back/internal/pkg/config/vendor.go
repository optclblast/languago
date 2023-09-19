package config

type (
	Logger interface {
		Warn(kv ...interface{})
		Debug(kv ...interface{})
		Info(kv ...interface{})
		Log(kv ...interface{})
	}
)
