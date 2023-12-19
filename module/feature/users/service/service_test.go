package service

import (
	"errors"
	"fmt"
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/users/dto"
	userMocks "github.com/capstone-kelompok-7/backend-disappear/module/feature/users/mocks"
	utils "github.com/capstone-kelompok-7/backend-disappear/utils/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
)

func setupTestService(t *testing.T) (
	*UserService,
	*userMocks.RepositoryUserInterface,
	*utils.HashInterface) {

	repo := userMocks.NewRepositoryUserInterface(t)
	hash := utils.NewHashInterface(t)
	service := NewUserService(repo, hash)

	return service.(*UserService), repo, hash
}

func TestUserService_PaginationFunctions(t *testing.T) {
	service := &UserService{}

	t.Run("CalculatePaginationValues", func(t *testing.T) {
		// Test case 1
		pageInt, totalPages := service.CalculatePaginationValues(1, 20, 5)
		assert.Equal(t, 1, pageInt)
		assert.Equal(t, 4, totalPages)

		// Test case 2
		pageInt, totalPages = service.CalculatePaginationValues(-1, 15, 5)
		assert.Equal(t, 1, pageInt)
		assert.Equal(t, 3, totalPages)

		// Test case 3
		pageInt, totalPages = service.CalculatePaginationValues(7, 50, 10)
		assert.Equal(t, 5, pageInt)
		assert.Equal(t, 5, totalPages)
	})

	t.Run("GetNextPage", func(t *testing.T) {
		// Test case 1
		nextPage := service.GetNextPage(3, 5)
		assert.Equal(t, 4, nextPage)

		// Test case 2
		nextPage = service.GetNextPage(5, 5)
		assert.Equal(t, 5, nextPage)

		// Test case 3
		nextPage = service.GetNextPage(8, 10)
		assert.Equal(t, 9, nextPage)
	})

	t.Run("GetPrevPage", func(t *testing.T) {
		// Test case 1
		prevPage := service.GetPrevPage(3)
		assert.Equal(t, 2, prevPage)

		// Test case 2
		prevPage = service.GetPrevPage(1)
		assert.Equal(t, 1, prevPage)

		// Test case 3
		prevPage = service.GetPrevPage(7)
		assert.Equal(t, 6, prevPage)
	})
}

func TestUserService_GetUsersById(t *testing.T) {
	userID := uint64(1)
	user := &entities.UserModels{
		ID:   userID,
		Name: "User 1",
		Role: "customer",
	}
	service, repo, _ := setupTestService(t)
	t.Run("Failed Case - User Not Found", func(t *testing.T) {
		expectedErr := errors.New("pengguna tidak ditemukan")
		repo.On("GetUsersById", userID).Return(nil, expectedErr).Once()

		result, err := service.GetUsersById(userID)

		assert.Nil(t, result)
		assert.Equal(t, expectedErr, err)
		repo.AssertExpectations(t)
	})

	t.Run("Success Case", func(t *testing.T) {
		repo.On("GetUsersById", userID).Return(user, nil)

		result, err := service.GetUsersById(userID)

		assert.NoError(t, err)
		assert.NotNil(t, result)

		repo.AssertExpectations(t)
	})

}

func TestUserService_GetUsersByEmail(t *testing.T) {
	userID := uint64(1)
	userEmail := "email@gmail.com"
	user := &entities.UserModels{
		ID:    userID,
		Email: userEmail,
		Name:  "User 1",
		Role:  "customer",
	}
	service, repo, _ := setupTestService(t)
	t.Run("Failed Case - User Not Found", func(t *testing.T) {
		expectedErr := errors.New("pengguna tidak ditemukan")
		repo.On("GetUsersByEmail", userEmail).Return(nil, expectedErr).Once()

		result, err := service.GetUsersByEmail(userEmail)

		assert.Nil(t, result)
		assert.Equal(t, expectedErr, err)
		repo.AssertExpectations(t)
	})

	t.Run("Success Case", func(t *testing.T) {
		repo.On("GetUsersByEmail", userEmail).Return(user, nil)

		result, err := service.GetUsersByEmail(userEmail)

		assert.NoError(t, err)
		assert.NotNil(t, result)

		repo.AssertExpectations(t)
	})

}

