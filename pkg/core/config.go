package core

func (d *Driver) GetDefaultBuildReference() string {
	return d.conf.Defaults.Build
}
