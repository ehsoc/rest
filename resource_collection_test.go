package rest_test

import (
	"reflect"
	"sort"
	"testing"

	"github.com/ehsoc/rest"
)

func TestGetResources(t *testing.T) {
	r := rest.ResourceCollection{}
	r.Resource("car", func(r *rest.Resource) {
		r.Resource("fiat", nil)
		r.Resource("citroen", nil)
		r.Resource("ford", nil)
	})
	rootRs := r.Resources()
	rs := rootRs[0].Resources()
	sort.Slice(rs, func(i, j int) bool {
		return rs[i].Path() < rs[j].Path()
	})
	if len(rs) != 3 {
		t.Errorf("got: %v want: %v", len(rs), 3)
	}
	assertStringEqual(t, rs[0].Path(), "citroen")
	assertStringEqual(t, rs[1].Path(), "fiat")
	assertStringEqual(t, rs[2].Path(), "ford")
}

func TestResource(t *testing.T) {
	collection := rest.ResourceCollection{}
	collection.Resource("find", func(r *rest.Resource) {
		r.Resource("left", func(r *rest.Resource) {
		})
		r.Resource("right", func(r *rest.Resource) {
		})
	})
	findNode := collection.Resources()[0]
	if findNode.Path() != "find" {
		t.Errorf("got : %v want: %v", findNode.Path(), "find")
	}
	if len(findNode.Resources()) != 2 {
		t.Fatalf("expecting 2 sub nodes got: %v", len(findNode.Resources()))
	}
	directionResources := findNode.Resources()
	sort.Slice(directionResources, func(i, j int) bool {
		return directionResources[i].Path() < directionResources[j].Path()
	})
	if directionResources[0].Path() != "left" {
		t.Errorf("got : %v want: %v", findNode.Resources()[0].Path(), "left")
	}
	if directionResources[1].Path() != "right" {
		t.Errorf("got : %v want: %v", findNode.Resources()[1].Path(), "right")
	}
}

func TestUriParamResource(t *testing.T) {
	collection := rest.ResourceCollection{}

	carIDParam := rest.NewURIParameter("carId", reflect.String)

	mo := rest.NewMethodOperation(&OperationStub{}, rest.NewResponse(200))

	collection.ResourceP(carIDParam, func(r *rest.Resource) {
		r.Resource("left", func(r *rest.Resource) {
			r.Get(mo, mustGetJSONContentType()).AddParameter(carIDParam)
		})
		r.Resource("right", func(r *rest.Resource) {
		})
	})

	findNode := collection.Resources()[0]
	if findNode.Path() != "{carId}" {
		t.Errorf("got : %v want: %v", findNode.Path(), "{carId}")
	}
	if len(findNode.Resources()) != 2 {
		t.Fatalf("expecting 2 sub nodes got: %v", len(findNode.Resources()))
	}
	directionResources := findNode.Resources()
	sort.Slice(directionResources, func(i, j int) bool {
		return directionResources[i].Path() < directionResources[j].Path()
	})
	if directionResources[0].Path() != "left" {
		t.Errorf("got : %v want: %v", findNode.Resources()[0].Path(), "left")
	}
	if directionResources[1].Path() != "right" {
		t.Errorf("got : %v want: %v", findNode.Resources()[1].Path(), "right")
	}
}
