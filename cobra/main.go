package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var (
	startPort int
	endPort   int
	showEmpty bool
	debug     bool
)

type portInfo struct {
	isOpen  bool
	service string
}

func main() {
	var rootCmd = &cobra.Command{
		Use:   "portscan",
		Short: "A port scanner CLI application",
		Long:  `A port scanner that checks for open ports in a given range and displays either empty or occupied ports.`,
		Run:   runScan,
	}

	rootCmd.Flags().IntVarP(&startPort, "start", "s", 1, "Start port for scanning")
	rootCmd.Flags().IntVarP(&endPort, "end", "e", 1024, "End port for scanning")
	rootCmd.Flags().BoolVarP(&showEmpty, "empty", "m", false, "Show empty ports instead of occupied")
	rootCmd.Flags().BoolVarP(&debug, "debug", "d", false, "Enable debug output")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func runScan(cmd *cobra.Command, args []string) {
	if startPort > endPort {
		fmt.Println("Error: Start port must be less than or equal to end port")
		os.Exit(1)
	}

	fmt.Printf("Scanning ports %d to %d...\n", startPort, endPort)

	results := make(map[int]portInfo)
	for port := startPort; port <= endPort; port++ {
		info := scanPort(port)
		results[port] = info
		if debug {
			fmt.Printf("Port %d: %v\n", port, info)
		} else {
			fmt.Printf("Progress: %d/%d\r", port-startPort+1, endPort-startPort+1)
		}
	}

	fmt.Println("\nScan complete!")

	portType := "occupied"
	if showEmpty {
		portType = "empty"
	}
	fmt.Printf("\nList of %s ports:\n", portType)

	count := 0
	for port := startPort; port <= endPort; port++ {
		info := results[port]
		if info.isOpen != showEmpty {
			if showEmpty {
				fmt.Printf("- %d\n", port)
			} else {
				fmt.Printf("- %d: %s\n", port, info.service)
			}
			count++
		}
	}

	if count == 0 {
		fmt.Printf("No %s ports found in the specified range.\n", portType)
	} else {
		fmt.Printf("\nTotal %s ports found: %d\n", portType, count)
	}
}

func scanPort(port int) portInfo {
	address := fmt.Sprintf("localhost:%d", port)
	conn, err := net.DialTimeout("tcp", address, 500*time.Millisecond)
	if err != nil {
		return portInfo{isOpen: false}
	}
	defer conn.Close()
	return portInfo{isOpen: true, service: getServiceName(port)}
}

func getServiceName(port int) string {
	switch port {
	case 80:
		return "HTTP"
	case 443:
		return "HTTPS"
	case 22:
		return "SSH"
	case 21:
		return "FTP"
	// Add more well-known ports as needed
	default:
		return fmt.Sprintf("Unknown (%d)", port)
	}
}
