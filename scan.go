package scan

import (
	"arping"
	"encoding/binary"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func Isport(port string) bool {
	var validator = regexp.MustCompile("^((6553[0-5])|(655[0-2][0-9])|(65[0-4][0-9]{2})|(6[0-4][0-9]{3})|([1-5][0-9]{4})|([0-5]{0,5})|([0-9]{1,4}))$")
	return validator.MatchString(port)
}

func domainchecker(hostname string) bool {
	domain := regexp.MustCompile(`^(([a-zA-Z]{1})|([a-zA-Z]{1}[a-zA-Z]{1})|([a-zA-Z]{1}[0-9]{1})|([0-9]{1}[a-zA-Z]{1})|([a-zA-Z0-9][a-zA-Z0-9-_]{1,61}[a-zA-Z0-9]))\.([a-zA-Z]{2,6}|[a-zA-Z0-9-]{2,30}\.[a-zA-Z
 ]{2,3})$`)
	return domain.MatchString(hostname)
}

func Get_ip(ip string) []string {
	var (
		err    error
		ips    []net.IP
		all_ip []string
	)
	if domainchecker(ip) {
		ips, err = net.LookupIP(ip)
		for _, ip := range ips {
			if ipv4 := ip.To4(); ipv4 != nil {
				all_ip = append(all_ip, ipv4.String())
			}

		}
		if err != nil {
			ip = ""
		}
	} else if net.ParseIP(ip) == nil {
		ip = ""
	} else {
		all_ip = append(all_ip, ip)
	}
	return all_ip
}

func socket(ip string, port int) (socket string) {
	socket = ip + ":" + strconv.Itoa(port)
	return socket
}

func Tcp_scan(ip string, port int) int {
	connection, err := net.DialTimeout("tcp", socket(ip, port), 1*time.Second)
	if err != nil {
		if strings.Contains(err.Error(), "too many open files") {
			time.Sleep(1 * time.Second)
			print(err)
			Tcp_scan(ip, port)
		} else {
			return 0
		}
	}

	defer connection.Close()
	return port

}
func Cdirgetter(cidr string) ([]string, error) {
	var hosts []string
	_, subnet, err := net.ParseCIDR(cidr)
	mascara := binary.BigEndian.Uint32(subnet.Mask)
	fAddr := binary.BigEndian.Uint32(subnet.IP)
	lAddr := (fAddr & mascara) | (mascara ^ 0xffffffff)
	for i := fAddr; i <= lAddr; i++ {
		ip := make(net.IP, 4)
		binary.BigEndian.PutUint32(ip, i)
		ips := ip.String()
		hosts = append(hosts, ips)
	}
	return hosts, err
}
func Arpscan_lan(ips string) string {
	ip := net.ParseIP(ips)
	_, _, err := arping.Ping(ip)
	if err == arping.ErrTimeout {
		return ""
	} else if err != nil {
		if strings.Contains(err.Error(), "operation not") {
			print("Please run as root\n")
		} else if strings.Contains(err.Error(), "ip+net") {
			return "Fail in net resources occurred Running again" + "\n"
			Arpscan_lan(ips)

		} else if strings.Contains(err.Error(), "no usable interface found") {
			return "Probably you put a CIDR outside ur net" + "\n"
			os.Exit(1)
		} else {
			return "Running again: Unknown Error succedeed for " + ips + "\n"
			Arpscan_lan(ips)
		}
	} else {
		return ips
	}
	return "Error"
}
