package service

import (
	"errors"
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/fcm/mocks"
	"github.com/capstone-kelompok-7/backend-disappear/utils/sendnotif"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestFcmService_GetFcmById(t *testing.T) {
	repo := mocks.NewRepositoryFcmInterface(t)
	service := NewFcmService(repo)

	fcm := &entities.FcmModels{
		ID:    1,
		Title: "Test title 1",
		Body:  "Test body 1",
	}

	expectedFcm := &entities.FcmModels{
		ID:    fcm.ID,
		Title: fcm.Title,
		Body:  fcm.Body,
	}

	t.Run("Success case - Message found", func(t *testing.T) {
		fcmId := uint64(1)
		repo.On("GetFcmById", fcmId).Return(expectedFcm, nil).Once()

		result, err := service.GetFcmById(fcmId)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedFcm.ID, result.ID)
		assert.Equal(t, expectedFcm.Title, result.Title)
		assert.Equal(t, expectedFcm.Body, result.Body)
	})

	t.Run("Failed case - Message not found", func(t *testing.T) {
		fcmId := uint64(2)

		expectedErr := errors.New("Notifikasi tidak ditemukan")
		repo.On("GetFcmById", fcmId).Return(nil, expectedErr).Once()

		result, err := service.GetFcmById(fcmId)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedErr, err)

		repo.AssertExpectations(t)
	})
}

func TestFcmService_GetFcmByIdUser(t *testing.T) {
	repo := mocks.NewRepositoryFcmInterface(t)
	service := NewFcmService(repo)

	userID := uint64(1)

	expectedFcms := []*entities.FcmModels{
		{
			ID:    1,
			Title: "Test title 1",
			Body:  "Test body 1",
		},
	}

	t.Run("Success case - Notifications found", func(t *testing.T) {
		repo.On("GetFcmByIdUser", userID).Return(expectedFcms, nil).Once()

		result, err := service.GetFcmByIdUser(userID)

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedFcms, result)

		repo.AssertExpectations(t)
	})

	t.Run("Failed case - Notifications not found", func(t *testing.T) {
		expectedErr := errors.New("Notifikasi tidak ditemukan")
		repo.On("GetFcmByIdUser", userID).Return(nil, expectedErr).Once()

		result, err := service.GetFcmByIdUser(userID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedErr, err)

		repo.AssertExpectations(t)
	})
}

func TestFcmService_DeleteFcmById(t *testing.T) {
	repo := mocks.NewRepositoryFcmInterface(t)
	service := NewFcmService(repo)

	fcmID := uint64(1)

	t.Run("Success case - Fcm found and deleted", func(t *testing.T) {
		repo.On("GetFcmById", fcmID).Return(&entities.FcmModels{ID: fcmID}, nil).Once()
		repo.On("DeleteFcmById", fcmID).Return(nil).Once()

		err := service.DeleteFcmById(fcmID)

		assert.Nil(t, err)

		repo.AssertExpectations(t)
	})

	t.Run("Failed case - Fcm not found", func(t *testing.T) {
		repo.On("GetFcmById", fcmID).Return(nil, errors.New("fcm tidak ditemukan")).Once()

		err := service.DeleteFcmById(fcmID)

		assert.Error(t, err)
		assert.Equal(t, errors.New("fcm tidak ditemukan"), err)

		repo.AssertNotCalled(t, "DeleteFcmById")
	})

	t.Run("Failed case - Failed to delete Fcm", func(t *testing.T) {
		repo.On("GetFcmById", fcmID).Return(&entities.FcmModels{ID: fcmID}, nil).Once()
		repo.On("DeleteFcmById", fcmID).Return(errors.New("fcm id tidak ditemukan")).Once()

		err := service.DeleteFcmById(fcmID)

		assert.Error(t, err)
		assert.Equal(t, errors.New("fcm id tidak ditemukan"), err)

		repo.AssertExpectations(t)
	})
}

func TestFcmService_CreateFcm(t *testing.T) {
	repoMock := mocks.NewRepositoryFcmInterface(t)
	service := NewFcmService(repoMock)

	request := sendnotif.SendNotificationRequest{
		OrderID: "123",
		UserID:  456,
		Title:   "Test Title",
		Body:    "Test Body",
	}

	t.Run("Success case - Fcm created and notification sent", func(t *testing.T) {
		repoMock.On("SendMessageNotification", request).Return("success", nil).Once()
		repoMock.On("CreateFcm", mock.AnythingOfType("*entities.FcmModels")).Return(&entities.FcmModels{}, nil).Once()

		sendSuccess, response, err := service.CreateFcm(request)

		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, "success", sendSuccess)

		repoMock.AssertExpectations(t)
	})

	t.Run("Failed case - Failed to send notification", func(t *testing.T) {
		expectedErr := errors.New("failed to send notification")
		repoMock.On("SendMessageNotification", request).Return("", expectedErr).Once()

		sendSuccess, response, err := service.CreateFcm(request)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "", sendSuccess)

		repoMock.AssertExpectations(t)
	})

	t.Run("Failed case - Failed to create Fcm in the database", func(t *testing.T) {
		sendSuccess := "success"
		expectedErr := errors.New("failed to create Fcm in the database")
		repoMock.On("SendMessageNotification", request).Return(sendSuccess, nil).Once()
		repoMock.On("CreateFcm", mock.AnythingOfType("*entities.FcmModels")).Return(nil, expectedErr).Once()

		sendSuccess, response, err := service.CreateFcm(request)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, sendSuccess, "")

		repoMock.AssertExpectations(t)
	})
}