func TestUserService_ChangePassword(t *testing.T) {
	userID := uint64(1)
	user := &entities.UserModels{
		ID:   userID,
		Name: "User 1",
		Role: "customer",
	}
	request := dto.UpdatePasswordRequest{
		OldPassword:     "pass123",
		NewPassword:     "pass1234",
		ConfirmPassword: "pass1234",
	}

	service, repo, hash := setupTestService(t)
	t.Run("Failed Case - User Not Found", func(t *testing.T) {
		expectedErr := errors.New("pengguna tidak ditemukan")
		repo.On("GetUsersById", userID).Return(nil, expectedErr).Once()

		err := service.ChangePassword(userID, request)

		assert.NotNil(t, err)
		assert.EqualError(t, err, "pengguna tidak ditemukan")
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Error Hash Password", func(t *testing.T) {
		expectedErr := errors.New("gagal hash password")
		repo.On("GetUsersById", userID).Return(user, nil).Once()
		hash.On("GenerateHash", request.NewPassword).Return("", expectedErr).Once()

		err := service.ChangePassword(userID, request)

		assert.NotNil(t, err)
		assert.EqualError(t, err, "gagal hash password")
		repo.AssertExpectations(t)
		hash.AssertExpectations(t)
	})

	t.Run("Failed Case - Error Change Password", func(t *testing.T) {
		hashedPassword := "hashedpassword"
		expectedErr := errors.New("some error")
		repo.On("GetUsersById", userID).Return(user, nil).Once()

		hash.On("GenerateHash", request.NewPassword).Return(hashedPassword, nil).Once()
		repo.On("ChangePassword", userID, hashedPassword).Return(expectedErr).Once()

		err := service.ChangePassword(userID, request)

		assert.NotNil(t, err)
		assert.Equal(t, err, expectedErr)
		repo.AssertExpectations(t)
		hash.AssertExpectations(t)
	})

	t.Run("Success Case - Password Changed", func(t *testing.T) {
		hashedPassword := "hashedpassword"
		repo.On("GetUsersById", userID).Return(user, nil).Once()

		hash.On("GenerateHash", request.NewPassword).Return(hashedPassword, nil).Once()
		repo.On("ChangePassword", userID, hashedPassword).Return(nil).Once()

		err := service.ChangePassword(userID, request)

		assert.Nil(t, err)
		repo.AssertExpectations(t)
		hash.AssertExpectations(t)
	})
}

func TestUserService_EditProfile(t *testing.T) {
	userID := uint64(1)
	user := &entities.UserModels{
		ID:           userID,
		Name:         "User 1",
		Phone:        "09888781237772",
		PhotoProfile: "photo.jpg",
	}
	request := dto.EditProfileRequest{
		Name:         "User 2",
		Phone:        "123123123312",
		PhotoProfile: "photos.jpg",
	}

	updatedUser := &entities.UserModels{
		Name:         "User 2",
		Phone:        "123123123312",
		PhotoProfile: "photos.jpg",
	}

	service, repo, _ := setupTestService(t)
	t.Run("Failed Case - User Not Found", func(t *testing.T) {
		expectedErr := errors.New("pengguna tidak ditemukan")
		repo.On("GetUsersById", userID).Return(nil, expectedErr).Once()

		result, err := service.EditProfile(userID, request)

		assert.Nil(t, result)
		assert.Equal(t, expectedErr, err)
		repo.AssertExpectations(t)
	})

	t.Run("Success Case - Profile Edited", func(t *testing.T) {
		repo.On("GetUsersById", userID).Return(user, nil).Once()
		repo.On("EditProfile", userID, request).Return(user, nil).Once()

		result, err := service.EditProfile(userID, request)

		assert.NotNil(t, result)
		assert.Nil(t, err)
		assert.Equal(t, request.Name, updatedUser.Name)
		assert.Equal(t, request.Phone, updatedUser.Phone)
		assert.Equal(t, request.PhotoProfile, updatedUser.PhotoProfile)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Error Editing Profile", func(t *testing.T) {
		expectedErr := errors.New("gagal mengedit profil")
		repo.On("GetUsersById", userID).Return(user, nil).Once()
		repo.On("EditProfile", userID, request).Return(nil, expectedErr).Once()

		result, err := service.EditProfile(userID, request)

		assert.Nil(t, result)
		assert.Equal(t, expectedErr, err)
		repo.AssertExpectations(t)
	})
}

