package exporter

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)


var C chan interface{}
var exitSignal chan bool

func DumpData(fileName string){
	fmt.Println("Started dump routine")
	C = make(chan interface{}, 20)
	exitSignal = make(chan bool)
	timer := time.NewTimer(2 * time.Second)
	var items []interface{}
outerloop:
	for true {
		select {
		case itemObject := <- C:{
			items = append(items, itemObject)
		}
		case <- timer.C:
			timer.Reset(2 * time.Second)
			updateJson(&items, fileName)
			fmt.Println("Dumping after 2 sec.")
			fmt.Println("Done dumping")

		case <- exitSignal:
			updateJson(&items, fileName)
			fmt.Println("Dumping at end.")
			fmt.Println("Done dumping for real this time.")
			break outerloop
		}
	}
	exitSignal <- true

}

func updateJson(items *[]interface{}, fileName string){
	_, err := os.Stat(fileName)
	if err != nil {
		result, _ := json.Marshal(*items)
		ioutil.WriteFile(fileName, result, 0777)
		return
	}

	jsonText, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println(err)
	}

	var idents []interface{}
	err = json.Unmarshal([]byte(jsonText), &idents)
	if err != nil {
		fmt.Println(err)
	}

	idents = append(idents, *items...)
	*items = nil
	result, _ := json.Marshal(&idents)
	ioutil.WriteFile(fileName, result, 0777)
}

func Wait(){
	exitSignal <- true
	select {
		case <-exitSignal:
			fmt.Println("Succesfully imported data.")
	}
}
