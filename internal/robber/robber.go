package robber

// TODO: Comment
type Robber struct {
	Args   *Args
	Config *Config
	Logger *Logger
}

// TODO: Comment
func NewRobber() *Robber {
	r := &Robber{
		Args: parseArgs(),
	}
	r.Config = newConfig(r)
	r.Logger = newLogger()
	return r
}
