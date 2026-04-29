package types

type LockUnlockCommand struct {
	Path   string `mapstructure:"path"`
	Status int    `mapstructure:"status"`
}

type CheckCommand struct {
	Path string `mapstructure:"path"`
}
