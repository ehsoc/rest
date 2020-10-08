package simple_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/ehsoc/resource/_example/simple"
)

func TestSimpleServer(t *testing.T) {
	t.Run("get car", func(t *testing.T) {
		carId := "101"
		s := simple.GenerateServer()
		testServer := httptest.NewServer(s)
		cli := http.Client{}
		res, _ := cli.Get(testServer.URL + "/v1/car/" + carId)
		if res.StatusCode != 200 {
			t.Errorf("got: %d want:%d", res.StatusCode, 200)
		}
		wantCar := simple.Car{carId, "Foo"}
		gotCar := simple.Car{}
		json.NewDecoder(res.Body).Decode(&gotCar)
		if !reflect.DeepEqual(gotCar, wantCar) {
			t.Errorf("got: %v want: %v", gotCar, wantCar)
		}
	})
	t.Run("car not found", func(t *testing.T) {
		carId := "2"
		s := simple.GenerateServer()
		testServer := httptest.NewServer(s)
		cli := http.Client{}
		res, _ := cli.Get(testServer.URL + "/v1/car/" + carId)
		if res.StatusCode != 404 {
			t.Errorf("got: %d want:%d", res.StatusCode, 404)
		}
		wantCar := simple.Car{carId, "Foo"}
		gotCar := simple.Car{}
		json.NewDecoder(res.Body).Decode(&gotCar)
		if reflect.DeepEqual(gotCar, wantCar) {
			t.Errorf("got: %v want: %v", gotCar, wantCar)
		}
	})
	t.Run("car not found", func(t *testing.T) {
		carId := "error"
		s := simple.GenerateServer()
		testServer := httptest.NewServer(s)
		cli := http.Client{}
		res, _ := cli.Get(testServer.URL + "/v1/car/" + carId)
		if res.StatusCode != 500 {
			t.Errorf("got: %d want:%d", res.StatusCode, 500)
		}
		wantCar := simple.Car{carId, "Foo"}
		gotCar := simple.Car{}
		json.NewDecoder(res.Body).Decode(&gotCar)
		if reflect.DeepEqual(gotCar, wantCar) {
			t.Errorf("got: %v want: %v", gotCar, wantCar)
		}
	})
}
