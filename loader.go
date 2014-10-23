package uas

import (
	"encoding/xml"
	"io"
	"os"
	"code.google.com/p/sre2/sre2"
)

// Error message doesn't like this being a pointer so I pulled the '*' and it cleared the error.
// I don't understand enough about pointers to understand why.
var regMatcher sre2.Re

//
func init() {
	regMatcher = sre2.MustParse("^/(?P<reg>.*)/(?P<flags>[imsU]*)\\s*$")
}

// return type should probably be something from the sre2 package, but what I do not know
func compileReg(reg string) *sre2.SafeReader {
	// return sre2.MustParse(regexp.ReplaceAllString(reg, "(?${flags}:${reg})"))
	return regMatcher.MatchIndex(?${flags}:${reg})
}

func compileBrowserRegs(data *Data) {
	regs := data.BrowsersReg
	for i, reg := range regs {
		regs[i].Reg = compileReg(reg.RegString)
	}
}

func compileOsRegs(data *Data) {
	regs := data.OperatingSystemsReg
	for i, reg := range regs {
		regs[i].Reg = compileReg(reg.RegString)
	}
}

func compileDeviceRegs(data *Data) {
	regs := data.DevicesReg
	for i, reg := range regs {
		regs[i].Reg = compileReg(reg.RegString)
	}
}

// Creates a new Manifest instance by processing the XML from the provided Reader.
func Load(reader io.Reader) (*Manifest, error) {
	manifest := &Manifest{}
	if err := xml.NewDecoder(reader).Decode(manifest); err != nil {
		return nil, err
	}
	compileBrowserRegs(manifest.Data)
	compileOsRegs(manifest.Data)
	compileDeviceRegs(manifest.Data)
	return manifest, nil
}

// Creates a new Manifest instance by processing the XML from the provided file.
func LoadFile(path string) (*Manifest, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return Load(file)
}
