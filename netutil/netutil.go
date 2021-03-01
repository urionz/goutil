package netutil

import (
	"bytes"
	"net"
	"net/http"
	"strings"
)

// InternalIP get internal IP
func InternalIP() (ip string) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		panic("Oops: " + err.Error())
	}

	for _, a := range addrs {
		if ipNet, ok := a.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				// os.Stdout.WriteString(ipNet.IP.String() + "\n")
				ip = ipNet.IP.String()
				return
			}
		}
	}

	// os.Exit(0)
	return
}

func GetClientIP(r *http.Request, headers ...string) string {
	for _, header := range headers {
		ip := r.Header.Get(header)
		if ip != "" {
			return strings.Split(ip, ",")[0]
		}
	}
	return strings.Split(r.RemoteAddr, ":")[0]
}

func IPInRange(from, to, target string) bool {
	ipFrom := net.ParseIP(from)
	ipTo := net.ParseIP(to)
	ipTarget := net.ParseIP(target)

	if ipFrom == nil || ipTo == nil || ipTarget == nil {
		return false
	}

	from16 := ipFrom.To16()
	to16 := ipTo.To16()
	target16 := ipTarget.To16()
	if from16 == nil || to16 == nil || target16 == nil {
		return false
	}

	if bytes.Compare(target16, from16) >= 0 && bytes.Compare(target16, to16) <= 0 {
		return true
	}
	return false
}
