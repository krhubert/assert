package assert

import (
	"testing"
)

type testSuite struct {
	r int
}

func (s *testSuite) Setup(t *testing.T) {
	s.r = 4
}

func (s *testSuite) Teardown(t *testing.T) {
	s.r = 0
}

func TestSuite(t *testing.T) {
	var suite *testSuite

	t.Run("TestSuite", func(t *testing.T) {
		suite = Setup[testSuite](t)
		if suite.r != 4 {
			t.Errorf("Expected r to be 4 after setup, got %d", suite.r)
		}
	})

	if suite == nil {
		t.Fatal("Expected suite to be initialized, got nil")
	}
	if suite.r != 0 {
		t.Fatalf("Expected r to be 0 after teardown, got %d", suite.r)
	}
}
