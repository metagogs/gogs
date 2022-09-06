package snow

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

func GetLocalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		return "", err
	}

	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}

	return "", fmt.Errorf("can not get the ip")
}

func IP4toInt16(ip string) int64 {
	bits := strings.Split(ip, ".")

	b2, _ := strconv.Atoi(bits[2])
	b3, _ := strconv.Atoi(bits[3])

	var sum int64

	sum += int64(b2) << 8
	sum += int64(b3)

	return sum
}
