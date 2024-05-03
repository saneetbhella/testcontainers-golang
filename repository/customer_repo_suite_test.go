package repository

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"log"
	"testcontainers-golang/model"
	"testcontainers-golang/testhelpers"
	"testing"
)

type CustomerRepoTestSuite struct {
	suite.Suite
	pgContainer *testhelpers.PostgresContainer
	repository  *Repository
	ctx         context.Context
}

func (suite *CustomerRepoTestSuite) SetupSuite() {
	ctx := context.Background()

	pgContainer, err := testhelpers.CreatePostgresContainer(ctx)
	if err != nil {
		log.Fatal("error setting up postgres testcontainer")
	}

	repository, err := NewRepository(ctx, pgContainer.ConnectionString)
	if err != nil {
		log.Fatal("error setting up repository", err)
	}

	suite.pgContainer = pgContainer
	suite.repository = repository
	suite.ctx = ctx
}

func (suite *CustomerRepoTestSuite) TearDownSuite() {
	if err := suite.repository.conn.Close(suite.ctx); err != nil {
		log.Fatal("error closing repository")
	}

	if err := suite.pgContainer.Postgres.Terminate(suite.ctx); err != nil {
		log.Fatal("error terminating postgres testcontainer")
	}
}

func (suite *CustomerRepoTestSuite) TestCreateCustomer() {
	t := suite.T()

	customer, err := suite.repository.CreateCustomer(suite.ctx, model.Customer{
		Name:  "Henry",
		Email: "henry@gmail.com",
	})
	assert.NoError(t, err)
	assert.NotNil(t, customer.Id)
}

func (suite *CustomerRepoTestSuite) TestGetCustomerByEmail() {
	t := suite.T()

	customer, err := suite.repository.GetCustomerByEmail(suite.ctx, "john@gmail.com")
	assert.NoError(t, err)
	assert.NotNil(t, customer)
	assert.Equal(t, "John", customer.Name)
	assert.Equal(t, "john@gmail.com", customer.Email)
}

func TestCustomerRepoTestSuite(t *testing.T) {
	suite.Run(t, new(CustomerRepoTestSuite))
}