func TestUserService_DeleteAccount(t *testing.T) {
	userID := uint64(1)
	user := &entities.UserModels{
		ID:   userID,
		Name: "User 1",
		Role: "customer",
	}

	service, repo, _ := setupTestService(t)

	t.Run("Failed Case - User Not Found", func(t *testing.T) {
		expectedErr := errors.New("pengguna tidak ditemukan")
		repo.On("GetUsersById", userID).Return(nil, expectedErr).Once()

		err := service.DeleteAccount(userID)

		assert.Equal(t, expectedErr, err)
		repo.AssertExpectations(t)
	})

	t.Run("Success Case - Account Deleted", func(t *testing.T) {
		repo.On("GetUsersById", userID).Return(user, nil).Once()
		repo.On("DeleteAccount", userID).Return(nil).Once()

		err := service.DeleteAccount(userID)

		assert.Nil(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Error Deleting Account", func(t *testing.T) {
		expectedErr := errors.New("gagal menghapus akun")
		repo.On("GetUsersById", userID).Return(user, nil).Once()
		repo.On("DeleteAccount", userID).Return(expectedErr).Once()

		err := service.DeleteAccount(userID)

		assert.Equal(t, expectedErr, err)
		repo.AssertExpectations(t)
	})

}

func TestDetermineLevel(t *testing.T) {
	testCases := []struct {
		exp      uint64
		expected string
	}{
		{exp: 400, expected: "Bronze"},
		{exp: 700, expected: "Silver"},
		{exp: 1500, expected: "Gold"},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("Exp %d", testCase.exp), func(t *testing.T) {
			result := determineLevel(testCase.exp)
			assert.Equal(t, testCase.expected, result)
		})
	}
}

func TestUserService_UpdateUserExp(t *testing.T) {
	userID := uint64(1)
	exp := uint64(800)
	user := &entities.UserModels{
		ID:   userID,
		Name: "User 1",
		Exp:  500,
	}

	t.Run("Failed Case - Error Updating Level", func(t *testing.T) {
		service, repo, _ := setupTestService(t)
		expectedErr := errors.New("gagal memperbarui Level")
		repo.On("GetUsersById", userID).Return(user, nil).Once()
		repo.On("UpdateUserExp", userID, exp).Return(user, nil).Once()
		repo.On("UpdateUserLevel", userID, "Silver").Return(expectedErr).Once()

		result, err := service.UpdateUserExp(userID, exp)

		assert.Nil(t, result)
		assert.Equal(t, expectedErr, err)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - User Not Found", func(t *testing.T) {
		service, repo, _ := setupTestService(t)
		expectedErr := errors.New("pengguna tidak ditemukan")
		repo.On("GetUsersById", userID).Return(nil, expectedErr).Once()

		result, err := service.UpdateUserExp(userID, exp)

		assert.Nil(t, result)
		assert.Equal(t, expectedErr, err)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Error Updating Exp", func(t *testing.T) {
		service, repo, _ := setupTestService(t)
		expectedErr := errors.New("gagal memperbarui Exp")
		repo.On("GetUsersById", userID).Return(user, nil).Once()
		repo.On("UpdateUserExp", userID, exp).Return(nil, expectedErr).Once()

		result, err := service.UpdateUserExp(userID, exp)

		assert.Nil(t, result)
		assert.Equal(t, expectedErr, err)
		repo.AssertExpectations(t)
	})

	t.Run("Success Case - Exp Updated", func(t *testing.T) {
		service, repo, _ := setupTestService(t)
		user.Exp = exp
		repo.On("GetUsersById", userID).Return(user, nil).Once()
		repo.On("UpdateUserExp", userID, exp).Return(user, nil).Once()

		result, err := service.UpdateUserExp(userID, exp)

		assert.NotNil(t, result)
		assert.NoError(t, err)
		assert.Equal(t, "Silver", result.Level)
		repo.AssertExpectations(t)
	})

}

