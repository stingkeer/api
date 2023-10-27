package api

import "gitee.com/fast_api/api/dwarf"

type ServerConfig struct {
	listen                                string
	dwarf                                 *dwarf.DwarfMaker
	loadPath                              *string
	caPEMBlock, certPEMBlock, keyPEMBlock []byte
}

func (c *ServerConfig) Listen() string {
	return c.listen
}

func (c *ServerConfig) Dwarf() *dwarf.DwarfMaker {
	return c.dwarf
}

func (c *ServerConfig) LoadPath() *string {
	return c.loadPath
}

func (c *ServerConfig) SetDwarfMaker(dwarf *dwarf.DwarfMaker) {
	c.dwarf = dwarf
}

func (c *ServerConfig) SetListen(ser string) {
	c.listen = ser
}

func (c *ServerConfig) AddIncludeRegex(regex ...string) {
	for _, s := range regex {
		c.dwarf.AddIncludeRegex(s)
	}
}

func (c *ServerConfig) AddExclude(pkg ...string) {
	for _, s := range pkg {
		c.dwarf.AddExclude(s)
	}
}

func (c *ServerConfig) SetDwarfMode(mode dwarf.FilterMode) {
	c.dwarf.SetFilterMode(mode)
}

func (c *ServerConfig) SetExePath(path string) {
	c.loadPath = &path
}
