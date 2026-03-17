package stats

import "testing"

func TestGetValueReturnsBoolAsString(t *testing.T) {
	value, err := GetValue("Online", []byte(`{"Online":false}`))
	if err != nil {
		t.Fatalf("GetValue returned error: %v", err)
	}

	if value != "false" {
		t.Fatalf("expected false, got %q", value)
	}
}

func TestGetValueSupportsNestedKeys(t *testing.T) {
	value, err := GetValue("Used.CPU", []byte(`{"Used":{"CPU":"75%"}}`))
	if err != nil {
		t.Fatalf("GetValue returned error: %v", err)
	}

	if value != "75%" {
		t.Fatalf("expected 75%%, got %q", value)
	}
}

func TestSetValueUpdatesBool(t *testing.T) {
	updated, err := SetValue("Online", true, []byte(`{"Online":false}`))
	if err != nil {
		t.Fatalf("SetValue returned error: %v", err)
	}

	value, err := GetValue("Online", updated)
	if err != nil {
		t.Fatalf("GetValue returned error: %v", err)
	}

	if value != "true" {
		t.Fatalf("expected true, got %q", value)
	}
}