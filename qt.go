package main

type QingTing struct {
	Config *Config
	Login  bool
}

func initQingTing(config *Config) QingTing {
	return QingTing{
		Config: config,
	}
}

func (qt QingTing) SaveConfig() {
	SaveConfig(qt.Config)
}
