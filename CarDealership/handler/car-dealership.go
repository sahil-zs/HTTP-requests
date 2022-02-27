package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type car struct {
	Id       string `json:"Id"`
	Name     string `json:"name"`
	Year     int    `json:"Year"`
	Brand    string `json:"Brand"`
	FuelType string `json:"fuel_type"`
	Engine   Engine `json:"engine"`
}

type Engine struct {
	displacement    int `json:"displacement"`
	no_of_cylinders int `json:"no_of_cylinders"`
	range_of_car    int `json:"range_of_car"`
}

// connect function to get a database connection
func connect() *sql.DB {
	db, err := sql.Open("mysql", "sahil:password@/CarDealership")
	if err != nil {
		fmt.Println(err.Error())
	}
	return db
}

func GetbyId(w http.ResponseWriter, request *http.Request) {
	if request.Method == "GET" {
		id := request.URL.Query().Get("id")

		if id == "" {
			w.WriteHeader(http.StatusNotFound)

		} else {
			db := connect()
			er := db.QueryRow("select * from car where id=?; ", id)
			if er != nil {
				fmt.Println(er)

			}

			var c car
			var eng string
			switch err := er.Scan(&c.Id, &c.Name, &c.Year, &c.Brand, &c.FuelType, &eng); err {
			case sql.ErrNoRows:
				fmt.Println("error" + err.Error())
			case nil:
				_, err := fmt.Fprintln(w, c)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
				}
			default:
				panic(err)
			}
		}
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func GetbyBrand(w http.ResponseWriter, request *http.Request) {
	if request.Method == "GET" {
		brand := request.URL.Query().Get("brand")

		if brand == "" {
			w.WriteHeader(http.StatusNotFound)
		} else {
			db := connect()

			rows, er := db.Query("SELECT * FROM car WHERE brand = ?;", brand)
			if er != nil {
				fmt.Println(er.Error())
			}

			var cars []car

			for rows.Next() {
				var c car
				var eng string

				err := rows.Scan(&c.Id, &c.Name, &c.Year, &c.Brand, &c.FuelType, &eng)
				if err != nil {
					fmt.Println(err.Error())
				}
				cars = append(cars, c)
			}
			_, err := fmt.Fprintln(w, cars)
			if err != nil {
				fmt.Println(err.Error())
			}
		}
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)

}

func Create(w http.ResponseWriter, request *http.Request) {

	//fmt.Println("at beg")

	body, err := ioutil.ReadAll(request.Body)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	var c car
	err = json.Unmarshal(body, &c)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(err)
	}
	fmt.Println(c)

	if c.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
	}

	Year := time.Now().Year()
	if c.Year < 1980 || c.Year > Year {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	brands := [5]string{"Tesla", "Porsche", "Ferrari", "Mercedes", "BMW"}
	flag := false
	for i := range brands {
		if brands[i] == c.Brand {
			flag = true
			break
		}
	}
	if !flag || c.Brand == "" {
		w.WriteHeader(http.StatusBadRequest)
	}

	fuelType := [3]string{"Petrol", "Diesel", "Electric"}
	flag = false
	for i := range fuelType {
		if fuelType[i] == c.FuelType {
			flag = true
			break
		}
	}
	if !flag || c.FuelType == "" {
		w.WriteHeader(http.StatusBadRequest)
	}

	c.Id = uuid.NewString()

	db := connect()
	insert, er := db.Query("INSERT INTO car  VALUES (?, ?,?,?,?,?)",
		c.Id, c.Name, c.Year, c.Brand, c.FuelType, "12gg")

	if er != nil {
		fmt.Println(er.Error())
	}

	defer func(insert *sql.Rows) {
		err := insert.Close()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
	}(insert)

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
	}(db)
	fmt.Println(c)

	w.WriteHeader(http.StatusCreated)

}

func Update(w http.ResponseWriter, request *http.Request) {
	body, err := io.ReadAll(request.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	var c car
	eng := "engtest"

	err = json.Unmarshal(body, &c)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(err)
	}

	//ValidateCarCredentials(c, w)

	db := connect()

	_, er := db.Exec("update car set name = ?, year = ?, brand = ? ,FuelType = ?, Engine =? where id = ?;",
		c.Name, c.Year, c.Brand, c.FuelType, eng, c.Id)

	if er != nil {
		fmt.Println(er.Error())
	}

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
	}(db)

	w.WriteHeader(http.StatusCreated)

}

func Delete_(w http.ResponseWriter, request *http.Request) {
	db := connect()
	id := request.URL.Query().Get("id")

	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	delfrom, err := db.Prepare("DELETE FROM car WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	delfrom.Exec(id)
	w.WriteHeader(http.StatusOK)
	log.Println("DELETE")
	defer db.Close()
}
