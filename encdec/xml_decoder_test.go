package encdec_test

import (
	"bytes"
	"encoding/xml"
	"reflect"
	"testing"

	"github.com/ehsoc/resource/encdec"
)

func TestXMLDecoder(t *testing.T) {
	car := Car{Brand: "Fiat"}
	gotCar := Car{}
	buf := bytes.NewBuffer([]byte(""))
	xml.NewEncoder(buf).Encode(car)

	decoder := encdec.XMLDecoder{}

	decoder.Decode(buf, &gotCar)

	if !reflect.DeepEqual(gotCar, car) {
		t.Errorf("got:%v want:%v", gotCar, car)
	}
}
