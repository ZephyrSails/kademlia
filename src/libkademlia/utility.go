package libkademlia

import (
  "net"
  "strconv"
  "net/rpc"
  "log"
)


func StringToIpPort(laddr string) (ip net.IP, port uint16, err error) {
	hostString, portString, err := net.SplitHostPort(laddr)
	if err != nil {
		return
	}
	ipStr, err := net.LookupHost(hostString)
	if err != nil {
		return
	}
	for i := 0; i < len(ipStr); i++ {
		ip = net.ParseIP(ipStr[i])
		if ip.To4() != nil {
			break
		}
	}
	portInt, err := strconv.Atoi(portString)
	port = uint16(portInt)
	return
}

func IpPortToString(ip net.IP, port uint16) (ipStr string, portStr string) {
  ipStr = ip.String()
  portStr = strconv.Itoa(int(port))
  return
}

func getClient(host net.IP, port uint16) (client *rpc.Client) {
  hostStr, portStr := IpPortToString(host, port)

	client, err := rpc.DialHTTPPath("tcp", net.JoinHostPort(hostStr, portStr),
		rpc.DefaultRPCPath+hostStr+portStr)
	if err != nil {
		log.Fatal("DialHTTP: ", err)
	}
  return
}
