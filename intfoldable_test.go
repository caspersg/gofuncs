package funcs

import "testing"

func TestIntFoldableAppend(t *testing.T) {
	anItem := IntItem{Value: 5}

	aList := IntFoldable{List: []int{1, 2, 3, 4}}
	appendResult := aList.Append(anItem).(IntFoldable).List
	expected := []int{1, 2, 3, 4, 5}
	if !SliceEqual(appendResult, expected) {
		t.Errorf("result == %d expected %d", appendResult, expected)
	}
}

func TestIntFoldableConversions(t *testing.T) {
	aList := IntFoldable{List: []int{1, 2, 3, 4}}
	resultList := aList.AsItem().AsFoldable().(IntFoldable)
	if !SliceEqual(resultList.List, aList.List) {
		t.Errorf("result == %d expected %d", resultList, aList)
	}

	anItem := IntItem{Value: 1}
	resultItem := anItem.AsFoldable().AsItem().(IntFoldable)
	expected := []int{1}
	if !SliceEqual(resultItem.List, expected) {
		t.Errorf("result == %d expected %d", resultItem, expected)
	}
}

func TestIntFoldableMap(t *testing.T) {
	mapFunc := func(x Item) Item { return IntItem{Value: x.(IntItem).Value * 2} }
	in := []int{0, 1, 2}
	expected := []int{0, 2, 4}
	got := Map(IntFoldable{List: in}, mapFunc)
	if !SliceEqual(got.(IntFoldable).List, expected) {
		t.Errorf("result == %d expected %d", got, expected)
	}
}

func TestIntFoldableFilter(t *testing.T) {
	isNegative := func(x Item) bool { return x.(IntItem).Value < 0 }
	in := []int{0, -1, 1, 2, -30}
	expected := []int{-1, -30}
	got := Filter(IntFoldable{List: in}, isNegative)
	if !SliceEqual(got.(IntFoldable).List, expected) {
		t.Errorf("result == %d expected %d", got, expected)
	}
}

func TestIntFoldableLength(t *testing.T) {
	in := []int{0, -1, 1, 2, -30}
	length := Length(IntFoldable{List: in})
	if length != 5 {
		t.Errorf("result == %d expected %d", length, 5)
	}
}
func TestIntFoldableAll(t *testing.T) {
	isNegative := func(x Item) bool { return x.(IntItem).Value < 0 }
	in := []int{0, -1, 1, 2, -30}
	allNegative := All(IntFoldable{List: in}, isNegative)
	if allNegative {
		t.Errorf("result == %d expected %d", allNegative, false)
	}
}

func TestIntFoldableAny(t *testing.T) {
	isNegative := func(x Item) bool { return x.(IntItem).Value < 0 }
	in := []int{0, -1, 1, 2, -30}
	anyNegative := Any(IntFoldable{List: in}, isNegative)
	if !anyNegative {
		t.Errorf("result == %d expected %d", anyNegative, true)
	}
}

func TestIntFoldableConcat(t *testing.T) {
	a := []int{1, 2, 3}
	b := []int{4, 5, 6}
	expected := []int{1, 2, 3, 4, 5, 6}
	result := Concat(IntFoldable{List: a}, IntFoldable{List: b})
	if !SliceEqual(result.(IntFoldable).List, expected) {
		t.Errorf("result == %d expected %d", result, expected)
	}
}

func TestIntFoldableTake(t *testing.T) {
	x := []int{1, 2, 3, 4, 5, 6}
	expected := []int{1, 2, 3}
	result := Take(IntFoldable{List: x}, 3)
	if !SliceEqual(result.(IntFoldable).List, expected) {
		t.Errorf("result == %d expected %d", result, expected)
	}
}

func TestIntFoldableDrop(t *testing.T) {
	x := []int{1, 2, 3, 4, 5, 6}
	expected := []int{4, 5, 6}
	result := Drop(IntFoldable{List: x}, 3)
	if !SliceEqual(result.(IntFoldable).List, expected) {
		t.Errorf("result == %d expected %d", result, expected)
	}
}