func TestUserService_UpdateUserContribution(t *testing.T) {
	userID := uint64(1)
	gramPlastic := uint64(500)

	user := &entities.UserModels{
		ID:        userID,
		Name:      "User 1",
		TotalGram: 0,
	}

	service, repo, _ := setupTestService(t)

	t.Run("Failed Case - User Not Found", func(t *testing.T) {
		expectedErr := errors.New("pengguna tidak ditemukan")
		repo.On("GetUsersById", userID).Return(nil, expectedErr).Once()

		result, err := service.UpdateUserContribution(userID, gramPlastic)

		assert.Nil(t, result)
		assert.Equal(t, expectedErr, err)
		repo.AssertExpectations(t)
	})

	t.Run("Success Case - Contribution Updated", func(t *testing.T) {
		repo.On("GetUsersById", userID).Return(user, nil).Once()
		repo.On("UpdateUserContribution", userID, gramPlastic).Return(user, nil).Once()

		result, err := service.UpdateUserContribution(userID, gramPlastic)

		assert.NotNil(t, result)
		assert.NoError(t, err)
		assert.Equal(t, gramPlastic, result.TotalGram)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Error Updating Contribution", func(t *testing.T) {
		expectedErr := errors.New("gagal memperbarui kontribusi")
		repo.On("GetUsersById", userID).Return(user, nil).Once()
		repo.On("UpdateUserContribution", userID, gramPlastic).Return(nil, expectedErr).Once()

		result, err := service.UpdateUserContribution(userID, gramPlastic)

		assert.Nil(t, result)
		assert.Equal(t, expectedErr, err)
		repo.AssertExpectations(t)
	})
}

func TestUserService_GetUserLevel(t *testing.T) {
	userID := uint64(1)
	expectedLevel := "Silver"

	user := &entities.UserModels{
		ID:    userID,
		Name:  "User 1",
		Level: expectedLevel,
	}

	service, repo, _ := setupTestService(t)

	t.Run("Success Case - Get User Level", func(t *testing.T) {
		repo.On("GetUsersById", userID).Return(user, nil).Once()
		repo.On("GetUserLevel", userID).Return(expectedLevel, nil).Once()

		resultLevel, err := service.GetUserLevel(userID)

		assert.Nil(t, err)
		assert.Equal(t, expectedLevel, resultLevel)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - User Not Found", func(t *testing.T) {
		expectedErr := errors.New("pengguna tidak ditemukan")
		repo.On("GetUsersById", userID).Return(nil, expectedErr).Once()

		resultLevel, err := service.GetUserLevel(userID)

		assert.Empty(t, resultLevel)
		assert.Equal(t, expectedErr, err)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Error Getting User Level", func(t *testing.T) {
		expectedErr := errors.New("gagal mendapatkan level pengguna")
		repo.On("GetUsersById", userID).Return(user, nil).Once()
		repo.On("GetUserLevel", userID).Return("", expectedErr).Once()

		resultLevel, err := service.GetUserLevel(userID)

		assert.Empty(t, resultLevel)
		assert.Equal(t, expectedErr, err)
		repo.AssertExpectations(t)
	})
}

