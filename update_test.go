package solr

import (
	"testing"
)

func TestNewUpdateBuilder(t *testing.T) {
	u := NewUpdateBuilder()
	if len(u.commands) != 0 {
		t.Fatal("update builder should be initialized")
	}
}

func TestUpdateBuilderAdd(t *testing.T) {
	u := NewUpdateBuilder()
	input := map[string]interface{}{"test": "test"}
	u.Add(input)
	_, ok := u.commands[CommandAdd]
	if !ok {
		t.Fatal("add command not found!")
	}
}

func TestUpdateBuilderCommit(t *testing.T) {
	u := NewUpdateBuilder()
	u.Commit(nil)
	_, ok := u.commands[CommandCommit]
	if !ok {
		t.Fatal("commit command not found!")
	}
}

func TestUpdateBuilderOptimize(t *testing.T) {
	u := NewUpdateBuilder()
	u.Optimize(nil)
	_, ok := u.commands[CommandOptimize]
	if !ok {
		t.Fatal("optimize command not found!")
	}
}

func TestUpdateBuilderRollback(t *testing.T) {
	u := NewUpdateBuilder()
	u.Rollback()
	_, ok := u.commands[CommandRollback]
	if !ok {
		t.Fatal("rollback command not found!")
	}
}

func TestUpdateBuilderDelete(t *testing.T) {
	u := NewUpdateBuilder()
	input := map[string]interface{}{"test": "test"}
	u.Delete(input)
	_, ok := u.commands[CommandDelete]
	if !ok {
		t.Fatal("delete command not found!")
	}
}

func TestNewUpdateDocument(t *testing.T) {
	upd := NewUpdateDocument("test")
	id, ok := upd.fields["id"]
	if !ok {
		t.Fatal("Id property not found!")
	}
	if upd.fields["id"] != "test" {
		t.Fatalf("expected id to be %s but got %s", "test", id)
	}
}

func TestUpdateAddDistinct(t *testing.T) {
	upd := NewUpdateDocument("test")
	input := []string{"one", "two", "one"}
	upd.AddDistinct("field", input)
	innerMap := upd.fields["field"].(map[string]interface{})
	actual, ok := innerMap[ActionAddDistinct]
	if !ok {
		t.Fatal("Add distinct property not found!")
	}
	if len(actual.([]string)) != len(input) {
		t.Fatalf("expected property to contain %d values but instead found %d", len(input), len(actual.([]string)))
	}
}

func TestUpdateAdd(t *testing.T) {
	upd := NewUpdateDocument("test")
	input := []string{"one", "two"}
	upd.Add("field", input)
	innerMap := upd.fields["field"].(map[string]interface{})
	actual, ok := innerMap[ActionAdd]
	if !ok {
		t.Fatal("Add property not found!")
	}
	if len(actual.([]string)) != len(input) {
		t.Fatalf("expected property to contain %d values but instead found %d", len(input), len(actual.([]string)))
	}
}

func TestUpdateSet(t *testing.T) {
	upd := NewUpdateDocument("test")
	input := "settest"
	upd.Set("field", input)
	innerMap := upd.fields["field"].(map[string]interface{})
	actual, ok := innerMap[ActionSet]
	if !ok {
		t.Fatal("Set property not found!")
	}
	if actual.(string) != input {
		t.Fatalf("expected property to be %s but instead got %s", input, actual.(string))
	}
}

func TestUpdateRemove(t *testing.T) {
	upd := NewUpdateDocument("test")
	input := "removetest"
	upd.Remove("field", input)
	innerMap := upd.fields["field"].(map[string]interface{})
	actual, ok := innerMap[ActionRemove]
	if !ok {
		t.Fatal("Remove property not found!")
	}
	if actual.(string) != input {
		t.Fatalf("expected property to be %s but instead got %s", input, actual.(string))
	}
}

func TestUpdateRemoveRegex(t *testing.T) {
	upd := NewUpdateDocument("test")
	input := "set*"
	upd.RemoveRegex("field", input)
	innerMap := upd.fields["field"].(map[string]interface{})
	actual, ok := innerMap[ActionRemoveRegex]
	if !ok {
		t.Fatal("RemoveRegex property not found!")
	}
	if actual.(string) != input {
		t.Fatalf("expected property to be %s but instead got %s", input, actual.(string))
	}
}

func TestUpdateIncrementBy(t *testing.T) {
	upd := NewUpdateDocument("test")
	input := 2
	upd.IncrementBy("field", input)
	innerMap := upd.fields["field"].(map[string]interface{})
	actual, ok := innerMap[ActionIncrement]
	if !ok {
		t.Fatal("Increment property not found!")
	}
	if actual.(int) != input {
		t.Fatalf("expected property to be %d but instead got %d", input, actual.(int))
	}
}
