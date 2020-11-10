package encdec_test

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/ehsoc/rest/encdec"
)

func TestJSONEncoder(t *testing.T) {
	car := Car{Brand: "Fiat"}
	gotCar := Car{}
	encoder := encdec.JSONEncoder{}
	buf := bytes.NewBuffer([]byte(""))
	encoder.Encode(buf, car)
	json.NewDecoder(buf).Decode(&gotCar)

	if !reflect.DeepEqual(gotCar, car) {
		t.Errorf("got:%v want:%v", gotCar, car)
	}
}
