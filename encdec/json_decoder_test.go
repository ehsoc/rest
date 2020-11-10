package encdec_test

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/ehsoc/resource/encdec"
)

func TestJSONDecoder(t *testing.T) {
	car := Car{Brand: "Fiat"}
	gotCar := Car{}
	buf := bytes.NewBuffer([]byte(""))
	json.NewEncoder(buf).Encode(car)

	decoder := encdec.JSONDecoder{}

	decoder.Decode(buf, &gotCar)

	if !reflect.DeepEqual(gotCar, car) {
		t.Errorf("got:%v want:%v", gotCar, car)
	}
}
