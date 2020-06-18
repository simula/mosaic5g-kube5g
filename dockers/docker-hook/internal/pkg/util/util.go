package util

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/go-cmd/cmd"
)

// RunCmd will run external commands in sync. Return stdout[0].
func RunCmd(logger *log.Logger, cmdName string, args ...string) cmd.Status {
	fmt.Println("cmdName=", cmdName)
	// logger.Print("cmdName=", cmdName)
	for i := 0; i < len(args); i++ {
		fmt.Println("args[", i, "]=", args[i])
		// logger.Print("args[", i, "]=", args[i])
	}
	installSnap := cmd.NewCmd(cmdName, args...)
	finalStatus := <-installSnap.Start() // block and wait
	logger.Print("finalStatus=", finalStatus)
	logger.Print("finalStatus.Cmd=", finalStatus.Cmd)
	return finalStatus

	// logger.Print(finalStatus.Cmd)
	// logger.Print(finalStatus)

}

// RunCmdNonBlocking will run external commands in sync. Return stdout[0].
func RunCmdNonBlocking(logger *log.Logger, cmdName string, args ...string) {
	fmt.Println("cmdName=", cmdName)
	// logger.Print("cmdName=", cmdName)
	for i := 0; i < len(args); i++ {
		fmt.Println("args[", i, "]=", args[i])
		// logger.Print("args[", i, "]=", args[i])
	}
	installSnap := cmd.NewCmd(cmdName, args...)
	finalStatus := installSnap.Start() // do not wait
	logger.Print("finalStatus=", finalStatus)

	// logger.Print(finalStatus.Cmd)
	// logger.Print(finalStatus)

}

//CheckSnapPackageExist will return if this package is already exist or not
func CheckSnapPackageExist(logger *log.Logger, packageName string) (bool, error) {
	if len(packageName) <= 0 {
		return false, errors.New("Input package name is empty")
	}
	retStatus := RunCmd(logger, "snap", "list")
	if retStatus.Exit != 0 {
		return false, errors.New("snap list return non-zero")
	}
	for i := 0; i < len(retStatus.Stdout); i++ {
		if strings.Contains(retStatus.Stdout[i], packageName) {
			logger.Println("Package: ", packageName, " Exist")
			return true, nil
		}

	}
	logger.Println("Package: ", packageName, " does not Exist")
	return false, nil
}

//GetInterfaceIP will get the ip of the interface. If failed, it'll return a default (127.0.1.10) value
func GetInterfaceIP(logger *log.Logger, interfaceName string) (string, error) {
	ret := RunCmd(logger, "ifconfig", interfaceName)
	if ret.Exit != 0 {
		return "127.0.1.10", errors.New("Fail to run ifconfig")
	}
	if len(ret.Stdout) <= 0 {
		return "127.0.1.10", errors.New("Fail to get result")
	}
	i := 0
	space := " "
	for {
		if ret.Stdout[1][27+i+1] == space[0] {
			break
		}
		i++
	}
	return ret.Stdout[1][20 : 27+i+1], nil
}

//GetIPFromDomain will get the IP of the domain
func GetIPFromDomain(logger *log.Logger, domain string) (string, error) {
	addr, err := net.LookupHost(domain)
	if err != nil {
		logger.Print("Failed to get IP from domain,err: ", err)
		return "", err
	}
	return addr[0], nil
}

// GetOutboundIP gets preferred outbound ip of this machine
func GetOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String()
}

// GetInterfaceByIP can get interface name from IP
func GetInterfaceByIP(targetIP string) (string, error) {
	ifaces, err := net.Interfaces()
	// handle err
	if err != nil {
		return "", err
	}
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			return "", err
		}
		// handle err
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip.String() == targetIP {

				return i.Name, nil
			}
			// process IP address
		}
	}
	return "", err
}
