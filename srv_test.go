package main

import "testing"

func TestExeactUnit(t *testing.T) {
	s := ` UNIT                               LOAD   ACTIVE SUB     DESCRIPTION
	agv_bsc_centos.service             loaded active running agv_bsc_centos service
	atd.service                        loaded active running Job spooling tools
	auditd.service                     loaded active running Security Auditing Service
  ● charger_payment.service            loaded failed failed  charger_payment service
	charger_qrcode.service             loaded active running charger_qrcode service`
	t.Logf("%+v\n", ExtractServiceUnit(s))
}
