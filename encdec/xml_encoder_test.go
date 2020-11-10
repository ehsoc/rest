package encdec_test

import (
	"bytes"
	"encoding/xml"
	"reflect"
	"testing"

	"github.com/ehsoc/resource/encdec"
)

func TestXMLEncoder(t *testing.T) {
	car := Car{Brand: "Fiat"}
	gotCar := Car{}
	encoder := encdec.XMLEncoder{}
	buf := bytes.NewBuffer([]byte(""))
	encoder.Encode(buf, car)
	xml.NewDecoder(buf).Decode(&gotCar)

	if !reflect.DeepEqual(gotCar, car) {
		t.Errorf("got:%v want:%v", gotCar, car)
	}
}
