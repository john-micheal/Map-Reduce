package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type SlaveDevice struct {
	id     string
	ipAddr string
	data   string
}

var slave1 = SlaveDevice{id: "1", ipAddr: "http://192.168.1.5:8091/fasta"}
var slave2 = SlaveDevice{id: "2", ipAddr: "http://192.168.1.5:8092/fasta"}
var slave3 = SlaveDevice{id: "3", ipAddr: "http://192.168.1.5:8093/fasta"}
var slave4 = SlaveDevice{id: "4", ipAddr: "http://192.168.1.5:8094/fasta"}
var slave5 = SlaveDevice{id: "5", ipAddr: "http://192.168.1.5:8095"}
var reducer_Response = ""

func make_chunks(fileName string) {
	bytes, err := ioutil.ReadFile("Genome.fasta")
	panicOnError(err)
	// Allocate data to each slave device
	chunkSize := len(bytes) / 4

	slave1.data = string(bytes[:chunkSize])
	slave2.data = string(bytes[chunkSize : chunkSize*2])
	slave3.data = string(bytes[chunkSize*2 : chunkSize*3])
	slave4.data = string(bytes[chunkSize*3:])

}
func divide_chunks_on_slaves() {

	make_chunks("Big_Data.txt")

	{
		b := strings.NewReader(slave1.data)
		resp, err := http.Post(slave1.ipAddr, "text/plain", b)
		panicOnErrorM(err)
		defer resp.Body.Close()
		// get stautus code
		fmt.Println("Status code:", resp.StatusCode)
	}

	{
		b := strings.NewReader(slave2.data)
		resp, err := http.Post(slave2.ipAddr, "text/plain", b)
		panicOnErrorM(err)
		defer resp.Body.Close()

		// get stautus code
		fmt.Println("Status code:", resp.StatusCode)
	}

	{
		b := strings.NewReader(slave3.data)
		resp, err := http.Post(slave3.ipAddr, "text/plain", b)
		panicOnErrorM(err)
		defer resp.Body.Close()

		// get stautus code
		fmt.Println("Status code:", resp.StatusCode)
	}

	{
		b := strings.NewReader(slave4.data)
		resp, err := http.Post(slave4.ipAddr, "text/plain", b)
		panicOnErrorM(err)
		defer resp.Body.Close()

		// get stautus code
		fmt.Println("Status code:", resp.StatusCode)
	}

}

func map_request() {

	// mapResult := ""
	fmt.Println("Map and get Map Results from Slaves")
	//get map result from  slave1
	resp, err := http.Get("http://192.168.1.5:8091/chunk/baseCount")
	panicOnError(err)
	defer resp.Body.Close()
	// get stautus code
	fmt.Println("Status code:", resp.StatusCode)
	b, err := ioutil.ReadAll(resp.Body)
	result1_location := string(b)

	if strings.Contains(result1_location, "http://") {

		//get map result from  slave2
		resp, err := http.Get("http://192.168.1.5:8092/chunk/baseCount")
		panicOnError(err)
		defer resp.Body.Close()
		// get stautus code
		fmt.Println("Status code:", resp.StatusCode)
		b, err := ioutil.ReadAll(resp.Body)
		result2_location := string(b)

		if strings.Contains(result2_location, "http://") {

			//get map result from  slave3
			resp, err := http.Get("http://192.168.1.5:8093/chunk/baseCount")
			panicOnError(err)
			defer resp.Body.Close()
			// get stautus code
			fmt.Println("Status code:", resp.StatusCode)
			b, err := ioutil.ReadAll(resp.Body)
			result3_location := string(b)
			if strings.Contains(result3_location, "http://") {
				//get map result from  slave4
				resp, err := http.Get("http://192.168.1.5:8094/chunk/baseCount")
				panicOnError(err)
				defer resp.Body.Close()
				// get stautus code
				fmt.Println("Status code:", resp.StatusCode)
				b, err := ioutil.ReadAll(resp.Body)
				result4_location := string(b)
				if strings.Contains(result4_location, "http://") {
					mapResult := result1_location + "\n" + result2_location + "\n" + result3_location + "\n" + result4_location
					print(mapResult + "\n")
					send_mapResult_to_reducer(mapResult)
				}

			}

		}

	}

}

func send_mapResult_to_reducer(mapResult string) {

	b := strings.NewReader(mapResult)
	resp, err := http.Post(slave5.ipAddr+"/reduceResults", "text/plain", b)
	panicOnErrorM(err)
	defer resp.Body.Close()
	// get stautus code
	fmt.Println("Status code:", resp.StatusCode)
}

// master acts as client  above
// ------------------------------------main is here !! -------------------------------------------------------------
func main() {
	// divide_chunks_on_slaves()

	master_as_server()
}

//master acts server below

func master_as_server() {
	http.HandleFunc("/", indexM)
	http.HandleFunc("/fasta", get_Slave_ip)
	http.HandleFunc("/fasta/baseCount", get_finalReducer_ip)
	http.HandleFunc("/reducerResponse", handle_reducer_response)

	fmt.Println("starting server")
	http.ListenAndServe(":8090", nil)

}

func indexM(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Handling / req")
	fmt.Fprintf(w, "Hello from Master")
}
func get_Slave_ip(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Handling GET req")
	// http://localhost:8090/fasta?id=0
	query := req.URL.Query()
	id := query.Get("id")
	var result string
	var err error
	var idInt int

	if id != "" {
		idInt, err = strconv.Atoi(id)
		if idInt == 0 {
			result = "0\n" + slave1.ipAddr + "\n" + slave2.ipAddr + "\n" + slave3.ipAddr + "\n" + slave4.ipAddr
		} else if idInt == 1 {
			result = slave1.id + "\n" + slave1.ipAddr
		} else if idInt == 2 {
			result = slave2.id + "\n" + slave2.ipAddr
		} else if idInt == 3 {
			result = slave3.id + "\n" + slave3.ipAddr
		} else if idInt == 4 {
			result = slave4.id + "\n" + slave4.ipAddr
		} else {
			result = "Write id value from 0 to 4"
		}
	} else {
		result = "Write id value from 0 to 4 ,example:(http://localhost:8090/fasta?id=1)"
	}
	// if we had any error return status 500 and error
	if err != nil {
		errorHandlerM(w, req, http.StatusInternalServerError, err)
		return
	}
	// set header return data
	w.Header().Add("Content-Type", "text/plain")
	fmt.Fprintf(w, string(result))
}

func get_finalReducer_ip(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Handling GET req")
	map_request() //also map result is sent to reducer here if passed the constraints

	if reducer_Response != "" {
		fmt.Fprintf(w, "mapReduceResult\n"+reducer_Response)

	} else {
		fmt.Fprintf(w, "Reducer didn't respond")
	}

}

func handle_reducer_response(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Handling POST req")
	defer req.Body.Close()

	// read req body
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		errorHandlerM(w, req, http.StatusInternalServerError, err)
		return
	}
	reducer_Response = string(b)
	w.WriteHeader(http.StatusCreated)
}

func panicOnErrorM(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func errorHandlerM(w http.ResponseWriter, req *http.Request, status int, err error) {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "text/plain")
	fmt.Fprintf(w, `{error:%v}`, err.Error())
}

func panicOnError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
