package repository

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testcontainers-golang/model"
	"testcontainers-golang/testhelpers"
	"testing"
)

func TestNewRepository(t *testing.T) {
	ctx := context.Background()
	container, err := testhelpers.CreatePostgresContainer(ctx)

	t.Cleanup(func() {
		if err := container.Postgres.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate pgContainer: %s", err)
		}
	})

	assert.NoError(t, err)

	customerRepo, err := NewRepository(ctx, container.ConnectionString)
	assert.NoError(t, err)

	c, err := customerRepo.CreateCustomer(ctx, model.Customer{
		Name:  "Henry",
		Email: "henry@gmail.com",
	})
	assert.NoError(t, err)
	assert.NotNil(t, c)

	customer, err := customerRepo.GetCustomerByEmail(ctx, "henry@gmail.com")
	assert.NoError(t, err)
	assert.NotNil(t, customer)
	assert.Equal(t, "Henry", customer.Name)
	assert.Equal(t, "henry@gmail.com", customer.Email)
}
