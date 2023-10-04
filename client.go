package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func main() {
	fmt.Println("Get from Master")
	// get_chunck_req := "http://192.168.1.5:8090/fasta?id=1"
	// resp, err := http.Get(get_chunck_req)

	get_numOf_bases_req := "http://192.168.1.5:8090/fasta/baseCount"
	resp, err := http.Get(get_numOf_bases_req)

	panicOnError(err)
	defer resp.Body.Close()
	// get stautus code
	fmt.Println("Status code:", resp.StatusCode)
	b, err := ioutil.ReadAll(resp.Body)
	slaveIP := string(b)
	IPs := strings.Split(slaveIP, "\n")

	if strings.Contains(slaveIP, "http://") {
		for i := 1; i < len(IPs); i++ {
			fmt.Println("Get from Slave:" + IPs[i])
			resp, err = http.Get(IPs[i])
			panicOnError(err)
			defer resp.Body.Close()
			fmt.Println("Status code:", resp.StatusCode)
			b, err = ioutil.ReadAll(resp.Body)
			if IPs[0] == "0" {

				f, err := os.OpenFile("client"+IPs[0]+".fasta", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
				if err != nil {
					panicOnError(err)
				}
				defer f.Close()
				if _, err = f.WriteString(string(b)); err != nil {
					panicOnError(err)
				}
			} else if IPs[0] == "mapReduceResult" {
				err = ioutil.WriteFile("Nucleobases_Count.txt", b, 0644)
				panicOnError(err)
			} else {
				err = ioutil.WriteFile("client"+IPs[0]+".fasta", b, 0644)
				panicOnError(err)
			}
		}
	} else {
		fmt.Println("Master Returned :", slaveIP)
	}
}

func panicOnError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
