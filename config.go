package api

import "gitee.com/fast_api/api/dwarf"

type Config struct {
	listen                    string
	dwarf                     *dwarf.DwarfMaker
	loadPath                  *string
	caFile, certFile, keyFile string
}

func (c *Config) CaFile() string {
	return c.caFile
}

func (c *Config) SetCaFile(caFile string) {
	c.caFile = caFile
}

func (c *Config) CertFile() string {
	return c.certFile
}

func (c *Config) SetCertFile(certFile string) {
	c.certFile = certFile
}

func (c *Config) KeyFile() string {
	return c.keyFile
}

func (c *Config) SetKeyFile(keyFile string) {
	c.keyFile = keyFile
}

func (c *Config) Listen() string {
	return c.listen
}

func (c *Config) Dwarf() *dwarf.DwarfMaker {
	return c.dwarf
}

func (c *Config) LoadPath() *string {
	return c.loadPath
}

func (c *Config) SetDwarfMaker(dwarf *dwarf.DwarfMaker) {
	c.dwarf = dwarf
}

func (c *Config) SetListen(ser string) {
	c.listen = ser
}

func (c *Config) AddIncludeRegex(regex ...string) {
	for _, s := range regex {
		c.dwarf.AddIncludeRegex(s)
	}
}

func (c *Config) AddExclude(pkg ...string) {
	for _, s := range pkg {
		c.dwarf.AddExclude(s)
	}
}

func (c *Config) SetDwarfMode(mode dwarf.FilterMode) {
	c.dwarf.SetFilterMode(mode)
}

func (c *Config) SetExePath(path string) {
	c.loadPath = &path
}
