package main

import (
	"fmt"
	"net"
	"time"

	"github.com/cmatsuoka/ouidb"
	"github.com/mostlygeek/arp"
)

var db *ouidb.OuiDB

func init() {
	db = ouidb.New("/etc/manuf")
	if db == nil {
		db = ouidb.New("manuf")
	}
}

func getMAC(s string) (string, error) {
	ifaces, err := net.Interfaces()
	checkError(err)
	for _, i := range ifaces {
		if i.Name == s {
			return i.HardwareAddr.String(), nil
		}
	}
	return "", fmt.Errorf("%s: no such interface", s)
}

func getName(addr string) string {
	names, err := net.LookupAddr(addr)
	if err != nil {
		return ""
	}
	return names[0]
}

func getMACFromIP(addr string) string {
	arp.CacheUpdate()
	return arp.Search(addr)
}

func getVendor(mac string) string {
	v, _ := db.Lookup(mac)
	return v

}

func formatDuration(t time.Duration) string {
	s := int(t.Seconds()) % 60
	m := int(t.Minutes()) % 60
	h := int(t.Hours()) % 24
	d := int(t.Hours()) / 24
	return fmt.Sprintf("%dd%dh%dm%ds", d, h, m, s)
}
