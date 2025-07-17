//go:build ignore

package assert

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

func ExampleEqual(t *testing.T) {
	type User struct {
		Id        uuid.UUID
		Email     string
		CreatedAt time.Time
		Balance   decimal.Decimal

		active bool
	}

	loc, _ := time.LoadLocation("Europe/Warsaw")

	// db.CreateUser("test@example.com")
	createdAt := time.Now()
	user := User{
		Id:        uuid.New(),
		Email:     "test@example.com",
		CreatedAt: createdAt,
		Balance:   decimal.NewFromFloat(1),

		active: true,
	}

	Equal(
		t,
		user,
		User{
			Email:     "test@example.com",
			CreatedAt: createdAt.In(loc),
			Balance:   decimal.RequireFromString("1"),
		},
		IgnoreUnexported(),
		SkipEmptyFields(),
	)
}
