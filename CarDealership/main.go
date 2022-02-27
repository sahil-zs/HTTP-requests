package main

import (
	"HTTP/CarDealership/handler"
	"log"
	"net/http"
)

//var uuid uuid2.NullUUID

func main() {
	//id := uuid.New()
	//fmt.Println(uuid.NewString())
	//fmt.Println("%T", uuid.NewString())
	http.HandleFunc("/car", handler.Update)

	log.Fatal(http.ListenAndServe(":8000", nil))

}
