package dpdkinit

import (
	"github.com/njcx/gopacket131_dpdk/dpdk"
	"strconv"
)

var DPdk *dpdk.DPDKHandle

func DpdkInit() error {

	dpdkPort, _ := Parse("dpdk_port")
	dpdkStatus, _ := Parse("dpdk_status")
	if dpdkStatus == "enable" {
		err := dpdk.InitDPDK([]string{})
		if err != nil {
			return err
		}
		num16, _ := strconv.ParseUint(dpdkPort, 10, 16)
		port := uint16(num16)

		h, err := dpdk.NewDPDKHandle(port)
		if err != nil {
			h.Close()
			return err
		}
		DPdk = h
	}
	return nil
}
