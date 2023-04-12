package main

import (
	"embed"
	_ "embed"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"time"
)

type RetraceInfo struct {
	Mapping string `json:"mapping"`
	Log     string `json:"log"`
}

var version = "0.0.1"
var defalutPort = "8081"

func main() {
	port := flag.String("p", defalutPort, "端口号默认端口("+defalutPort+")")
	mapping := flag.String("m", "", "混淆mapping文件路径")
	flag.Parse()
	http.HandleFunc("/", index)
	http.HandleFunc("/retrace", retrace)
	ips, err := LocalIPv4s()
	if err != nil {
		panic(err)
	}
	fmt.Println("\n地址:")
	fmt.Printf("\thttp://127.0.0.1:%s\n", *port)
	for _, ip := range ips {
		fmt.Printf("\thttp://%s:%s\n", ip, *port)
	}
	go func() {
		after := time.After(time.Second * 2)
		<-after
		if runtime.GOOS == "darwin" {
			arg := "http://127.0.0.1:" + (*port)
			if mappingAbs, err := filepath.Abs(*mapping); err == nil {
				arg += "?path=" + mappingAbs
			}
			command := exec.Command("open", arg)
			command.Run()
		}
	}()
	panic(http.ListenAndServe(":"+(*port), nil))
}

func isNum(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func retrace(w http.ResponseWriter, r *http.Request) {
	// 检查是否为post请求
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "invalid_http_method")
		return
	}
	defer r.Body.Close()
	con, _ := io.ReadAll(r.Body) //获取post的数据
	requestInfo := RetraceInfo{}
	json.Unmarshal(con, &requestInfo)

	tempFile, err := os.CreateTemp("", "log")
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	tempFile.WriteString(requestInfo.Log)
	tempFile.Close()
	defer os.Remove(tempFile.Name())
	command := exec.Command("retrace", requestInfo.Mapping, tempFile.Name())
	androidHome := os.Getenv("ANDROID_HOME")
	if len(androidHome) > 0 {
		command.Dir = path.Join(androidHome, "cmdline-tools", "latest", "bin")
	}
	output, err := command.CombinedOutput()
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	w.Write(output)
}

//go:embed axios.min.js
var axios embed.FS

func indexFile(w http.ResponseWriter, r *http.Request) {
	http.FileServer(http.FS(axios)).ServeHTTP(w, r)
}

//go:embed input.html
var indexByte []byte

func index(w http.ResponseWriter, r *http.Request) {
	w.Write(indexByte)
}

// LocalIPs return all non-loopback IPv4 addresses
func LocalIPv4s() ([]string, error) {
	var ips []string
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ips, err
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			ips = append(ips, ipnet.IP.String())
		}
	}

	return ips, nil
}
