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
	resultItem := anItem.AsFoldable().AsItem().(IntListItem)
	expected := []int{1}
	if !SliceEqual(resultItem.Value, expected) {
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
