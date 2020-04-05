package main

import (
	"time"
	"fmt"
	"log"
	"net/http"
	"os"
	"io/ioutil"
	"encoding/json"
	"sort"
)
type Vehicle struct {
	Timedate time.Time
	Number string
	Speed float64
}
type Req struct {
	Date string
	Speed float64
} 
var v = make ([]Vehicle, 10000)

func handler(w http.ResponseWriter, r *http.Request){ 			
	
	err := r.ParseForm()

	if err != nil {
		fmt.Println("Error parsing form!")
	}

	var rq = Req {}
	err = json.NewDecoder(r.Body).Decode(&rq)

	if err != nil {
		fmt.Println("Error decoding file!")
	}
	
	
	if reqProcesser(rq) == nil {
		rqo, _ := json.Marshal("Nothing found!")
		fmt.Println("Nothing found!")
		w.Write(rqo)
	}else{
		rqo, _ := json.Marshal(reqProcesser(rq)) 
		fmt.Println("The result: ", reqProcesser(rq))
		w.Write(rqo)
	}
}

func reqProcesser(rq Req) (result []Vehicle) {
	
	var m []Vehicle							//slice containing cars passed in one day
	for _, l := range v {
		if l.Timedate.Format("2006-01-02") == rq.Date {	
			m = append(m,l)	
		}
	}

	sort.SliceStable(m, func(i, j int) bool {return m[i].Speed < m[j].Speed}) //sorts the slice m by speed 
											//and panics if interface is mot a slice
	if rq.Speed == -1 {						//if the request contains only the date, min-max case
		result = make([]Vehicle, 2)
		result[0] = m[0]					//result slice contains min and max speed rows
		result[1] = m[len(m)-1]				//aka the first and the last elements of sorted m
	}else{ 									//moreThan case
		var result []Vehicle 
		for _, l := range m {
			if l.Speed > rq.Speed {		//compares all elements of m with rq
				result = append(result, l)		//if m.Speed is more than rq.speed => writes in result
			}	
		}									//else moreThan case
	}
	return result
}


func reciever(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm() 

	if err != nil {
		fmt.Println("error:", err)
	}

	_, err = os.Create("cars.json")
	err = json.NewDecoder(r.Body).Decode(&v)

	if err != nil {
		fmt.Println("error:", err)
	}
	out, err := json.MarshalIndent(v, " ", "")
	err = ioutil.WriteFile("cars.json", out, 0644)

	if err != nil {
		fmt.Println("error:", err)
	}
}

func main() {
	
	http.HandleFunc("/data", handler)
	http.HandleFunc("/recieve", reciever)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
