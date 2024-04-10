package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
	"errors"

	"github.com/ariary/quicli/pkg/quicli"
)

func main() {
	quicli.Run(quicli.Cli{
		Usage:       "testnet [flags]",
		Description: "test network connectivity",
		Flags: quicli.Flags{
			{Name: "url", Default: "", Description: "destination endpoint to test connectivity (hostname/ip:port)"},
			{Name: "timeout", Default: 1, Description: "specify timeout in second for testing request"},
			{Name: "serve", Default: false, Description: "wait testnet connection with an http server"},
		},
		Function: Testnet,
	})
}

func Testnet(cfg quicli.Config) {
	//detect pod ip
	fromIp,err:= externalIP()
	if err != nil {
		fmt.Println("failed to reterieve pod IP")
		fromIp="<unknown>"
	}

	// server mode
	if cfg.GetBoolFlag("serve") {
		http.HandleFunc("/", waitTests(fromIp))
		http.ListenAndServe(":9292", nil)
	}
	//simple check
	url := cfg.GetStringFlag("url")
	if url == "" {
		fmt.Fprintf(os.Stderr, "Usage: you must specify an url with -u/--url\n")
		os.Exit(1)
	}
	timeout := cfg.GetIntFlag("timeout")
	success, failed := conectivityCheck(url, timeout, fromIp)
	if failed != "" {
		fmt.Fprintf(os.Stderr, "%s", failed)
		os.Exit(1)
	}
	fmt.Print(success)

}

func conectivityCheck(endpoint string, timeout int, fromIp string) (success, failed string) {
	conn, err := net.DialTimeout("tcp", endpoint, time.Duration(timeout)*time.Second)
	if err != nil {
		return "", fmt.Sprintln("Connection failed to", endpoint, "from", fromIp, ":", err)
	}
	defer conn.Close()

	//test write to conn
	if err:= conn.SetWriteDeadline(time.Now().Add(time.Duration(timeout) * time.Second));err!=nil{
		return "", fmt.Sprintln("Connection failed to", endpoint, "from", fromIp, ":", err)
	}
	if _, err := conn.Write([]byte("testnet")); err != nil {
		return "", fmt.Sprintln("Connection failed to", endpoint, "from", fromIp, ":", err)
	}

	return fmt.Sprintln("Connection successful to", endpoint, "from", fromIp), ""
}

func waitTests(fromIp string)http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request){
		args := strings.Split(r.URL.Path, "/")[1:] // always starts with /
		if len(args) > 0 {
			if len(args) == 1 {
				if args[0] == "shutdown" { // shutdown
					fmt.Fprintf(w, "Well received boss!")
					fmt.Println("shutdown server ..")
					os.Exit(0)
				}
			} else if len(args) == 2 { // check
				endpoint := args[0] + ":" + args[1]
				success, failed := conectivityCheck(endpoint, 1,fromIp)
				if failed != "" {
					fmt.Print(failed)
				} else {
					fmt.Print(success)
				}
			}
		}
	}
}


func externalIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String(), nil
		}
	}
	return "", errors.New("are you connected to the network?")
}
