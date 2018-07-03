package maps

import (
	"encoding/json"
	"testing"
)

func TestPathSetArgs(t *testing.T) {
	examples := []struct {
		given      interface{}
		errMessage string
	}{
		{
			nil,
			"nil obj",
		},
		{
			map[interface{}]interface{}{},
			"non-pointer obj: map[interface {}]interface {}",
		},
		{
			&struct{}{},
			"non-map obj: *struct {}",
		},
	}

	for _, x := range examples {
		err := PathSet(x.given, "A", 10)
		if x.errMessage == "" && err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if x.errMessage != "" {
			if err == nil {
				t.Errorf("expected err %q, got no error", x.errMessage)
			} else if x.errMessage != err.Error() {
				t.Errorf("expected err %q, got %v", x.errMessage, err)
			}
		}
	}
}

func TestSet(t *testing.T) {
	examples := []struct {
		given      string
		path       string
		value      interface{}
		expected   string
		errMessage string
	}{
		{`{}`, "", 0, `{}`, "invalid path"},
		{`{}`, "[0]", 1, "{}", "invalid path"},
		{`{}`, "0", 2, "{}", "invalid path"},
		{`{}`, "0.a", 3, "{}", "invalid path"},
		{`{}`, "[0].a", 4, "{}", "invalid path"},

		{`{}`, "a", 10, `{"a":10}`, ""},
		{`{}`, "a.b", 100, `{"a":{"b":100}}`, ""},
		{`{}`, "a.b.c.d", true, `{"a":{"b":{"c":{"d":true}}}}`, ""},
		{`{}`, "a[0]", 5, `{"a":[5]}`, ""},
		{`{}`, "a.0", 5, `{"a":[5]}`, ""},
		{`{}`, "a[0].b", 10, `{"a":[{"b":10}]}`, ""},
		{`{}`, "a.0.b", 10, `{"a":[{"b":10}]}`, ""},
		{`{}`, "a[3].c", "test", `{"a":[null,null,null,{"c":"test"}]}`, ""},
		{`{}`, "a.3.c", "test", `{"a":[null,null,null,{"c":"test"}]}`, ""},
		{`{}`, "a[2].b[1]", 1, `{"a":[null,null,{"b":[null,1]}]}`, ""},
		{`{}`, "a[2].b[1].c[0]", 1, `{"a":[null,null,{"b":[null,{"c":[1]}]}]}`, ""},
		{`{}`, "a[0].b[0].c", 2, `{"a":[{"b":[{"c":2}]}]}`, ""},
		{`{}`, "a[0].b.c[0]", 3, `{"a":[{"b":{"c":[3]}}]}`, ""},
		{`{}`, "a[0].b.c[0].d", 4, `{"a":[{"b":{"c":[{"d":4}]}}]}`, ""},
		{`{}`, "a[0].b[1].c[2].d[3]", 5, `{"a":[{"b":[null,{"c":[null,null,{"d":[null,null,null,5]}]}]}]}`, ""},
		{`{}`, "a[0].b[1].c[2].d[3].e", 6, `{"a":[{"b":[null,{"c":[null,null,{"d":[null,null,null,{"e":6}]}]}]}]}`, ""},
		{`{}`, "a.b.c", struct{}{}, `{"a":{"b":{"c":{}}}}`, ""},
		{`{}`, "a.b.c", []struct{}{}, `{"a":{"b":{"c":[]}}}`, ""},
		{`{}`, "a.0.0.0.b", 10, `{"a":[[[{"b":10}]]]}`, ""},
		{`{}`, "a.0.b.1.c.2.d", 5, `{"a":[{"b":[null,{"c":[null,null,{"d":5}]}]}]}`, ""},

		{`{"a":1,"b":2}`, "a", 10, `{"a":10,"b":2}`, ""},
		{`{"a":1,"b":2}`, "c", 10, `{"a":1,"b":2,"c":10}`, ""},
		{`{"a":1,"b":2}`, "a.c", 10, `{"a":{"c":10},"b":2}`, ""},
		{`{"a":{"b":{"d":5}}}`, "a.b.c", 10, `{"a":{"b":{"c":10,"d":5}}}`, ""},
		{`{"a":[{"b":2}]}`, "a.1.c", 10, `{"a":[{"b":2},{"c":10}]}`, ""},
		{`{"a":[{"b":2}]}`, "a.2.c", 10, `{"a":[{"b":2},null,{"c":10}]}`, ""},

		// hacks
		{`{}`, " ", 0, `{" ":0}`, ""},
		{`{}`, "a. 0 ", 5, `{"a":{" 0 ":5}}`, ""},
		{`{}`, " 0. 1. 2. 3", 5, `{" 0":{" 1":{" 2":{" 3":5}}}}`, ""},
	}

	for i, x := range examples {
		var obj map[string]interface{}

		if err := json.Unmarshal([]byte(x.given), &obj); err != nil {
			t.Fatal(err)
		}

		err := PathSet(&obj, x.path, x.value)

		if x.errMessage == "" && err != nil {
			t.Errorf("example %d: expected no error, got %v", i, err)
		}
		if x.errMessage != "" {
			if err == nil {
				t.Errorf("example %d: expected err %q, got no error", i, x.errMessage)
			} else if x.errMessage != err.Error() {
				t.Errorf("example %d: expected err %q, got %v", i, x.errMessage, err)
			}
		}

		res, err := json.Marshal(obj)
		if err != nil {
			t.Fatal(err)
		}
		if string(res) != x.expected {
			t.Errorf("example %d: expected result %s, got %s", i, x.expected, string(res))
		}
	}
}
