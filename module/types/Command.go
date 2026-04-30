package types

type LockUnlockCommand struct {
	Path   string `mapstructure:"path"`
	Status int    `mapstructure:"status"`
}

type CheckCommand struct {
	Path string `mapstructure:"path"`
}

type OnePathCommand struct {
	Path string `mapstructure:"path"`
}

type TwoPathCommand struct {
	Src  string `mapstructure:"src"`
	Dest string `mapstructure:"dest"`
}
