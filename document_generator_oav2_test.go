package resource_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/go-openapi/spec"
)

type Pet struct {
	// id
	ID int64 `json:"id,omitempty"`
	// name
	// Required: true
	Name string `json:"name"`
	// photo urls
	// Required: true
	PhotoUrls []string `json:"photoUrls" xml:"photoUrl"`
	// pet status in the store
	// Enum: [available pending sold]
	Status string `json:"status,omitempty"`
	// tags
	Tags []Tag `json:"tags" xml:"tag"`
}

// Tag tag
//
// swagger:model Tag
type Tag struct {
	// id
	ID int64 `json:"id,omitempty"`
	// name
	Name string `json:"name,omitempty"`
}

var petSchemaJson = `{
	"type": "object",
	"description" : "A Pet object.",
	"properties": {
	  "id": {
		"type": "integer",
		"format": "int64"
	  },
	  "name": {
		"type": "string"
	  },
	  "photoUrls": {
		"type": "array",
		"items": {
		  "type": "string"
		}
	  },
	  "tags": {
		"type": "array",
		"items":{
			"type": "object",
			"description": "A Tag object.",
			"properties": {
			  "id": {
				"type": "integer",
				"format": "int64"
			  },
			  "name": {
				"type": "string"
			  }
			}
		  }
	  },
	  "status": {
		"type": "string"
	  }
	}
  }`

// func TestToSchema(t *testing.T) {
// 	testCases := []struct {
// 		v    interface{}
// 		json string
// 	}{
// 		{"", `{"type":"string"}`},
// 		{[]string{}, `{"type":"array","items":{"type":"string"}}`},
// 		{Pet{}, petSchemaJson},
// 	}
// 	for _, test := range testCases {
// 		t.Run(fmt.Sprintf("%T", test.v), func(t *testing.T) {
// 			got := resource.ToSchema(test.v)
// 			gotJson, _ := got.MarshalJSON()
// 			assertJsonSchemaEqual(t, string(gotJson), test.json)
// 		})
// 	}
// }

func assertJsonSchemaEqual(t *testing.T, got, want string) {
	gotJson := spec.Schema{}
	err := json.Unmarshal([]byte(got), &gotJson)
	if err != nil {
		t.Fatalf("Not expecting error: %v", err)
	}
	wantJson := spec.Schema{}
	err = json.Unmarshal([]byte(want), &wantJson)
	if err != nil {
		t.Fatalf("Not expecting error: %v", err)
	}
	if !reflect.DeepEqual(gotJson, wantJson) {
		t.Errorf("\ngot: %v \nwant: %v", got, want)
	}
}
