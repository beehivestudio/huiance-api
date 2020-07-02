package utils

import "testing"

func TestUnion(t *testing.T) {
	slice1 := []string{"9527", "9313"}

	slice2 := []string{"9898", "9527"}

	slice3 := Union(slice1, slice2)

	if len(slice3) != 3 {
		t.Error("Union failed, msg:", slice3)
	}

	t.Log(slice3)
}

func TestIntersect(t *testing.T) {
	slice1 := []string{"9527", "9313"}

	slice2 := []string{"9898", "9527"}

	slice3 := Intersect(slice1, slice2)

	if len(slice3) != 1 {
		t.Error("Intersect failed, msg:", slice3)
	}

	t.Log(slice3)
}
