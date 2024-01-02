package ping_fx_module

import (
	"github.com/go-ping/ping"
	"go.uber.org/zap"
	"time"
)

type Ping struct {
	ip      string
	count   int
	timeout int
	maxLoss float64

	online bool

	log *zap.SugaredLogger
}

func NewPing(log *zap.SugaredLogger) *Ping {
	cfg, err := loadConfig(cfgPath)
	if err != nil {
		panic(err)
	}

	return &Ping{
		ip:      cfg.IP,
		count:   cfg.Count,
		timeout: cfg.Timeout,
		maxLoss: cfg.MaxLoss,
		online:  false,
		log:     log,
	}
}

func (p *Ping) newPinger() *ping.Pinger {
	pinger, err := ping.NewPinger(p.ip)
	if err != nil {
		p.log.Errorf("Failed to create pinger: %v", err)
	}

	pinger.SetPrivileged(true) // Must be set for Windows
	pinger.Count = p.count
	pinger.Timeout = time.Duration(p.timeout) * time.Second
	return pinger
}

func (p *Ping) checkOnline() bool {
	pinger := p.newPinger()

	if err := pinger.Run(); err != nil {
		p.log.Errorf("Failed to run pinger: %v", err)
		return false
	}

	p.log.Infof("sent: %v, recv: %v, loss: %v", pinger.PacketsSent, pinger.PacketsRecv, pinger.Statistics().PacketLoss)
	return pinger.Statistics().PacketLoss < p.maxLoss
}

func (p *Ping) Run(interval int, ch chan bool) {
	p.online = p.checkOnline()
	ch <- p.online

	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		curStatus := p.checkOnline()
		p.log.Infof("ping %s, online %v", p.ip, curStatus)
		if curStatus != p.online {
			p.online = curStatus
			ch <- p.online
		}
	}
}
