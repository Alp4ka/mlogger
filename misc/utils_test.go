package misc

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCoalesce_ShouldReturnNil_WhenAllNil(t *testing.T) {
	type A struct {
		Name string
	}

	ret := Coalesce((*A)(nil), (*A)(nil))
	assert.Nil(t, ret)
}

func TestCoalesce_ShouldReturnNil_WhenEmpty(t *testing.T) {
	type A struct {
		Name string
	}

	ret := Coalesce[*A]()
	assert.Nil(t, ret)
}

func TestCoalesce_ShouldReturnFirstElement_WhenOnlyOneElement(t *testing.T) {
	ret := Coalesce(100)
	assert.Equal(t, 100, ret)
}

func TestCoalesce_ShouldReturnFirstElement_WhenNoNilElems(t *testing.T) {
	ret := Coalesce(100, 300, 400, 500)
	assert.Equal(t, 100, ret)
}

func TestCoalesce_ShouldReturnFirstNonNilElement(t *testing.T) {
	type A struct {
		Name string
	}

	ret := Coalesce[*A](nil, &A{Name: "Vasya"}, &A{Name: "Kolya"})
	assert.Equal(t, "Vasya", ret.Name)
}

func TestCoalesce_ShouldReturnEmptyA_WhenZeroValue(t *testing.T) {
	type A struct {
		Name string
	}

	ret := Coalesce[*A](nil, &A{}, &A{Name: "Kolya"})
	assert.Equal(t, "", ret.Name)
}
