package assert

import (
	"testing"
)

// Suite is a noop implementation of a test suite.
//
// It provides Setup and Teardown methods that can be overridden by
// specific test suites. This struct is designed to be embedded in
// custom test suite types, allowing you to inherit the basic suite
// behavior and override only the methods you need.
//
// Suite can be used if you don't want to define both Setup and Teardown
// methods in your test suite - you can embed Suite and override only
// the methods you actually need.
//
// Example:
//
//	type DatabaseTestSuite struct {
//		assert.Suite  // Embed to inherit noop implementations
//		db *sql.DB
//	}
//
//	// Only override Setup, Teardown inherited as noop
//	func (s *DatabaseTestSuite) Setup(t *testing.T) {
//		s.db = setupTestDB(t)
//	}
type Suite struct{}

// Setup is a no-op method provided for embedding.
// Override this method in your custom suite to perform test setup logic.
func (s Suite) Setup(t *testing.T) {}

// Teardown is a no-op method provided for embedding.
// Override this method in your custom suite to perform cleanup logic.
func (s Suite) Teardown(t *testing.T) {}

// Suiter defines the interface that all test suites must implement.
//
// This interface establishes the contract for test suite lifecycle management,
// providing hooks for initialization and cleanup operations. Any type that
// implements these two methods can be used as a test suite with the [Setup] function.
type Suiter interface {
	// Setup initializes the test suite before a test runs.
	// It receives a *testing.T which can be used for logging, assertions,
	// or controlling test execution (e.g., t.Skip(), t.Fatal()).
	Setup(t *testing.T)

	// Teardown cleans up the test suite after a test completes.
	// It receives a *testing.T which can be used for logging cleanup
	// operations or reporting cleanup failures.
	Teardown(t *testing.T)
}

// Setup allocates, initializes, and returns a test suite of type S.
//
// It performs the following steps:
//
//  1. Allocates a new instance of the suite (S)
//  2. Calls the suite's Setup(t *testing.T) method.
//  3. Registers the suite's Teardown method using t.Cleanup.
//
// This ensures that test lifecycle hooks are consistently applied
// and automatically cleaned up, even on test failures.
func Setup[V any, S interface {
	*V
	Suiter
}](t *testing.T) S {
	t.Helper()

	s := S(new(V))
	s.Setup(t)
	t.Cleanup(func() {
		s.Teardown(t)
	})
	return s
}
