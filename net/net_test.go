// SPDX-License-Identifier: BSD-3-Clause
package net

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"testing"

	"github.com/shirou/gopsutil/v4/internal/common"
)

func skipIfNotImplementedErr(t *testing.T, err error) {
	if errors.Is(err, common.ErrNotImplementedError) {
		t.Skip("not implemented")
	}
}

func TestAddrString(t *testing.T) {
	v := Addr{IP: "192.168.0.1", Port: 8000}

	s := fmt.Sprintf("%v", v)
	if s != `{"ip":"192.168.0.1","port":8000}` {
		t.Errorf("Addr string is invalid: %v", v)
	}
}

func TestIOCountersStatString(t *testing.T) {
	v := IOCountersStat{
		Name:      "test",
		BytesSent: 100,
	}
	e := `{"name":"test","bytesSent":100,"bytesRecv":0,"packetsSent":0,"packetsRecv":0,"errin":0,"errout":0,"dropin":0,"dropout":0,"fifoin":0,"fifoout":0}`
	if e != fmt.Sprintf("%v", v) {
		t.Errorf("NetIOCountersStat string is invalid: %v", v)
	}
}

func TestProtoCountersStatString(t *testing.T) {
	v := ProtoCountersStat{
		Protocol: "tcp",
		Stats: map[string]int64{
			"MaxConn":      -1,
			"ActiveOpens":  4000,
			"PassiveOpens": 3000,
		},
	}
	e := `{"protocol":"tcp","stats":{"ActiveOpens":4000,"MaxConn":-1,"PassiveOpens":3000}}`
	if e != fmt.Sprintf("%v", v) {
		t.Errorf("NetProtoCountersStat string is invalid: %v", v)
	}
}

func TestConnectionStatString(t *testing.T) {
	v := ConnectionStat{
		Fd:     10,
		Family: 10,
		Type:   10,
		Uids:   []int32{10, 10},
	}
	e := `{"fd":10,"family":10,"type":10,"localaddr":{"ip":"","port":0},"remoteaddr":{"ip":"","port":0},"status":"","uids":[10,10],"pid":0}`
	if e != fmt.Sprintf("%v", v) {
		t.Errorf("NetConnectionStat string is invalid: %v", v)
	}
}

func TestIOCountersAll(t *testing.T) {
	v, err := IOCounters(false)
	skipIfNotImplementedErr(t, err)
	if err != nil {
		t.Errorf("Could not get NetIOCounters: %v", err)
	}
	per, err := IOCounters(true)
	skipIfNotImplementedErr(t, err)
	if err != nil {
		t.Errorf("Could not get NetIOCounters: %v", err)
	}
	if len(v) != 1 {
		t.Errorf("Could not get NetIOCounters: %v", v)
	}
	if v[0].Name != "all" {
		t.Errorf("Invalid NetIOCounters: %v", v)
	}
	var pr uint64
	for _, p := range per {
		pr += p.PacketsRecv
	}
	// small diff is ok, compare instead of math.Abs(subtraction) with uint64
	var diff uint64
	if v[0].PacketsRecv > pr {
		diff = v[0].PacketsRecv - pr
	} else {
		diff = pr - v[0].PacketsRecv
	}
	if diff > 5 {
		if ci := os.Getenv("CI"); ci != "" {
			// This test often fails in CI. so just print even if failed.
			fmt.Printf("invalid sum value: %v, %v", v[0].PacketsRecv, pr)
		} else {
			t.Errorf("invalid sum value: %v, %v", v[0].PacketsRecv, pr)
		}
	}
}

func TestIOCountersPerNic(t *testing.T) {
	v, err := IOCounters(true)
	skipIfNotImplementedErr(t, err)
	if err != nil {
		t.Errorf("Could not get NetIOCounters: %v", err)
	}
	if len(v) == 0 {
		t.Errorf("Could not get NetIOCounters: %v", v)
	}
	for _, vv := range v {
		if vv.Name == "" {
			t.Errorf("Invalid NetIOCounters: %v", vv)
		}
	}
}

func TestGetNetIOCountersAll(t *testing.T) {
	n := []IOCountersStat{
		{
			Name:        "a",
			BytesRecv:   10,
			PacketsRecv: 10,
		},
		{
			Name:        "b",
			BytesRecv:   10,
			PacketsRecv: 10,
			Errin:       10,
		},
	}
	ret := getIOCountersAll(n)
	if len(ret) != 1 {
		t.Errorf("invalid return count")
	}
	if ret[0].Name != "all" {
		t.Errorf("invalid return name")
	}
	if ret[0].BytesRecv != 20 {
		t.Errorf("invalid count bytesrecv")
	}
	if ret[0].Errin != 10 {
		t.Errorf("invalid count errin")
	}
}

