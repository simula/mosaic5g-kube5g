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
	PrintFunc(logger, "cmdName= "+cmdName)
	for i := 0; i < len(args); i++ {
		PrintFunc(logger, "args["+string(i)+"]="+args[i])
	}
	installSnap := cmd.NewCmd(cmdName, args...)
	finalStatus := <-installSnap.Start() // block and wait
	PrintFunc(logger, finalStatus)
	PrintFunc(logger, finalStatus.Cmd)

	return finalStatus
}

// RunCmdNonBlocking will run external commands in sync. Return stdout[0].
func RunCmdNonBlocking(logger *log.Logger, cmdName string, args ...string) {
	PrintFunc(logger, "cmdName="+cmdName)
	// logger.Print("cmdName=", cmdName)
	for i := 0; i < len(args); i++ {
		PrintFunc(logger, "args["+string(i)+"]="+args[i])
	}
	installSnap := cmd.NewCmd(cmdName, args...)
	finalStatus := installSnap.Start() // do not wait
	PrintFunc(logger, finalStatus)
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
			PrintFunc(logger, "Package: "+packageName+" Exist")
			return true, nil
		}

	}
	PrintFunc(logger, "Package: "+packageName+" does not Exist")
	return false, nil
}

//GetIPFromDomain will get the IP of the domain
func GetIPFromDomain(logger *log.Logger, domain string) (string, error) {
	addr, err := net.LookupHost(domain)
	if err != nil {
		PrintFunc(logger, "Failed to get IP from domain,err: ", err)
		return "", err
	}
	return addr[0], nil
}

// GetOutboundIP gets preferred outbound ip of this machine
func GetOutboundIP(logger *log.Logger) string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		PrintFuncFatal(logger, err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String()
}

// GetInterfaceByIP can get interface name from IP
func GetInterfaceByIP(logger *log.Logger, targetIP string) (string, error) {
	ifaces, err := net.Interfaces()
	// handle err
	if err != nil {
		return "", err
	}
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			PrintFunc(logger, err)
			continue
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

//PrintFunc will return if this package is already exist or not
func PrintFunc(logger *log.Logger, args ...interface{}) {
	switch len(args) {
	case 1:
		logger.Print(args[0])
		fmt.Println(args[0])
	case 2:
		logger.Print(args[0], args[1])
		fmt.Println(args[0], args[1])
	default:
		fmt.Println("lemlen(args)=", len(args), args)
		logger.Print("Unexpected number of variables")
		panic("Unexpected number of variables")
	}
}

//PrintFuncFatal will return if this package is already exist or not
func PrintFuncFatal(logger *log.Logger, args ...interface{}) {
	logger.Fatalln(args[0])
	panic(args[0])
}
