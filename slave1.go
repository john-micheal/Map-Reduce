package main

import (
	// "encoding/json"
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func main() {

	http.HandleFunc("/", index)
	http.HandleFunc("/fasta", handleFasta)
	http.HandleFunc("/chunk/baseCount", handle_bases_count)
	http.HandleFunc("/countResult", handle_count_Result)
	fmt.Println("starting server")
	http.ListenAndServe(":8091", nil)
}

func index(w http.ResponseWriter, req *http.Request) {
	file_bytes, err := ioutil.ReadFile("slave1.fasta")
	fmt.Println("Handling / req")
	fmt.Fprintf(w, string(file_bytes))
	if err != nil {
		errorHandler(w, req, http.StatusInternalServerError, err)
		return
	}
}

func handleFasta(w http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		// fmt.Fprintf(w,"Nothing to Get yet")
		get(w, req)
	} else if req.Method == "POST" {
		post(w, req)
	} else {
		fmt.Println("Handling invalid /users req")
		errorHandler(w, req, http.StatusMethodNotAllowed, fmt.Errorf("Invalid Method"))
	}
}

func get(w http.ResponseWriter, req *http.Request) {
	file_bytes, err := ioutil.ReadFile("slave1.fasta")
	fmt.Println("Handling / req")
	fmt.Fprintf(w, string(file_bytes))
	if err != nil {
		errorHandler(w, req, http.StatusInternalServerError, err)
		return
	}

}

func post(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Handling POST req")
	defer req.Body.Close()

	// read req body
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		errorHandler(w, req, http.StatusInternalServerError, err)
		return
	}

	err = ioutil.WriteFile("slave1.fasta", b, 0644)
	if err != nil {
		log.Fatal(err)
	}

	w.WriteHeader(http.StatusCreated)
}
func handle_bases_count(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Handling / req")

	Nucleobases := map[string]int{"A": 0, "T": 0, "C": 0, "G": 0}

	// read file line by line
	f, err := os.Open("slave1.fasta")
	panicOnError(err)
	defer f.Close()

	s := bufio.NewScanner(f)
	itr := 0
	for s.Scan() {
		// fmt.Println(s.Text())
		if itr != 0 {
			Nucleobases["A"] += strings.Count(s.Text(), "A")
			Nucleobases["T"] += strings.Count(s.Text(), "T")
			Nucleobases["C"] += strings.Count(s.Text(), "C")
			Nucleobases["G"] += strings.Count(s.Text(), "G")
		}
		itr = 1
	}
	err = s.Err()
	panicOnError(err)

	countResult := "A:" + strconv.Itoa(Nucleobases["A"]) + "\nT:" + strconv.Itoa(Nucleobases["T"]) +
		"\nC:" + strconv.Itoa(Nucleobases["C"]) + "\nG:" + strconv.Itoa(Nucleobases["G"])

	f, err = os.Create("countResult1.txt")
	panicOnError(err)
	defer f.Close()
	f.WriteString(countResult)

	my_count_Result_location := "http://192.168.1.5:8091/countResult"
	fmt.Fprintf(w, my_count_Result_location)
	if err != nil {
		errorHandler(w, req, http.StatusInternalServerError, err)
		return
	}
}
func handle_count_Result(w http.ResponseWriter, req *http.Request) {
	file_bytes, err := ioutil.ReadFile("countResult1.txt")
	fmt.Println("Handling / req")
	fmt.Fprintf(w, string(file_bytes))
	if err != nil {
		errorHandler(w, req, http.StatusInternalServerError, err)
		return
	}
}

func errorHandler(w http.ResponseWriter, req *http.Request, status int, err error) {
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
