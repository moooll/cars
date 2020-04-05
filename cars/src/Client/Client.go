package main

import (
	"net/http"
	"io/ioutil"
	"time"
	"math/rand"
	"fmt"
	"strconv"
	"log"
	"encoding/json"
	"bytes"
	"os"
	"bufio"
	"strings"
)

type Vehicle struct {
	Timedate time.Time
	Number string
	Speed float64
}

var cars = make ([]Vehicle, 10000)

func makeCars(cars []Vehicle){
	maxCars := rand.Intn(300)
	
	for i := range cars {
		cars[i].Number = strconv.Itoa(rand.Intn(9999-1111)+1111)+" "+"PP-"+strconv.Itoa(rand.Intn(7-1)+1)  	//РР тоже должно задаваться рандомом + 
		cars[i].Speed = float64(rand.Int63n(200-20))+float64(maxCars/6)	//super pseudo rand func
		if i == 0 {
			cars[i].Timedate = time.Date(2020, 3, 27, 0, 0, 0, 0, time.UTC)
		}else if  i != 0 && i % maxCars == 0 {			
			cars[i].Timedate = cars[i-1].Timedate.AddDate(0, 0, 1)
		}else{
			cars[i].Timedate = cars[i-1].Timedate.Add(time.Second * time.Duration(rand.Intn(86400))) // randomly changes time within 24 hours
		}
	}
}

func makeReq() {
	type req struct {
		Date string
		Speed float64
	}

	fmt.Println("1. Learn the info about range speed for a day\n2. Info about all cars with speed over the set value\n\nPlease, make up your mind")
	reader := bufio.NewReader(os.Stdin)
	ch, _ := reader.ReadString('\n')

	if strings.Compare(ch,"1\n")==0 { 

		fmt.Println("Enter the date (yyyy-mm-dd)") 			
		date,_ := reader.ReadString('\n')  			
		date = strings.Replace(date, "\n","", -1)	
		var s float64 = -1			
		r := req {date, s}
		b, err := json.Marshal(r)

		if err != nil {
			fmt.Println("Error marshalling")
		}

		resp, err := http.Post("http://localhost:8080/data", "application/json", bytes.NewBuffer(b)) 	
		
		if err != nil {
			fmt.Println("Error sending request")
		}

		body, err := ioutil.ReadAll(resp.Body)
		
		if err != nil {
			fmt.Println("Error reading from request")
		}

		err = ioutil.WriteFile("rspns.json", body, 0644)

		if err != nil {
			fmt.Println("Error writing to file")
		}
		
	}else if string(ch) == "2\n" {
		
		fmt.Println("Enter the date (yyyy-mm-dd)") 			
		date,_ := reader.ReadString('\n')  			
		date = strings.Replace(date, "\n","", -1)	
		fmt.Println("Enter the speed") 			
		speed,_ := reader.ReadString('\n')  			
		speed = strings.Replace(date, "\n","", -1)	
		s, _ := strconv.ParseFloat(speed, 64) 
		r := req {date, s}
		b, err := json.Marshal(r)

		if err != nil {
			fmt.Println("Error marshalling ", err)
		}

		resp, err := http.Post("http://localhost:8080/data", "application/json", bytes.NewBuffer(b)) 	
	
		if err != nil {
			fmt.Println("Error sending request")
		}

		rs, err := ioutil.ReadAll(resp.Body)
		
		if err != nil {
			fmt.Println("Error unmarshalling")
		}

		defer resp.Body.Close()
		
		err = ioutil.WriteFile("rspns.json", rs, 0644)

		if err != nil {
			fmt.Println("Error writing in file!")
		}
	}else{
		fmt.Println ("Something's wrong! Check the input, please.'1' and '2' values permitted only.")
	}	
}

func main() {

	makeCars(cars)
	client := &http.Client {}
	data, err := json.Marshal(cars)

	if err != nil {
		log.Fatal(err)
	}

	_, err = client.Post("http://localhost:8080/recieve", "application/json", bytes.NewBuffer(data))  
	
	if err != nil {
		log.Fatal(err)
		fmt.Println("error: ", err)
	}
	makeReq()
}








