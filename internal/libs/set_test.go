package libs

import "testing"

func TestSet_Add(t *testing.T) {
	st := NewSet[int]()
	st.Add(1)

	if st.elements == nil {
		t.Error("Set[int] is nil")
	}

	if len(st.elements) != 1 {
		t.Error("Set[int] has incorrect number of elements")
	}

	if _, ok := st.elements[1]; !ok {
		t.Error("Set[int] has incorrect element")
	}
}

func TestSet_Contains(t *testing.T) {
	st := NewSet[int]()
	st.Add(1)

	got := st.Contains(1)
	if !got {
		t.Error("Set[int] should contain the element")
	}
}

func TestSet_Contains_NotContain(t *testing.T) {
	st := NewSet[int]()
	st.Add(1)

	got := st.Contains(2)
	if got {
		t.Error("Set[int] should not contain the element")
	}
}

func TestSet_Remove(t *testing.T) {
	st := NewSet[int]()
	st.Add(1)
	st.Remove(1)

	if len(st.elements) != 0 {
		t.Error("Set[int] has incorrect number of elements")
	}
}

func TestSet_Size(t *testing.T) {
	st := NewSet[int]()
	st.Add(1)
	st.Add(1)
	st.Add(1)

	got := st.Size()
	want := 1
	if got != want {
		t.Error("Set[int] has incorrect size")
	}
}
