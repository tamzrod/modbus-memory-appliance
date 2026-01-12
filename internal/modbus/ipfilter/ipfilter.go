// internal/modbus/ipfilter/ipfilter.go
package ipfilter

import (
	"fmt"
	"net"
	"strings"
)

type Filter struct {
	allow []*net.IPNet
	deny  []*net.IPNet
	// If allow is non-empty => default deny.
	hasAllow bool
}

func Compile(allow, deny []string) (*Filter, error) {
	f := &Filter{}
	var err error

	f.allow, err = parseCIDRs(allow)
	if err != nil {
		return nil, fmt.Errorf("allow list: %w", err)
	}
	f.deny, err = parseCIDRs(deny)
	if err != nil {
		return nil, fmt.Errorf("deny list: %w", err)
	}
	f.hasAllow = len(f.allow) > 0
	return f, nil
}

func (f *Filter) Enabled() bool {
	return len(f.allow) > 0 || len(f.deny) > 0
}

func (f *Filter) Allowed(ip net.IP) bool {
	if ip == nil {
		return false
	}
	ip = normalizeIP(ip)

	// deny wins
	for _, n := range f.deny {
		if n.Contains(ip) {
			return false
		}
	}

	// if allow list exists, must match
	if f.hasAllow {
		for _, n := range f.allow {
			if n.Contains(ip) {
				return true
			}
		}
		return false
	}

	// no allow list => permit by default
	return true
}

func parseCIDRs(items []string) ([]*net.IPNet, error) {
	var out []*net.IPNet
	for _, raw := range items {
		s := strings.TrimSpace(raw)
		if s == "" {
			continue
		}

		// Single IP -> /32 or /128
		if !strings.Contains(s, "/") {
			ip := net.ParseIP(s)
			if ip == nil {
				return nil, fmt.Errorf("invalid ip %q", s)
			}
			out = append(out, ipToSingleIPNet(ip))
			continue
		}

		_, n, err := net.ParseCIDR(s)
		if err != nil {
			return nil, fmt.Errorf("invalid cidr %q: %w", s, err)
		}
		out = append(out, n)
	}
	return out, nil
}

func ipToSingleIPNet(ip net.IP) *net.IPNet {
	ip = normalizeIP(ip)
	bits := 32
	if ip.To4() == nil {
		bits = 128
	}
	return &net.IPNet{
		IP:   ip,
		Mask: net.CIDRMask(bits, bits),
	}
}

func normalizeIP(ip net.IP) net.IP {
	if v4 := ip.To4(); v4 != nil {
		return v4
	}
	return ip.To16()
}