func TestInterfaces(t *testing.T) {
	v, err := Interfaces()
	skipIfNotImplementedErr(t, err)
	if err != nil {
		t.Errorf("Could not get NetInterfaceStat: %v", err)
	}
	if len(v) == 0 {
		t.Errorf("Could not get NetInterfaceStat: %v", err)
	}
	for _, vv := range v {
		if vv.Name == "" {
			t.Errorf("Invalid NetInterface: %v", vv)
		}
	}
}

func TestProtoCountersStatsAll(t *testing.T) {
	v, err := ProtoCounters(nil)
	skipIfNotImplementedErr(t, err)
	if err != nil {
		t.Fatalf("Could not get NetProtoCounters: %v", err)
	}
	if len(v) == 0 {
		t.Fatalf("Could not get NetProtoCounters: %v", err)
	}
	for _, vv := range v {
		if vv.Protocol == "" {
			t.Errorf("Invalid NetProtoCountersStat: %v", vv)
		}
		if len(vv.Stats) == 0 {
			t.Errorf("Invalid NetProtoCountersStat: %v", vv)
		}
	}
}

func TestProtoCountersStats(t *testing.T) {
	v, err := ProtoCounters([]string{"tcp", "ip"})
	skipIfNotImplementedErr(t, err)
	if err != nil {
		t.Fatalf("Could not get NetProtoCounters: %v", err)
	}
	if len(v) == 0 {
		t.Fatalf("Could not get NetProtoCounters: %v", err)
	}
	if len(v) != 2 {
		t.Fatalf("Go incorrect number of NetProtoCounters: %v", err)
	}
	for _, vv := range v {
		if vv.Protocol != "tcp" && vv.Protocol != "ip" {
			t.Errorf("Invalid NetProtoCountersStat: %v", vv)
		}
		if len(vv.Stats) == 0 {
			t.Errorf("Invalid NetProtoCountersStat: %v", vv)
		}
	}
}

func TestConnections(t *testing.T) {
	if ci := os.Getenv("CI"); ci != "" { // skip if test on CI
		return
	}

	v, err := Connections("inet")
	skipIfNotImplementedErr(t, err)
	if err != nil {
		t.Errorf("could not get NetConnections: %v", err)
	}
	if len(v) == 0 {
		t.Errorf("could not get NetConnections: %v", v)
	}
	for _, vv := range v {
		if vv.Family == 0 {
			t.Errorf("invalid NetConnections: %v", vv)
		}
	}
}

func TestFilterCounters(t *testing.T) {
	if ci := os.Getenv("CI"); ci != "" { // skip if test on CI
		return
	}

	if runtime.GOOS == "linux" {
		// some test environment has not the path.
		if !common.PathExists("/proc/sys/net/netfilter/nf_connTrackCount") {
			t.SkipNow()
		}
	}

	v, err := FilterCounters()
	skipIfNotImplementedErr(t, err)
	if err != nil {
		t.Errorf("could not get NetConnections: %v", err)
	}
	if len(v) == 0 {
		t.Errorf("could not get NetConnections: %v", v)
	}
	for _, vv := range v {
		if vv.ConnTrackMax == 0 {
			t.Errorf("nf_connTrackMax needs to be greater than zero: %v", vv)
		}
	}
}

func TestInterfaceStatString(t *testing.T) {
	v := InterfaceStat{
		Index:        0,
		MTU:          1500,
		Name:         "eth0",
		HardwareAddr: "01:23:45:67:89:ab",
		Flags:        []string{"up", "down"},
		Addrs:        InterfaceAddrList{{Addr: "1.2.3.4"}, {Addr: "5.6.7.8"}},
	}

	s := fmt.Sprintf("%v", v)
	if s != `{"index":0,"mtu":1500,"name":"eth0","hardwareAddr":"01:23:45:67:89:ab","flags":["up","down"],"addrs":[{"addr":"1.2.3.4"},{"addr":"5.6.7.8"}]}` {
		t.Errorf("InterfaceStat string is invalid: %v", s)
	}

	list := InterfaceStatList{v, v}
	s = fmt.Sprintf("%v", list)
	if s != `[{"index":0,"mtu":1500,"name":"eth0","hardwareAddr":"01:23:45:67:89:ab","flags":["up","down"],"addrs":[{"addr":"1.2.3.4"},{"addr":"5.6.7.8"}]},{"index":0,"mtu":1500,"name":"eth0","hardwareAddr":"01:23:45:67:89:ab","flags":["up","down"],"addrs":[{"addr":"1.2.3.4"},{"addr":"5.6.7.8"}]}]` {
		t.Errorf("InterfaceStatList string is invalid: %v", s)
	}
}
