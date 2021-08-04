package user

import (
	"context"
	"github.com/google/uuid"
	"github.com/sposadas/mystorage/internal/domain"
	"github.com/sposadas/mystorage/util"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_sqlRepository_Store(t *testing.T) {
	db, err := util.InitDb()
	assert.NoError(t, err)

	myRepository := NewRepository(db)
	ctx := context.TODO()

	userId := uuid.New()
	user := domain.User{
		UUID: userId,
	}

	_, err = myRepository.Store(ctx, &user)
	assert.NoError(t, err)

	getResult, err := myRepository.GetOne(ctx, uuid.New())
	assert.NoError(t, err)
	assert.Nil(t, getResult)

	getResult, err = myRepository.GetOne(ctx, userId)
	assert.NoError(t, err)
	assert.NotNil(t, getResult)
	assert.Equal(t, user.UUID, getResult.UUID)
}

func Test_sqlRepository_Update(t *testing.T) {
	db, err := util.InitDb()
	assert.NoError(t, err)

	myRepository, ctx := NewRepository(db), context.TODO()

	userId := uuid.New()
	user := domain.User{
		UUID: userId,
		Username: "testing",
		Email: "test@mercadolibre.com",
	}

	_, err = myRepository.Store(ctx, &user)
	assert.NoError(t, err)

	getResult, err := myRepository.GetOne(ctx, userId)
	assert.NoError(t, err)
	assert.Equal(t, user.UUID, getResult.UUID)
	assert.Equal(t, user.Username, getResult.Username)
	assert.Equal(t, user.Email, getResult.Email)

	updatedUser := domain.User{
		ID: getResult.ID,
		UUID: getResult.UUID,
		Username: "test",
		Email: "newemail@mercadolibre.com",
		Active: true,
	}

	getResult, err = myRepository.Update(ctx, &updatedUser)
	assert.NoError(t, err)
	assert.Equal(t, user.UUID, getResult.UUID)
	assert.NotEqual(t, user.Username, getResult.Username)
	assert.NotEqual(t, user.Email, getResult.Email)
	assert.NotEqual(t, user.Active, getResult.Active)
}

func TestDelete(t *testing.T) {
	db, err := util.InitDb()
	assert.NoError(t, err)

	myRepository, ctx := NewRepository(db), context.TODO()

	userId := uuid.New()
	user := domain.User{
		UUID: userId,
		Username: "testing",
		Email: "test@mercadolibre.com",
	}

	_, err = myRepository.Store(ctx, &user)
	assert.NoError(t, err)

	getResult, err := myRepository.GetOne(ctx, userId)
	assert.NoError(t, err)
	assert.Equal(t, user.UUID, getResult.UUID)
	assert.Equal(t, user.Username, getResult.Username)
	assert.Equal(t, user.Email, getResult.Email)

	allUsers, err := myRepository.GetAll(ctx)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(allUsers))

	err = myRepository.Delete(ctx, userId)
	assert.NoError(t, err)

	getResult, err = myRepository.GetOne(ctx, userId)
	assert.NoError(t, err)
	assert.Nil(t, getResult)

	allUsers, err = myRepository.GetAll(ctx)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(allUsers))
}