func TestUserService_UpdateUserChallengeFollow(t *testing.T) {
	userID := uint64(1)
	expectedChallenge := uint64(5)

	user := &entities.UserModels{
		ID:             userID,
		TotalChallenge: expectedChallenge,
	}

	service, repo, _ := setupTestService(t)

	t.Run("Success Case - Update User Challenge Follow", func(t *testing.T) {
		repo.On("UpdateUserChallengeFollow", userID, expectedChallenge).Return(user, nil).Once()

		resultUser, err := service.UpdateUserChallengeFollow(userID, expectedChallenge)

		assert.Nil(t, err)
		assert.Equal(t, expectedChallenge, resultUser.TotalChallenge)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Error Updating User Challenge Follow", func(t *testing.T) {
		expectedErr := errors.New("gagal memperbarui challenge pengguna")
		repo.On("UpdateUserChallengeFollow", userID, expectedChallenge).Return(nil, expectedErr).Once()

		resultUser, err := service.UpdateUserChallengeFollow(userID, expectedChallenge)

		assert.Nil(t, resultUser)
		assert.Equal(t, expectedErr, err)
		repo.AssertExpectations(t)
	})
}

func TestUserService_GetLeaderboardByExp(t *testing.T) {
	limit := 10
	expectedUsers := []*entities.UserModels{
		&entities.UserModels{
			ID:   1,
			Name: "User 1",
		},
		&entities.UserModels{
			ID:   2,
			Name: "User 2",
		},
	}

	service, repo, _ := setupTestService(t)

	t.Run("Success Case - Get Leaderboard by Exp", func(t *testing.T) {
		repo.On("GetLeaderboardByExp", limit).Return(expectedUsers, nil).Once()

		users, err := service.GetLeaderboardByExp(limit)

		assert.Nil(t, err)
		assert.Equal(t, len(expectedUsers), len(users))
		for i := range expectedUsers {
			assert.Equal(t, expectedUsers[i].ID, users[i].ID)
		}
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Error Getting Leaderboard by Exp", func(t *testing.T) {
		expectedErr := errors.New("gagal mengambil leaderboard")
		repo.On("GetLeaderboardByExp", limit).Return(nil, expectedErr).Once()

		users, err := service.GetLeaderboardByExp(limit)

		assert.Nil(t, users)
		assert.Equal(t, expectedErr, err)
		repo.AssertExpectations(t)
	})
}

func TestUserService_GetUserTransactionActivity(t *testing.T) {
	userID := uint64(1)
	expectedSuccessOrder := 5
	expectedFailedOrder := 2
	expectedTotalOrder := 7

	service, repo, _ := setupTestService(t)

	t.Run("Success Case - Get User Transaction Activity", func(t *testing.T) {
		repo.On("GetUsersById", userID).Return(&entities.UserModels{ID: userID}, nil).Once()
		repo.On("GetUserTransactionActivity", userID).Return(expectedSuccessOrder, expectedFailedOrder, expectedTotalOrder, nil).Once()

		successOrder, failedOrder, totalOrder, err := service.GetUserTransactionActivity(userID)

		assert.Nil(t, err)
		assert.Equal(t, expectedSuccessOrder, successOrder)
		assert.Equal(t, expectedFailedOrder, failedOrder)
		assert.Equal(t, expectedTotalOrder, totalOrder)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - User Not Found", func(t *testing.T) {
		expectedErr := errors.New("pengguna tidak ditemukan")
		repo.On("GetUsersById", userID).Return(nil, expectedErr).Once()

		successOrder, failedOrder, totalOrder, err := service.GetUserTransactionActivity(userID)

		assert.EqualError(t, err, expectedErr.Error())
		assert.Zero(t, successOrder)
		assert.Zero(t, failedOrder)
		assert.Zero(t, totalOrder)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Error Getting User Transaction Activity", func(t *testing.T) {
		expectedErr := errors.New("gagal mengambil aktivitas transaksi pengguna")
		repo.On("GetUsersById", userID).Return(&entities.UserModels{ID: userID}, nil).Once()
		repo.On("GetUserTransactionActivity", userID).Return(0, 0, 0, expectedErr).Once()

		successOrder, failedOrder, totalOrder, err := service.GetUserTransactionActivity(userID)

		assert.EqualError(t, err, expectedErr.Error())
		assert.Zero(t, successOrder)
		assert.Zero(t, failedOrder)
		assert.Zero(t, totalOrder)
		repo.AssertExpectations(t)
	})
}

