package handler

import (
	"bytes"
	"encoding/json"
	uuid2 "github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestGetbyid this is test function for GetbyId()
func TestGetByiD(t *testing.T) {
	var UUID = "c6c2de38-7a30-46f4-8842-f474b2afa888"
	op1 := car{UUID, "model 1", 2019, "Tesla", "Diesel", Engine{3, 4, 100}}

	testCases := []struct {
		desc   string
		url    string
		op     car
		status int
	}{
		{"success", "localhost:8000/" + UUID, op1, http.StatusOK},
		{"invalid uuid", "localhost:8000/car/" + "abcde", car{}, http.StatusBadRequest},
		{"invalid uuid", "localhost:8000/car/" + "123445", car{}, http.StatusBadRequest},
	}
	for _, v := range testCases {
		req := httptest.NewRequest(http.MethodGet, v.url, nil)
		w := httptest.NewRecorder()
		GetbyId(w, req)

		d, err := ioutil.ReadAll(w.Body)
		if err != nil {
			return
		}

		var data car

		err = json.Unmarshal(d, &data)
		if err != nil {
			log.Println(err)
			return
		}
		assert.Equal(t, v.op, data)
	}
}

// TestGetbybrand this is test function for GetbyBrand()
func TestGetbybrand(t *testing.T) {
	output1 := []car{{uuid2.NewString(), "md2", 2020, "tesla", "electric", Engine{3, 0, 200}},
		{uuid2.NewString(), "md3", 2002, "tesla", "petrol", Engine{6, 3, 100}}}
	output2 := []car{{uuid2.NewString(), "md11", 2020, "tesla", "electric", Engine{}}, {uuid2.NewString(), "md3", 2002, "tesla", "petrol", Engine{}}}

	testCases := []struct {
		desc   string
		url    string
		op     []car
		status int
	}{
		{"Normal get by Brand", "localhost:8000/car/?Brand=tesla&engine=included", output1, http.StatusOK},
		{"passed empty Brand", "localhost:8000/car/?Brand=", []car{}, http.StatusBadRequest},
		{"passed invalid Brand", "localhost:8000/car/?Brand=fdgfg", []car{}, http.StatusBadRequest},
		{"Normal get by Brand without engine", "localhost:8000/car/?Brand=tesla", output2, http.StatusOK},
		{"passed empty engine", "localhost:8000/car/?Brand=ferrari&engine=", []car{}, http.StatusBadRequest},
	}
	for _, v := range testCases {
		req := httptest.NewRequest(http.MethodGet, "/car/Brand=", nil)
		w := httptest.NewRecorder()
		GetbyBrand(w, req)

		d, _ := ioutil.ReadAll(w.Body)

		var data []car

		err := json.Unmarshal(d, &data)
		if err != nil {
			log.Println(err)
		}
		assert.Equal(t, v.op, data)
	}
}

// TestPost this is test function fot post request
func TestPost(t *testing.T) {
	const UUID = "7718d22b-ab39-4fed-abcc-b2be39ca2ef0"

	carip := car{UUID, "model 1", 2019, "Tesla", "Diesel", Engine{3, 4, 100}}
	//emptyCar := car{}
	invalidIp := car{"7747", "model 11", 1999, "Ferrari", "Diesel", Engine{3, 4, 100}}
	invalidIdIp1 := car{"sdjfds", "model 11", 1999, "Ferrari", "Diesel", Engine{3, 4, 100}}

	testCases := []struct {
		desc   string
		ip     car
		op     car
		output int
	}{
		{"success", carip, carip, http.StatusCreated},
		{"empty entry", car{}, car{}, http.StatusBadRequest},
		{"duplicate Id", carip, carip, http.StatusBadRequest},
		{"invalid Id", invalidIp, car{}, http.StatusBadRequest},
		{"invalid Id", invalidIdIp1, car{}, http.StatusBadRequest},
	}
	for i, v := range testCases {
		body, err := json.Marshal(v.ip)
		if err != nil {
			t.Errorf("Marshall err,test:%v failed", i)
			continue
		}
		req := httptest.NewRequest(http.MethodPost, "/tc/", bytes.NewReader(body))
		w := httptest.NewRecorder()
		Create(w, req)
		res := w.Result()
		respBody, _ := ioutil.ReadAll(res.Body)

		assert.Equal(t, body, respBody)
	}
}

// TestPut this is test function for put request
func TestPut(t *testing.T) {
	const UUID = "7718d22b-ab39-4fed-abcc-b2be39ca2ef0"

	ip1 := car{UUID, "model 1", 2019, "Tesla", "Diesel", Engine{3, 4, 100}}
	ip2 := car{uuid2.NewString(), "model 1", 2019, "Tesla", "Diesel", Engine{3, 4, 100}}
	invalidIdIp1 := car{"7747", "model 11", 1999, "Ferrari", "Diesel", Engine{3, 4, 100}}
	invalidIdIp2 := car{"sdjfds", "model 11", 1999, "Ferrari", "Diesel", Engine{3, 4, 100}}

	testCases := []struct {
		desc   string
		input  car
		op     car
		status int
	}{
		{"success", ip1, ip1, http.StatusOK},
		{"record empty", car{}, car{}, http.StatusBadRequest},
		{"uuid valid but not present in database", ip2, ip2, http.StatusBadRequest},
		{"invalid uuid", invalidIdIp1, car{}, http.StatusCreated},
		{"invalid uuid", invalidIdIp2, car{}, http.StatusCreated},
	}
	for i, tc := range testCases {
		body, err := json.Marshal(tc.input)
		if err != nil {
			t.Errorf("Marshall err,test:%v failed", i)

			continue
		}
		req := httptest.NewRequest(http.MethodPut, "/car/"+tc.input.Id, bytes.NewReader(body))
		w := httptest.NewRecorder()
		Update(w, req)
		res := w.Result()
		respBody, _ := ioutil.ReadAll(res.Body)

		assert.Equal(t, body, respBody)
	}
}

// TestDelete this is test function for delete request
func TestDelete(t *testing.T) {
	var uuid = "aa146b8a-53af-40a9-a9b3-6edc46db6f78"
	testCases := []struct {
		desc   string
		id     string
		status int
	}{
		{"success", uuid, http.StatusNoContent},
		{"sent empty entry", "", http.StatusBadRequest},
		{"Invalid UUID", "fduygs", http.StatusBadRequest},
		{"Valid UUID not present in database", uuid2.NewString(), http.StatusInternalServerError},
	}
	for i, tc := range testCases {
		req := httptest.NewRequest(http.MethodPut, "/car/"+tc.id, nil)
		w := httptest.NewRecorder()
		Delete_(w, req)
		if w.Result().StatusCode != tc.status {
			t.Errorf("Test Failed %v. Expected %v GOT %v", i, tc.status, w.Result().StatusCode)
		}
	}
}
