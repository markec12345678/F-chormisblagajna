package userio

import (
	"testing"
)

func TestToggleSelectedTreeElement_Found(t *testing.T) {
	tree := []PromptTreeElement{
		{Title: "A", Selected: false, CounterIndex: 0},
		{Title: "B", Selected: false, CounterIndex: 1},
		{Title: "C", Selected: true, CounterIndex: 2},
	}

	result, found := ToggleSelectedTreeElement(1, tree)

	if !found {
		t.Fatal("expected to find element at index 1")
	}
	if !result[1].Selected {
		t.Error("element at index 1 should be selected after toggle")
	}
	if result[0].Selected {
		t.Error("element at index 0 should remain unselected")
	}
	if !result[2].Selected {
		t.Error("element at index 2 should remain selected")
	}
}

func TestToggleSelectedTreeElement_NotFound(t *testing.T) {
	tree := []PromptTreeElement{
		{Title: "A", Selected: false, CounterIndex: 0},
		{Title: "B", Selected: false, CounterIndex: 1},
	}

	result, found := ToggleSelectedTreeElement(99, tree)

	if found {
		t.Fatal("should not find element at index 99")
	}
	if len(result) != 2 {
		t.Errorf("tree should remain unchanged, got %d elements", len(result))
	}
}

func TestToggleSelectedTreeElement_Children(t *testing.T) {
	tree := []PromptTreeElement{
		{
			Title: "Parent", Selected: false, CounterIndex: 0,
			SubElements: []PromptTreeElement{
				{Title: "Child1", Selected: false, CounterIndex: 1},
				{Title: "Child2", Selected: false, CounterIndex: 2},
			},
		},
	}

	result, found := ToggleSelectedTreeElement(0, tree)

	if !found {
		t.Fatal("expected to find parent element")
	}
	if !result[0].Selected {
		t.Error("parent should be selected")
	}
	for i, child := range result[0].SubElements {
		if !child.Selected {
			t.Errorf("child %d should be selected when parent is selected", i)
		}
	}
}

func TestToggleSelectedTreeElement_DeselectParent(t *testing.T) {
	tree := []PromptTreeElement{
		{
			Title: "Parent", Selected: true, CounterIndex: 0,
			SubElements: []PromptTreeElement{
				{Title: "Child1", Selected: true, CounterIndex: 1},
			},
		},
	}

	result, found := ToggleSelectedTreeElement(0, tree)

	if !found {
		t.Fatal("expected to find parent element")
	}
	if result[0].Selected {
		t.Error("parent should be deselected")
	}
	if result[0].SubElements[0].Selected {
		t.Error("child should be deselected when parent is deselected")
	}
}

func TestToggleSelectedTreeElement_NestedChildren(t *testing.T) {
	tree := []PromptTreeElement{
		{
			Title: "Root", Selected: false, CounterIndex: 0,
			SubElements: []PromptTreeElement{
				{
					Title: "Child", Selected: false, CounterIndex: 1,
					SubElements: []PromptTreeElement{
						{Title: "Grandchild", Selected: false, CounterIndex: 2},
					},
				},
			},
		},
	}

	result, found := ToggleSelectedTreeElement(1, tree)

	if !found {
		t.Fatal("expected to find child element")
	}
	if result[0].Selected {
		t.Error("root should not be selected when child is toggled")
	}
	if !result[0].SubElements[0].Selected {
		t.Error("child should be selected")
	}
	if !result[0].SubElements[0].SubElements[0].Selected {
		t.Error("grandchild should be selected (direct children of toggled element propagate)")
	}
}

func TestPromptTreeElement_DefaultValues(t *testing.T) {
	e := PromptTreeElement{}

	if e.Title != "" {
		t.Errorf("default Title should be empty, got %q", e.Title)
	}
	if e.Selected {
		t.Error("default Selected should be false")
	}
	if e.Level != 0 {
		t.Errorf("default Level should be 0, got %d", e.Level)
	}
	if e.CounterIndex != 0 {
		t.Errorf("default CounterIndex should be 0, got %d", e.CounterIndex)
	}
	if e.SubElements != nil {
		t.Error("default SubElements should be nil")
	}
}