func TestUserService_GetUserChallengeActivity(t *testing.T) {
	userID := uint64(1)
	expectedSuccessChallenge := 8
	expectedFailedChallenge := 3
	expectedTotalChallenge := 11

	service, repo, _ := setupTestService(t)

	t.Run("Success Case - Get User Challenge Activity", func(t *testing.T) {
		repo.On("GetUsersById", userID).Return(&entities.UserModels{ID: userID}, nil).Once()
		repo.On("GetUserChallengeActivity", userID).Return(expectedSuccessChallenge, expectedFailedChallenge, expectedTotalChallenge, nil).Once()

		successChallenge, failedChallenge, totalChallenge, err := service.GetUserChallengeActivity(userID)

		assert.Nil(t, err)
		assert.Equal(t, expectedSuccessChallenge, successChallenge)
		assert.Equal(t, expectedFailedChallenge, failedChallenge)
		assert.Equal(t, expectedTotalChallenge, totalChallenge)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - User Not Found", func(t *testing.T) {
		expectedErr := errors.New("pengguna tidak ditemukan")
		repo.On("GetUsersById", userID).Return(nil, expectedErr).Once()

		successChallenge, failedChallenge, totalChallenge, err := service.GetUserChallengeActivity(userID)

		assert.EqualError(t, err, expectedErr.Error())
		assert.Zero(t, successChallenge)
		assert.Zero(t, failedChallenge)
		assert.Zero(t, totalChallenge)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Error Getting User Challenge Activity", func(t *testing.T) {
		expectedErr := errors.New("gagal mengambil aktivitas tantangan pengguna")
		repo.On("GetUsersById", userID).Return(&entities.UserModels{ID: userID}, nil).Once()
		repo.On("GetUserChallengeActivity", userID).Return(0, 0, 0, expectedErr).Once()

		successChallenge, failedChallenge, totalChallenge, err := service.GetUserChallengeActivity(userID)

		assert.EqualError(t, err, expectedErr.Error())
		assert.Zero(t, successChallenge)
		assert.Zero(t, failedChallenge)
		assert.Zero(t, totalChallenge)
		repo.AssertExpectations(t)
	})
}

func TestUserService_GetUserProfile(t *testing.T) {
	userID := uint64(1)
	expectedUser := &entities.UserModels{
		ID:    userID,
		Name:  "John Doe",
		Email: "john@example.com",
		Phone: "123456789",
		Role:  "user",
	}

	service, repo, _ := setupTestService(t)

	t.Run("Success Case - Get User Profile", func(t *testing.T) {
		repo.On("GetUsersById", userID).Return(expectedUser, nil).Once()

		userProfile, err := service.GetUserProfile(userID)

		assert.Nil(t, err)
		assert.Equal(t, expectedUser, userProfile)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Error Getting User Profile", func(t *testing.T) {
		expectedErr := errors.New("failed to get user profile")
		repo.On("GetUsersById", userID).Return(nil, expectedErr).Once()

		userProfile, err := service.GetUserProfile(userID)

		assert.EqualError(t, err, expectedErr.Error())
		assert.Nil(t, userProfile)
		repo.AssertExpectations(t)
	})
}

func TestUserService_GetUsersBySearchAndFilter(t *testing.T) {
	page := 1
	perPage := 10
	search := "John"
	levelFilter := "Gold"

	expectedUsers := []*entities.UserModels{
		&entities.UserModels{ID: 1, Name: "John Doe", Level: "Gold"},
		&entities.UserModels{ID: 2, Name: "John Smith", Level: "Gold"},
	}

	expectedTotalUsers := int64(len(expectedUsers))

	service, repo, _ := setupTestService(t)

	t.Run("Success Case - Get Users By Search And Filter", func(t *testing.T) {
		repo.On("GetAllUsersBySearchAndFilter", page, perPage, search, levelFilter).Return(expectedUsers, expectedTotalUsers, nil).Once()

		users, totalUsers, err := service.GetUsersBySearchAndFilter(page, perPage, search, levelFilter)

		assert.Nil(t, err)
		assert.Equal(t, expectedTotalUsers, totalUsers)
		assert.Equal(t, expectedUsers, users)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Error Getting Users By Search And Filter", func(t *testing.T) {
		expectedErr := errors.New("failed to get users by search and filter")
		repo.On("GetAllUsersBySearchAndFilter", page, perPage, search, levelFilter).Return(nil, int64(0), expectedErr).Once()

		users, totalUsers, err := service.GetUsersBySearchAndFilter(page, perPage, search, levelFilter)

		assert.EqualError(t, err, expectedErr.Error())
		assert.Nil(t, users)
		assert.Zero(t, totalUsers)
		repo.AssertExpectations(t)
	})
}

