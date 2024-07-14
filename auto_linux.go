package main

import (
	"io"
	"log"

	"github.com/citilinkru/libudev"
	"github.com/citilinkru/libudev/matcher"
)

func NewAuto() (*CashDrawer, error) {
	sc := libudev.NewScanner()
	err, devices := sc.ScanDevices()
	if err != nil {
		log.Fatal(err)
	}

	m := matcher.NewMatcher()
	m.SetStrategy(matcher.StrategyAnd)
	m.AddRule(matcher.NewRuleEnv("ID_VENDOR", "Prolific_Technology_Inc."))
	m.AddRule(matcher.NewRuleEnv("DEVNAME", "tty"))

	filteredDevices := m.Match(devices)

	if len(filteredDevices) == 0 {
		return nil, io.EOF
	}

	return New("/dev/"+filteredDevices[0].Env["DevName"], 9600)
}
