package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

// Struct to store system info
type SystemInfo struct {
	HostName     string
	CPUName      string
	DiskCapacity uint64
	DiskUsage    uint64
	DiskFree     uint64
	RamCapacity  uint64
	RamAvailable uint64
}

func getSystemInfo(data *SystemInfo) {
	// Get cpu, disk, and host info
	cpuInfo, _ := cpu.Info()
	hostInfo, _ := host.Info()
	ramInfo, _ := mem.VirtualMemory()
	var diskInfo *disk.UsageStat

	switch hostInfo.OS {
	case "darwin":
		diskInfo, _ = disk.Usage("/")
	case "windows":
		diskInfo, _ = disk.Usage("\\")
	default:
		diskInfo, _ = disk.Usage("/")
	}

	// Store data into struct
	data.HostName = hostInfo.Hostname
	data.CPUName = cpuInfo[0].ModelName
	data.DiskCapacity = diskInfo.Total / 1024 / 1024
	data.DiskUsage = diskInfo.Used / 1024 / 1024
	data.DiskFree = diskInfo.Free / 1024 / 1024
	data.RamCapacity = ramInfo.Total / 1024 / 1024
	data.RamAvailable = ramInfo.Available / 1024 / 1024
}

// Function to Jsonify data nd write it to a file
func saveData(data *SystemInfo) {
	jsonify, _ := json.Marshal(data)
	_ = os.WriteFile("output.json", []byte(jsonify), 0644)
}

// Homepage Route
func homePageHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Write([]byte("<h1>Hello World!</h1>"))
}

// Route to get data and return as json
func getDataHandler(w http.ResponseWriter, r *http.Request) {
	data := SystemInfo{}
	getSystemInfo(&data)
	saveData(&data)
	json.NewEncoder(w).Encode(data)
}

// Request handler
func handleRequests(mux *http.ServeMux) {
	mux.HandleFunc("/", homePageHandler)
	mux.HandleFunc("/api", getDataHandler)
	log.Fatal(http.ListenAndServe(":3000", mux))
}

func main() {
	// Run the server
	mux := http.NewServeMux()
	fmt.Println("Server running at localhost:3000")
	handleRequests(mux)
}
