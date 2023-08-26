package helpers

import (
	"io"
	"net"
	"net/http"
)

func GetPublicIP() (string, error) {

	//ipify API returns IPv4 IP address in plain/text
	res, err := http.Get("https://api.ipify.org/")
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	//Read response bytes
	ip, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(ip), nil
}

func GetLocalIP() (net.IP, error) {
	
	//Connect to 8.8.8.8
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	//Get local address from the connection 
	localAddress := conn.LocalAddr().(*net.UDPAddr)

	return localAddress.IP, nil
}

func GetLocalIPs() ([]net.IP, error) {

	var ips []net.IP

	//Get all interface addresses
	addresses, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	//Loop throught them and get the IPv4 addresses
	for _, addr := range addresses {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ips = append(ips, ipnet.IP)
			}
		}
	}
	
	return ips, nil
}
