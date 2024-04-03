package webdriver

import (
	"fmt"

	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

type ChromeDriver struct {
	BrowserName  string
	BrowserPath  string
	DebuggerAddr string
	DriverPort   int
	DriverAddr   string
}

func NewChromeDriver(browserName, browserPath, debuggerAddr, driverAddr string, driverPort int) ChromeDriver {
	return ChromeDriver{
		BrowserName:  browserName,
		BrowserPath:  browserPath,
		DebuggerAddr: debuggerAddr,
		DriverAddr:   driverAddr,
		DriverPort:   driverPort,
	}
}

func (c *ChromeDriver) StartChrome() (*selenium.Service, error) {
	opts := []selenium.ServiceOption{}
	return selenium.NewChromeDriverService(c.DriverAddr, c.DriverPort, opts...)
}

func (c *ChromeDriver) WebDriver() (selenium.WebDriver, error) {
	caps := selenium.Capabilities{
		"browserName": c.BrowserName,
	}

	chromeCaps := chrome.Capabilities{}

	if len(c.BrowserPath) > 0 {
		chromeCaps.Path = c.BrowserPath
	}
	if len(c.DebuggerAddr) > 0 {
		chromeCaps.DebuggerAddr = c.DebuggerAddr
	}
	caps.AddChrome(chromeCaps)
	return selenium.NewRemote(caps, fmt.Sprintf("http://%s:%d/wd/hub", c.DriverAddr, c.DriverPort))
}