func TestUserService_GetUsersByLevel(t *testing.T) {
	page := 1
	perPage := 10
	level := "Gold"

	expectedUsers := []*entities.UserModels{
		&entities.UserModels{ID: 1, Name: "John Doe", Level: "Gold"},
		&entities.UserModels{ID: 2, Name: "Jane Smith", Level: "Gold"},
	}

	expectedTotalUsers := int64(len(expectedUsers))

	service, repo, _ := setupTestService(t)

	t.Run("Success Case - Get Users By Level", func(t *testing.T) {
		repo.On("GetFilterLevel", page, perPage, level).Return(expectedUsers, expectedTotalUsers, nil).Once()

		users, totalUsers, err := service.GetUsersByLevel(page, perPage, level)

		assert.Nil(t, err)
		assert.Equal(t, expectedTotalUsers, totalUsers)
		assert.Equal(t, expectedUsers, users)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Error Getting Users By Level", func(t *testing.T) {
		expectedErr := errors.New("failed to get users by level")
		repo.On("GetFilterLevel", page, perPage, level).Return(nil, int64(0), expectedErr).Once()

		users, totalUsers, err := service.GetUsersByLevel(page, perPage, level)

		assert.EqualError(t, err, expectedErr.Error())
		assert.Nil(t, users)
		assert.Zero(t, totalUsers)
		repo.AssertExpectations(t)
	})
}

func TestUserService_GetUsersByName(t *testing.T) {
	page := 1
	perPage := 10
	name := "John Doe"

	expectedUsers := []*entities.UserModels{
		&entities.UserModels{ID: 1, Name: "John Doe"},
		&entities.UserModels{ID: 2, Name: "John Doe"},
	}

	expectedTotalItems := int64(len(expectedUsers))

	service, repo, _ := setupTestService(t)

	t.Run("Success Case - Get Users By Name", func(t *testing.T) {
		repo.On("FindByName", page, perPage, name).Return(expectedUsers, nil).Once()
		repo.On("GetTotalUserCountByName", name).Return(expectedTotalItems, nil).Once()

		users, totalItems, err := service.GetUsersByName(page, perPage, name)

		assert.Nil(t, err)
		assert.Equal(t, expectedTotalItems, totalItems)
		assert.Equal(t, expectedUsers, users)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Error Getting Users By Name", func(t *testing.T) {
		expectedErr := errors.New("failed to get users by name")
		repo.On("FindByName", page, perPage, name).Return(nil, expectedErr).Once()

		users, totalItems, err := service.GetUsersByName(page, perPage, name)

		assert.EqualError(t, err, expectedErr.Error())
		assert.Nil(t, users)
		assert.Zero(t, totalItems)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Error Getting Total User Count By Name", func(t *testing.T) {
		expectedErr := errors.New("failed to get total user count by name")
		repo.On("FindByName", page, perPage, name).Return(expectedUsers, nil).Once()
		repo.On("GetTotalUserCountByName", name).Return(int64(0), expectedErr).Once()

		users, totalItems, err := service.GetUsersByName(page, perPage, name)

		assert.EqualError(t, err, expectedErr.Error())
		assert.Nil(t, users)
		assert.Zero(t, totalItems)
		repo.AssertExpectations(t)
	})
}

