package system

import (
	"fmt"
	"syscall"

	"github.com/qiniu/logkit/metric"
)

type NetStats struct {
	ps PS
}

func (_ *NetStats) Name() string {
	return "netstat"
}

func (s *NetStats) Collect() (datas []map[string]interface{}, err error) {
	netconns, err := s.ps.NetConnections()
	if err != nil {
		return nil, fmt.Errorf("error getting net connections info: %s", err)
	}
	counts := make(map[string]int)
	counts["UDP"] = 0

	// TODO: add family to results
	for _, netcon := range netconns {
		if netcon.Type == syscall.SOCK_DGRAM {
			counts["UDP"] += 1
			continue // UDP has no status
		}
		c, ok := counts[netcon.Status]
		if !ok {
			counts[netcon.Status] = 0
		}
		counts[netcon.Status] = c + 1
	}

	fields := map[string]interface{}{
		"tcp_established": counts["ESTABLISHED"],
		"tcp_syn_sent":    counts["SYN_SENT"],
		"tcp_syn_recv":    counts["SYN_RECV"],
		"tcp_fin_wait1":   counts["FIN_WAIT1"],
		"tcp_fin_wait2":   counts["FIN_WAIT2"],
		"tcp_time_wait":   counts["TIME_WAIT"],
		"tcp_close":       counts["CLOSE"],
		"tcp_close_wait":  counts["CLOSE_WAIT"],
		"tcp_last_ack":    counts["LAST_ACK"],
		"tcp_listen":      counts["LISTEN"],
		"tcp_closing":     counts["CLOSING"],
		"tcp_none":        counts["NONE"],
		"udp_socket":      counts["UDP"],
	}
	datas = append(datas, fields)
	return
}

func init() {
	metric.Add("netstat", func() metric.Collector {
		return &NetStats{ps: newSystemPS()}
	})
}
