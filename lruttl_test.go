/*
Copyright 2013 Google Inc.
Copyright 2017 Dual Inventive Technology Centre B.V.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package lruttl

import (
	"testing"
	"time"
)

type simpleStruct struct {
	int
	string
}

type complexStruct struct {
	int
	simpleStruct
}

var getTests = []struct {
	name       string
	keyToAdd   interface{}
	keyToGet   interface{}
	expectedOk bool
}{
	{"string_hit", "myKey", "myKey", true},
	{"string_miss", "myKey", "nonsense", false},
	{"simple_struct_hit", simpleStruct{1, "two"}, simpleStruct{1, "two"}, true},
	{"simeple_struct_miss", simpleStruct{1, "two"}, simpleStruct{0, "noway"}, false},
	{"complex_struct_hit", complexStruct{1, simpleStruct{2, "three"}},
		complexStruct{1, simpleStruct{2, "three"}}, true},
}

func TestGet(t *testing.T) {
	for _, tt := range getTests {
		lru := New(0, time.Hour)
		lru.Add(tt.keyToAdd, 1234)
		val, ok := lru.Get(tt.keyToGet)
		if ok != tt.expectedOk {
			t.Fatalf("%s: cache hit = %v; want %v", tt.name, ok, !ok)
		} else if ok && val != 1234 {
			t.Fatalf("%s expected get to return 1234 but got %v", tt.name, val)
		}
	}
}

func TestRemove(t *testing.T) {
	lru := New(0, time.Hour)
	lru.Add("myKey", 1234)
	if val, ok := lru.Get("myKey"); !ok {
		t.Fatal("TestRemove returned no match")
	} else if val != 1234 {
		t.Fatalf("TestRemove failed.  Expected %d, got %v", 1234, val)
	}

	lru.Remove("myKey")
	if _, ok := lru.Get("myKey"); ok {
		t.Fatal("TestRemove returned a removed entry")
	}
}

func TestTTL(t *testing.T) {
	lru := New(0, time.Millisecond*100)
	lru.Add("myKey", 1234)
	if val, ok := lru.Get("myKey"); !ok {
		t.Fatal("TestTTL returned no match")
	} else if val != 1234 {
		t.Fatalf("TestTTL failed.  Expected %d, got %v", 1234, val)
	}

	time.Sleep(1 * time.Second)
	if _, ok := lru.Get("myKey"); ok {
		t.Fatal("TestTTL returned a removed entry")
	}
}

func TestTTLReset(t *testing.T) {
	lru := New(0, time.Second)
	lru.Add("myKey", 1234)
	if val, ok := lru.Get("myKey"); !ok {
		t.Fatal("TestTTLReset returned no match")
	} else if val != 1234 {
		t.Fatalf("TestTTLReset failed.  Expected %d, got %v", 1234, val)
	}

	time.Sleep(500 * time.Millisecond)
	lru.Add("myKey", 5678)
	if val, ok := lru.Get("myKey"); !ok {
		t.Fatal("TestTTLReset returned no match")
	} else if val != 5678 {
		t.Fatalf("TestTTLReset failed.  Expected %d, got %v", 5678, val)
	}

	time.Sleep(500 * time.Millisecond)
	if val, ok := lru.Get("myKey"); !ok {
		t.Fatal("TestTTLReset returned no match")
	} else if val != 5678 {
		t.Fatalf("TestTTLReset failed.  Expected %d, got %v", 5678, val)
	}
	time.Sleep(700 * time.Millisecond)

	if _, ok := lru.Get("myKey"); ok {
		t.Fatal("TestTTLReset returned a removed entry")
	}
}

func TestAddNilCache(t *testing.T) {
	c := New(0, time.Hour)
	c.Clear()
	c.Add("a", 1)
	c.Clear()
	if _, ok := c.Get("a"); ok {
		t.Fatal("TestAddNilCache returned a removed entry")
	}
	c.Clear()
	c.Remove("a")
	c.Clear()
	c.RemoveOldest()
}

func TestAutoPrune(t *testing.T) {
	c := New(1, time.Hour)
	c.Add("a", 1)
	if _, ok := c.Get("a"); !ok {
		t.Fatal("TestAddNilCache returned no entry")
	}
	c.Add("b", 2)
	if _, ok := c.Get("a"); ok {
		t.Fatal("TestAddNilCache returned no entry")
	}
}
