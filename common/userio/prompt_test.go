package userio

import (
	"testing"

	"github.com/nutrixpos/pos/common/logger"
)

type nopLogger struct{}

func (nopLogger) Info(string, ...interface{})  {}
func (nopLogger) Warning(string, ...interface{}) {}
func (nopLogger) Error(string, ...interface{})   {}

func testLogger() logger.ILogger {
	return nopLogger{}
}

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

func TestPropagateCounterIndexToTree_Flat(t *testing.T) {
	m := &BubbleTeaSeedablesPrompter{Logger: testLogger()}
	tree := []PromptTreeElement{
		{Title: "A", Selected: false},
		{Title: "B", Selected: false},
		{Title: "C", Selected: false},
	}

	total, result := m.PropagateCounterIndexToTree(0, tree)

	if total != 3 {
		t.Errorf("expected total=3, got %d", total)
	}
	for i, elem := range result {
		if elem.CounterIndex != i {
			t.Errorf("element %d: expected CounterIndex=%d, got %d", i, i, elem.CounterIndex)
		}
	}
}

func TestPropagateCounterIndexToTree_WithOffset(t *testing.T) {
	m := &BubbleTeaSeedablesPrompter{Logger: testLogger()}
	tree := []PromptTreeElement{
		{Title: "X"},
		{Title: "Y"},
	}

	total, result := m.PropagateCounterIndexToTree(5, tree)

	if total != 7 {
		t.Errorf("expected total=7, got %d", total)
	}
	if result[0].CounterIndex != 5 {
		t.Errorf("first element: expected CounterIndex=5, got %d", result[0].CounterIndex)
	}
	if result[1].CounterIndex != 6 {
		t.Errorf("second element: expected CounterIndex=6, got %d", result[1].CounterIndex)
	}
}

func TestPropagateCounterIndexToTree_Nested(t *testing.T) {
	m := &BubbleTeaSeedablesPrompter{Logger: testLogger()}
	tree := []PromptTreeElement{
		{
			Title: "Parent",
			SubElements: []PromptTreeElement{
				{Title: "Child1"},
				{Title: "Child2"},
			},
		},
	}

	total, result := m.PropagateCounterIndexToTree(0, tree)

	if total != 3 {
		t.Errorf("expected total=3, got %d", total)
	}
	if result[0].CounterIndex != 0 {
		t.Errorf("parent: expected CounterIndex=0, got %d", result[0].CounterIndex)
	}
	if result[0].SubElements[0].CounterIndex != 1 {
		t.Errorf("child1: expected CounterIndex=1, got %d", result[0].SubElements[0].CounterIndex)
	}
	if result[0].SubElements[1].CounterIndex != 2 {
		t.Errorf("child2: expected CounterIndex=2, got %d", result[0].SubElements[1].CounterIndex)
	}
}

func TestPropagateCounterIndexToTree_DeepNested(t *testing.T) {
	m := &BubbleTeaSeedablesPrompter{Logger: testLogger()}
	tree := []PromptTreeElement{
		{
			Title: "Root",
			SubElements: []PromptTreeElement{
				{
					Title: "Branch",
					SubElements: []PromptTreeElement{
						{Title: "Leaf1"},
						{Title: "Leaf2"},
					},
				},
				{Title: "Leaf3"},
			},
		},
	}

	total, result := m.PropagateCounterIndexToTree(0, tree)

	if total != 6 {
		t.Errorf("expected total=6, got %d", total)
	}
	if result[0].CounterIndex != 0 {
		t.Errorf("root: expected 0, got %d", result[0].CounterIndex)
	}
	if result[0].SubElements[0].CounterIndex != 1 {
		t.Errorf("branch: expected 1, got %d", result[0].SubElements[0].CounterIndex)
	}
	if result[0].SubElements[0].SubElements[0].CounterIndex != 2 {
		t.Errorf("leaf1: expected 2, got %d", result[0].SubElements[0].SubElements[0].CounterIndex)
	}
	if result[0].SubElements[0].SubElements[1].CounterIndex != 3 {
		t.Errorf("leaf2: expected 3, got %d", result[0].SubElements[0].SubElements[1].CounterIndex)
	}
	if result[0].SubElements[1].CounterIndex != 5 {
		t.Errorf("leaf3: expected 5, got %d", result[0].SubElements[1].CounterIndex)
	}
}

func TestPropagateCounterIndexToTree_Empty(t *testing.T) {
	m := &BubbleTeaSeedablesPrompter{Logger: testLogger()}

	total, result := m.PropagateCounterIndexToTree(0, []PromptTreeElement{})

	if total != 0 {
		t.Errorf("expected total=0, got %d", total)
	}
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d elements", len(result))
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
