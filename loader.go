package uas

import (
	"encoding/xml"
	"io"
	"os"
	"regexp"
	"code.google.com/p/sre2/sre2"
)

var regMatcher *regexp.Regexp

func init() {
	regMatcher = sre2.MustParse("^/(?P<reg>.*)/(?P<flags>[imsU]*)\\s*$")
}

func compileReg(reg string) *sre2.SafeReader {
	return sre2.MustParse(regMatcher.ReplaceAllString(reg, "(?${flags}:${reg})"))
}

func compileBrowserRegs(regs []*BrowserReg) {
	for i, reg := range regs {
		regs[i].Reg = compileReg(reg.RegString)
	}
}

func compileOsRegs(regs []*OsReg) {
	for i, reg := range regs {
		regs[i].Reg = compileReg(reg.RegString)
	}
}

func compileDeviceRegs(regs []*DeviceReg) {
	for i, reg := range regs {
		regs[i].Reg = compileReg(reg.RegString)
	}
}

func mapBrowserTypeToBrowser(manifest *Manifest) {
	for _, browser := range manifest.Data.Browsers {
		browserType, found := manifest.GetBrowserType(browser.TypeId)
		if !found {
			browserType = manifest.OtherBrowserType()
		}
		browser.Type = browserType
	}
}

func mapOsToBrowser(manifest *Manifest) {
	for _, browser := range manifest.Data.Browsers {
		browser.Os, _ = manifest.GetOsForBrowser(browser.Id)
	}
}

// Creates a new Manifest instance by processing the XML from the provided Reader.
func Load(reader io.Reader) (*Manifest, error) {
	manifest := &Manifest{}
	if err := xml.NewDecoder(reader).Decode(manifest); err != nil {
		return nil, err
	}
	compileBrowserRegs(manifest.Data.BrowsersReg)
	compileOsRegs(manifest.Data.OperatingSystemsReg)
	compileDeviceRegs(manifest.Data.DevicesReg)
	mapBrowserTypeToBrowser(manifest)
	mapOsToBrowser(manifest)
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
