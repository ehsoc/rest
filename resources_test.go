package resource_test

import (
	"reflect"
	"sort"
	"testing"

	"github.com/ehsoc/resource"
)

func TestGetResources(t *testing.T) {
	r := resource.Resources{}
	r.Resource("car", func(r *resource.Resource) {
		r.Resource("fiat", nil)
		r.Resource("citroen", nil)
		r.Resource("ford", nil)
	})
	rootRs := r.GetResources()
	rs := rootRs[0].GetResources()
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
	collection := resource.Resources{}
	collection.Resource("find", func(r *resource.Resource) {
		r.Resource("left", func(r *resource.Resource) {
		})
		r.Resource("right", func(r *resource.Resource) {
		})
	})
	findNode := collection.GetResources()[0]
	if findNode.Path() != "find" {
		t.Errorf("got : %v want: %v", findNode.Path(), "find")
	}
	if len(findNode.GetResources()) != 2 {
		t.Fatalf("expecting 2 sub nodes got: %v", len(findNode.GetResources()))
	}
	directionResources := findNode.GetResources()
	sort.Slice(directionResources, func(i, j int) bool {
		return directionResources[i].Path() < directionResources[j].Path()
	})
	if directionResources[0].Path() != "left" {
		t.Errorf("got : %v want: %v", findNode.GetResources()[0].Path(), "left")
	}
	if directionResources[1].Path() != "right" {
		t.Errorf("got : %v want: %v", findNode.GetResources()[1].Path(), "right")
	}
}
func TestAddResource(t *testing.T) {
	collection := resource.Resources{}
	findNode := resource.NewResource("find")
	leftNode := resource.NewResource("left")
	rightNode := resource.NewResource("right")
	findNode.AddResource(rightNode)
	findNode.AddResource(leftNode)
	collection.AddResource(findNode)

	gotFindNode := collection.GetResources()[0]
	if !reflect.DeepEqual(gotFindNode, findNode) {
		t.Errorf("got : %v \nwant: %v", gotFindNode, findNode)
	}
	if len(findNode.GetResources()) != 2 {
		t.Fatalf("expecting 2 sub nodes got: %v", len(findNode.GetResources()))
	}
	directionResources := findNode.GetResources()
	sort.Slice(directionResources, func(i, j int) bool {
		return directionResources[i].Path() < directionResources[j].Path()
	})
	if directionResources[0].Path() != "left" {
		t.Errorf("got : %v want: %v", findNode.GetResources()[0].Path(), "left")
	}
	if directionResources[1].Path() != "right" {
		t.Errorf("got : %v want: %v", findNode.GetResources()[1].Path(), "right")
	}
}