func TestUserService_GetAllUsers(t *testing.T) {
	page := 1
	perPage := 10

	expectedUsers := []*entities.UserModels{
		&entities.UserModels{ID: 1, Name: "John Doe"},
		&entities.UserModels{ID: 2, Name: "Jane Smith"},
	}

	expectedTotalItems := int64(len(expectedUsers))

	service, repo, _ := setupTestService(t)

	t.Run("Success Case - Get All Users", func(t *testing.T) {
		repo.On("FindAll", page, perPage).Return(expectedUsers, nil).Once()
		repo.On("GetTotalUserCount").Return(expectedTotalItems, nil).Once()

		users, totalItems, err := service.GetAllUsers(page, perPage)

		assert.Nil(t, err)
		assert.Equal(t, expectedTotalItems, totalItems)
		assert.Equal(t, expectedUsers, users)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Error Getting All Users", func(t *testing.T) {
		expectedErr := errors.New("failed to get all users")
		repo.On("FindAll", page, perPage).Return(nil, expectedErr).Once()

		users, totalItems, err := service.GetAllUsers(page, perPage)

		assert.EqualError(t, err, expectedErr.Error())
		assert.Nil(t, users)
		assert.Zero(t, totalItems)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Error Getting Total User Count", func(t *testing.T) {
		expectedErr := errors.New("failed to get total user count")
		repo.On("FindAll", page, perPage).Return(expectedUsers, nil).Once()
		repo.On("GetTotalUserCount").Return(int64(0), expectedErr).Once()

		users, totalItems, err := service.GetAllUsers(page, perPage)

		assert.EqualError(t, err, expectedErr.Error())
		assert.Nil(t, users)
		assert.Zero(t, totalItems)
		repo.AssertExpectations(t)
	})
}

func TestUserService_ValidatePassword(t *testing.T) {
	userID := uint64(1)
	oldPassword := "oldPass123"
	newPassword := "newPass123"
	confirmPassword := "newPass123"

	service, repo, hash := setupTestService(t)

	t.Run("Success Case - Password Validation", func(t *testing.T) {
		storedPassword := "$2a$10$12345678901234567890123456789012345678901234567890"
		repo.On("GetUsersPassword", userID).Return(storedPassword, nil).Once()
		hash.On("ComparePassword", storedPassword, oldPassword).Return(true, nil).Once()

		err := service.ValidatePassword(userID, oldPassword, newPassword, confirmPassword)

		assert.Nil(t, err)
		repo.AssertExpectations(t)
		hash.AssertExpectations(t)
	})

	t.Run("Failed Case - Invalid Old Password", func(t *testing.T) {
		expectedErr := errors.New("Password lama tidak valid")
		repo.On("GetUsersPassword", userID).Return("", expectedErr).Once()

		err := service.ValidatePassword(userID, oldPassword, newPassword, confirmPassword)

		assert.EqualError(t, err, expectedErr.Error())
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Old Password Doesn't Match", func(t *testing.T) {
		expectedErr := errors.New("Password lama tidak valid")
		storedPassword := "$2a$10$12345678901234567890123456789012345678901234567890"
		repo.On("GetUsersPassword", userID).Return(storedPassword, nil).Once()
		hash.On("ComparePassword", storedPassword, oldPassword).Return(false, nil).Once()

		err := service.ValidatePassword(userID, oldPassword, newPassword, confirmPassword)

		assert.EqualError(t, err, expectedErr.Error())
		repo.AssertExpectations(t)
		hash.AssertExpectations(t)
	})

	t.Run("Failed Case - New Password Same as Old Password", func(t *testing.T) {
		expectedErr := errors.New("Password baru tidak boleh sama dengan password lama")
		err := service.ValidatePassword(userID, oldPassword, oldPassword, oldPassword)

		assert.EqualError(t, err, expectedErr.Error())
	})

	t.Run("Failed Case - New Password and Confirm Password Mismatch", func(t *testing.T) {
		expectedErr := errors.New("Password baru dan konfirmasi password tidak cocok")

		err := service.ValidatePassword(userID, oldPassword, newPassword, "mismatchedPass")

		assert.EqualError(t, err, expectedErr.Error())
	})
}
