package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"gopkg.in/yaml.v2"
)

type Customer struct {
	Id         int    `json:"id"`
	Email      string `json:"email"`
	First      string `json:"first"`
	Last       string `json:"last"`
	Company    string `json:"company"`
	Created_at time.Time
	Country    string `json:"country"`
}

//readfuile

func readFile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if _, err := os.Stat("customer.json"); err == nil {
		file, err := os.Open("customer.json")
		if err != nil {
			log.Fatal(err)
		}
		b, _ := ioutil.ReadAll(file)
		var customers []Customer
		json.Unmarshal([]byte(string(b)), &customers)
		json.NewEncoder(w).Encode(customers)
	} else {
		fmt.Fprintln(w, "fileNotExist")
	}
}

//create file

func createFile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if _, err := os.Stat("customer.json"); err == nil {
		file, err := os.Open("customer.json")
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		b, _ := ioutil.ReadAll(file)
		var customers []Customer
		json.Unmarshal([]byte(string(b)), &customers)
		var customer Customer
		_ = json.NewDecoder(r.Body).Decode(&customer)
		customer.Created_at = time.Now()
		customer.Id = len(customers) + 1
		customers = append(customers, customer)
		jsonData, _ := json.Marshal(customers)
		err = ioutil.WriteFile("customer.json", jsonData, 0644)
		if err != nil {
			panic(err)
		}
		json.NewEncoder(w).Encode(customer)
	} else {
		os.Create("customer.json")
		file, err := os.Open("customer.json")
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		b, _ := ioutil.ReadAll(file)
		var customers []Customer
		json.Unmarshal([]byte(string(b)), &customers)
		var customer Customer
		_ = json.NewDecoder(r.Body).Decode(&customer)
		customer.Created_at = time.Now()
		customers = append(customers, customer)
		jsonData, _ := json.Marshal(customers)
		err = ioutil.WriteFile("customer.json", jsonData, 0644)
		if err != nil {
			panic(err)
		}
		json.NewEncoder(w).Encode(customer)
	}

}

//modify and save the file

func modifyAndSaveFile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])
	if _, err := os.Stat("customer.json"); err == nil {
		file, err := os.Open("customer.json")
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		b, _ := ioutil.ReadAll(file)
		var customers []Customer
		json.Unmarshal([]byte(string(b)), &customers)
		for index, item := range customers {
			if item.Id == id {
				// fmt.Fprint(w, "id matched"+params["id"])
				customers = append(customers[:index], customers[index+1:]...)
				var customer Customer
				_ = json.NewDecoder(r.Body).Decode(&customer)
				customer.Created_at = time.Now()
				customer.Id = item.Id
				customers = append(customers, customer)
				jsonData, _ := json.Marshal(customers)
				err = ioutil.WriteFile("customer.json", jsonData, 0644)
				// fmt.Fprint(w, "record modified")
				json.NewEncoder(w).Encode(customer)
				if err != nil {
					panic(err)
				}
				break

			}
		}
	} else {
		fmt.Fprint(w, "fileNotExist")
	}

}

//pagination
func pagenateAndReadFile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	page, _ := strconv.Atoi(params["page"])
	if _, err := os.Stat("customer.json"); err == nil {
		file, err := os.Open("customer.json")
		if err != nil {
			log.Fatal(err)
		}
		b, _ := ioutil.ReadAll(file)
		var customers []Customer
		json.Unmarshal([]byte(string(b)), &customers)
		start := (page - 1) * 20
		stop := start + 20
		if start > len(customers) {
			fmt.Fprintf(w, "nodata")
			return
		}

		if stop > len(customers) {
			stop = len(customers)
		}
		json.NewEncoder(w).Encode(customers[start:stop])

	} else {
		fmt.Fprintln(w, "fileNotExist")
	}
}

//convert to yaml

func convertToYaml(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/yaml")
	if _, err := os.Stat("customer.json"); err == nil {
		file, err := os.Open("customer.json")
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		b, _ := ioutil.ReadAll(file)
		var customers []Customer
		json.Unmarshal([]byte(string(b)), &customers)
		yamlData, _ := yaml.Marshal(&customers)
		err = ioutil.WriteFile("customer.yaml", yamlData, 0644)
		if err != nil {
			panic(err)
		}
		fmt.Fprint(w, string(yamlData))
	} else {
		fmt.Fprint(w, "fileNotExist")

	}

}

//convert to xml

func convertToXml(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/xml")
	if _, err := os.Stat("customer.json"); err == nil {
		file, err := os.Open("customer.json")
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		b, _ := ioutil.ReadAll(file)
		var customers []Customer
		json.Unmarshal([]byte(string(b)), &customers)
		xmlData, _ := xml.Marshal(&customers)
		err = ioutil.WriteFile("customer.xml", xmlData, 0644)
		if err != nil {
			panic(err)
		}
		fmt.Fprint(w, string(xmlData))
	} else {
		fmt.Fprint(w, "fileNotExist")

	}

}

//delete specified file

func deleteFile(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	if _, err := os.Stat(params["filename"]); err == nil {
		err := os.Remove(params["filename"])
		if err != nil {
			log.Fatal(err)
		} else {
			fmt.Fprint(w, "Deleted: "+params["filename"])
		}
	} else {
		fmt.Fprint(w, "fileNotExist")
	}
}

// handleRequests
func handleRequests() {

	myRouter := mux.NewRouter()
	myRouter.HandleFunc("/api/customers", readFile).Methods("GET")
	myRouter.HandleFunc("/api/customers", createFile).Methods("POST")
	myRouter.HandleFunc("/api/customers/{id}", modifyAndSaveFile).Methods("POST")
	myRouter.HandleFunc("/api/customers/toxml", convertToXml).Methods("GET")
	myRouter.HandleFunc("/api/customers/toyaml", convertToYaml).Methods("GET")
	myRouter.HandleFunc("/api/delete/{filename}", deleteFile).Methods("POST")
	myRouter.HandleFunc("/api/customers/page/{page}", pagenateAndReadFile).Methods("GET")
	handler := cors.Default().Handler(myRouter)
	fmt.Println("Starting server on port :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))

}

func main() {
	handleRequests()
}
