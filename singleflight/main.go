package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/sync/singleflight"
)

func main() {
	var requestGroup singleflight.Group
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("calling the endpoint")
		response, err := http.Get("https://jsonplaceholder.typicode.com/photos")

		if err != nil {
			fmt.Print(err.Error())
			os.Exit(1)
		}

		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(2 * time.Second)
		w.Write(responseData)
	})
	http.HandleFunc("/singleflight", func(w http.ResponseWriter, r *http.Request) {
		res, err, shared := requestGroup.Do("singleflight", func() (interface{}, error) {
			fmt.Println("calling the endpoint")
			response, err := http.Get("https://jsonplaceholder.typicode.com/photos")

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return nil, err
			}
			responseData, err := ioutil.ReadAll(response.Body)
			if err != nil {
				log.Fatal(err)
			}
			time.Sleep(2 * time.Second)
			return string(responseData), err
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		result := res.(string)
		fmt.Println("shared = ", shared)
		fmt.Fprintf(w, "%q", result)
	})

	http.ListenAndServe(":3000", nil)
}
