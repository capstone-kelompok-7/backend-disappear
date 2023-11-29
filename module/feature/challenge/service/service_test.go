package service

import (
	"errors"
	"testing"
	"time"

	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/challenge/dto"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/challenge/mocks"
	user_mock "github.com/capstone-kelompok-7/backend-disappear/module/feature/users/mocks"
	user_service "github.com/capstone-kelompok-7/backend-disappear/module/feature/users/service"
	"github.com/capstone-kelompok-7/backend-disappear/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestChallengeService_CalculatePaginationValues(t *testing.T) {
	service := &ChallengeService{}

	t.Run("Page less than or equal to zero should default to 1", func(t *testing.T) {
		page, totalPages := service.CalculatePaginationValues(0, 100, 8)

		assert.Equal(t, 1, page)
		assert.Equal(t, 13, totalPages)
	})

	t.Run("Page exceeds total pages should set to total pages", func(t *testing.T) {
		page, totalPages := service.CalculatePaginationValues(15, 100, 8)

		assert.Equal(t, 13, page)
		assert.Equal(t, 13, totalPages)
	})

	t.Run("Page within limits should return correct values", func(t *testing.T) {
		page, totalPages := service.CalculatePaginationValues(2, 100, 8)

		assert.Equal(t, 2, page)
		assert.Equal(t, 13, totalPages)
	})

	t.Run("Total items not perfectly divisible by perPage should round totalPages", func(t *testing.T) {
		page, totalPages := service.CalculatePaginationValues(1, 95, 8)

		assert.Equal(t, 1, page)
		assert.Equal(t, 12, totalPages)
	})
}

func TestChallengeService_GetNextPage(t *testing.T) {
	service := &ChallengeService{}

	t.Run("Next Page Within Total Pages", func(t *testing.T) {
		currentPage := 3
		totalPages := 5

		nextPage := service.GetNextPage(currentPage, totalPages)

		assert.Equal(t, currentPage+1, nextPage)
	})

	t.Run("Next Page Equal to Total Pages", func(t *testing.T) {
		currentPage := 5
		totalPages := 5

		nextPage := service.GetNextPage(currentPage, totalPages)

		assert.Equal(t, totalPages, nextPage)
	})
}

func TestChallengelService_GetPrevPage(t *testing.T) {
	service := &ChallengeService{}

	t.Run("Previous Page Within Bounds", func(t *testing.T) {
		currentPage := 3

		prevPage := service.GetPrevPage(currentPage)

		assert.Equal(t, currentPage-1, prevPage)
	})

	t.Run("Previous Page at Lower Bound", func(t *testing.T) {
		currentPage := 1

		prevPage := service.GetPrevPage(currentPage)

		assert.Equal(t, currentPage, prevPage)
	})
}

func TestChallengeService_GetAll(t *testing.T) {
	repo := mocks.NewRepositoryChallengeInterface(t)
	repoUser := user_mock.NewRepositoryUserInterface(t)
	userService := user_service.NewUserService(repoUser, utils.NewHash())

	service := NewChallengeService(repo, userService)

	challenges := []*entities.ChallengeModels{
		{ID: 1, Title: "Challenge 1", Photo: "challenge1.jpg", StartDate: time.Now(), EndDate: time.Now().AddDate(0, 0, 7), Description: "Description 1", Status: "Status 1", Exp: 500},
		{ID: 2, Title: "Challenge 2", Photo: "challenge2.jpg", StartDate: time.Now(), EndDate: time.Now().AddDate(0, 0, 7), Description: "Description 2", Status: "Status 2", Exp: 200},
	}

	t.Run("Success Case - Challenge Found", func(t *testing.T) {
		expectedTotalItems := int64(10)
		repo.On("FindAll", 1, 10).Return(challenges, nil).Once()
		repo.On("GetTotalChallengeCount").Return(expectedTotalItems, nil).Once()

		result, totalItems, err := service.GetAllChallenges(1, 10)

		assert.NoError(t, err)
		assert.Equal(t, len(challenges), len(result))
		assert.Equal(t, expectedTotalItems, totalItems)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - GetTotalChallengeCount Error", func(t *testing.T) {
		expectedErr := errors.New("GetTotalChallengeCount Error")

		repo.On("FindAll", 1, 10).Return(challenges, nil).Once()
		repo.On("GetTotalChallengeCount").Return(int64(0), expectedErr).Once()

		result, totalItems, err := service.GetAllChallenges(1, 10)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, int64(0), totalItems)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - GetAll Error", func(t *testing.T) {
		expectedErr := errors.New("FindAll Error")
		repo.On("FindAll", 1, 10).Return(nil, expectedErr).Once()

		result, totalItems, err := service.GetAllChallenges(1, 10)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, int64(0), totalItems)
		repo.AssertExpectations(t)
	})
}

func TestChallengeService_GetChallengeByTitle(t *testing.T) {
	repo := mocks.NewRepositoryChallengeInterface(t)
	repoUser := user_mock.NewRepositoryUserInterface(t)
	userService := user_service.NewUserService(repoUser, utils.NewHash())

	service := NewChallengeService(repo, userService)

	challenges := []*entities.ChallengeModels{
		{ID: 1, Title: "Challenge 1", Photo: "challenge1.jpg", StartDate: time.Now(), EndDate: time.Now().AddDate(0, 0, 7), Description: "Description 1", Status: "Status 1", Exp: 500},
		{ID: 2, Title: "Challenge 2", Photo: "challenge2.jpg", StartDate: time.Now(), EndDate: time.Now().AddDate(0, 0, 7), Description: "Description 2", Status: "Status 2", Exp: 200},
	}
	title := "Test"

	t.Run("Success Case - Challenge Found by Name", func(t *testing.T) {
		expectedTotalItems := int64(10)
		repo.On("FindByTitle", 1, 10, title).Return(challenges, nil).Once()
		repo.On("GetTotalChallengeCountByTitle", title).Return(expectedTotalItems, nil).Once()

		result, totalItems, err := service.GetChallengeByTitle(1, 10, title)

		assert.NoError(t, err)
		assert.Equal(t, len(challenges), len(result))
		assert.Equal(t, expectedTotalItems, totalItems)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Error Finding Challenge by Title", func(t *testing.T) {
		expectedErr := errors.New("failed to find challenge by tittle")
		repo.On("FindByTitle", 1, 10, title).Return(nil, expectedErr).Once()

		result, totalItems, err := service.GetChallengeByTitle(1, 10, title)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Zero(t, totalItems)
		assert.Equal(t, expectedErr, err)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Error Getting Total Challenge Count by Title", func(t *testing.T) {
		expectedErrSubstring := "failed to get total challenge count by title"
		repo.On("FindByTitle", 1, 10, title).Return(challenges, nil).Once()
		repo.On("GetTotalChallengeCountByTitle", title).Return(int64(0), errors.New(expectedErrSubstring)).Once()

		result, totalItems, err := service.GetChallengeByTitle(1, 10, title)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, int64(0), totalItems)
		repo.AssertExpectations(t)
	})
}

func TestChallengeService_GetChallengeByStatus(t *testing.T) {
	repo := mocks.NewRepositoryChallengeInterface(t)
	repoUser := user_mock.NewRepositoryUserInterface(t)
	userService := user_service.NewUserService(repoUser, utils.NewHash())

	service := NewChallengeService(repo, userService)

	challenges := []*entities.ChallengeModels{
		{ID: 1, Title: "Challenge 1", Photo: "challenge1.jpg", StartDate: time.Now(), EndDate: time.Now().AddDate(0, 0, 7), Description: "Description 1", Status: "Status 1", Exp: 500},
		{ID: 2, Title: "Challenge 2", Photo: "challenge2.jpg", StartDate: time.Now(), EndDate: time.Now().AddDate(0, 0, 7), Description: "Description 2", Status: "Status 2", Exp: 200},
	}
	status := "Test"

	t.Run("Success Case - Challenge Found by Status", func(t *testing.T) {
		expectedTotalItems := int64(10)
		repo.On("FindByStatus", 1, 10, status).Return(challenges, nil).Once()
		repo.On("GetTotalChallengeCountByStatus", status).Return(expectedTotalItems, nil).Once()

		result, totalItems, err := service.GetChallengeByStatus(1, 10, status)

		assert.NoError(t, err)
		assert.Equal(t, len(challenges), len(result))
		assert.Equal(t, expectedTotalItems, totalItems)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Error Finding Challenge by Status", func(t *testing.T) {
		expectedErr := errors.New("failed to find challenge by status")
		repo.On("FindByStatus", 1, 10, status).Return(nil, expectedErr).Once()

		result, totalItems, err := service.GetChallengeByStatus(1, 10, status)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Zero(t, totalItems)
		assert.Equal(t, expectedErr, err)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Error Getting Total Challenge Count by Status", func(t *testing.T) {
		expectedErrSubstring := "failed to get total challenge count by status"
		repo.On("FindByStatus", 1, 10, status).Return(challenges, nil).Once()
		repo.On("GetTotalChallengeCountByStatus", status).Return(int64(0), errors.New(expectedErrSubstring)).Once()

		result, totalItems, err := service.GetChallengeByStatus(1, 10, status)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, int64(0), totalItems)
		repo.AssertExpectations(t)
	})
}

func TestChallengeService_CreateChallenge(t *testing.T) {
	repo := mocks.NewRepositoryChallengeInterface(t)
	repoUser := user_mock.NewRepositoryUserInterface(t)
	userService := user_service.NewUserService(repoUser, utils.NewHash())
	service := NewChallengeService(repo, userService)

	t.Run("Success Case - Belum Kadaluwarsa", func(t *testing.T) {
		challenges := &entities.ChallengeModels{
			Title:       "Challenge 1",
			Photo:       "challenge1.jpg",
			StartDate:   time.Now(),
			EndDate:     time.Now().AddDate(0, 0, 7),
			Description: "Description 1",
			Exp:         500,
		}

		expectedStatusBelumKadaluwarsa := "Belum Kadaluwarsa"

		repo.On("CreateChallenge", mock.AnythingOfType("*entities.ChallengeModels")).Run(func(args mock.Arguments) {
			challengeArg := args.Get(0).(*entities.ChallengeModels)
			assert.Equal(t, expectedStatusBelumKadaluwarsa, challengeArg.Status)
		}).Return(&entities.ChallengeModels{Status: expectedStatusBelumKadaluwarsa}, nil).Once()

		resultBelumKadaluwarsa, errBelumKadaluwarsa := service.CreateChallenge(challenges)

		assert.Nil(t, errBelumKadaluwarsa)
		assert.NotNil(t, resultBelumKadaluwarsa)
		assert.Equal(t, expectedStatusBelumKadaluwarsa, resultBelumKadaluwarsa.Status)

		repo.AssertExpectations(t)
	})

	t.Run("Success Case - Kadaluwarsa", func(t *testing.T) {
		challengesKadaluwarsa := &entities.ChallengeModels{
			Title:       "Challenge 2",
			Photo:       "challenge2.jpg",
			StartDate:   time.Now().AddDate(0, 0, -7),
			EndDate:     time.Now().AddDate(0, 0, -1),
			Description: "Description 2",
			Exp:         800,
		}

		expectedStatusKadaluwarsa := "Kadaluwarsa"

		repo.On("CreateChallenge", mock.AnythingOfType("*entities.ChallengeModels")).Run(func(args mock.Arguments) {
			challengeArg := args.Get(0).(*entities.ChallengeModels)
			assert.Equal(t, expectedStatusKadaluwarsa, challengeArg.Status)
		}).Return(&entities.ChallengeModels{Status: expectedStatusKadaluwarsa}, nil).Once()

		resultKadaluwarsa, errKadaluwarsa := service.CreateChallenge(challengesKadaluwarsa)

		assert.Nil(t, errKadaluwarsa)
		assert.NotNil(t, resultKadaluwarsa)
		assert.Equal(t, expectedStatusKadaluwarsa, resultKadaluwarsa.Status)

		repo.AssertExpectations(t)
	})

	t.Run("Failed Case", func(t *testing.T) {
		challengesFailed := &entities.ChallengeModels{
			Title: "Challenge Failed",
		}

		expectedError := errors.New("gagal menambahkan tantangan")

		repo.On("CreateChallenge", mock.AnythingOfType("*entities.ChallengeModels")).Return(nil, expectedError).Once()

		resultFailed, errFailed := service.CreateChallenge(challengesFailed)

		assert.Error(t, errFailed)
		assert.Nil(t, resultFailed)
		assert.Equal(t, expectedError, errFailed)

		repo.AssertExpectations(t)
	})
}

func TestChallengeService_GetChallengeById(t *testing.T) {
	repo := mocks.NewRepositoryChallengeInterface(t)
	repoUser := user_mock.NewRepositoryUserInterface(t)
	userService := user_service.NewUserService(repoUser, utils.NewHash())
	service := NewChallengeService(repo, userService)

	t.Run("Success Case - Found", func(t *testing.T) {
		expectedChallenge := &entities.ChallengeModels{
			ID:          1,
			Title:       "Challenge 1",
			Photo:       "challenge1.jpg",
			StartDate:   time.Now(),
			EndDate:     time.Now().AddDate(0, 0, 7),
			Description: "Description 1",
			Status:      "Belum Kadaluwarsa",
			Exp:         500,
		}

		repo.On("GetChallengeById", uint64(1)).Return(expectedChallenge, nil).Once()

		resultFound, errFound := service.GetChallengeById(uint64(1))

		assert.Nil(t, errFound)
		assert.NotNil(t, resultFound)
		assert.Equal(t, expectedChallenge, resultFound)

		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Not Found", func(t *testing.T) {
		repo.On("GetChallengeById", uint64(2)).Return(nil, errors.New("not found")).Once()

		resultNotFound, errNotFound := service.GetChallengeById(uint64(2))

		assert.Error(t, errNotFound)
		assert.Nil(t, resultNotFound)
		assert.Equal(t, errors.New("tantangan tidak ditemukan"), errNotFound)

		repo.AssertExpectations(t)
	})
}

func TestUpdateIfNotEmpty(t *testing.T) {
	t.Run("Non-empty value", func(t *testing.T) {
		target := "initial"
		value := "updated"
		updateIfNotEmpty(&target, value)
		assert.Equal(t, value, target)
	})

	t.Run("Empty value", func(t *testing.T) {
		target := "initial"
		value := ""
		updateIfNotEmpty(&target, value)
		assert.Equal(t, "initial", target)
	})
}

func TestUpdateIfNotZero(t *testing.T) {
	t.Run("Non-zero value", func(t *testing.T) {
		target := time.Now()
		value, err := time.Parse(time.RFC3339, "2023-11-24T19:05:45+07:00")
		assert.NoError(t, err)
		updateIfNotZero(&target, value)
		assert.True(t, target.Equal(value), "Expected time to be equal")
	})

	t.Run("Zero value", func(t *testing.T) {
		target := time.Now()
		value := time.Time{}
		updateIfNotZero(&target, value)
	})
}

func TestUpdateIfNotZeroUint64(t *testing.T) {
	t.Run("Non-zero value", func(t *testing.T) {
		var target uint64 = 42
		var value uint64 = 99
		updateIfNotZeroUint64(&target, value)
		assert.Equal(t, uint64(99), target)
	})

	t.Run("Zero value", func(t *testing.T) {
		var target uint64 = 42
		var value uint64 = 0
		updateIfNotZeroUint64(&target, value)
		assert.Equal(t, uint64(42), target)
	})
}

func TestChallengeService_UpdateChallenge(t *testing.T) {
	repo := mocks.NewRepositoryChallengeInterface(t)
	repoUser := user_mock.NewRepositoryUserInterface(t)
	userService := user_service.NewUserService(repoUser, utils.NewHash())
	service := NewChallengeService(repo, userService)

	t.Run("Status Change: Kadaluwarsa to Belum Kadaluwarsa", func(t *testing.T) {
		existingChallenge := &entities.ChallengeModels{
			ID:          1,
			Title:       "Existing Challenge",
			StartDate:   time.Now().AddDate(0, 0, -14),
			EndDate:     time.Now().AddDate(0, 0, -7),
			Description: "Existing Description",
			Exp:         100,
			Status:      "Kadaluwarsa",
		}

		updateData := &entities.ChallengeModels{
			EndDate: time.Now().AddDate(0, 0, 7),
		}

		repo.On("GetChallengeById", uint64(1)).Return(existingChallenge, nil).Once()
		repo.On("UpdateChallenge", uint64(1), mock.AnythingOfType("*entities.ChallengeModels")).Return(existingChallenge, nil).Once()

		updatedChallenge, err := service.UpdateChallenge(uint64(1), updateData)
		assert.Nil(t, err)
		assert.Equal(t, updateData.EndDate, updatedChallenge.EndDate)
		assert.Equal(t, "Belum Kadaluwarsa", updatedChallenge.Status)
		repo.AssertExpectations(t)
	})

	t.Run("Title Update", func(t *testing.T) {
		existingChallenge := &entities.ChallengeModels{
			ID:    1,
			Title: "Existing Challenge",
		}

		updateData := &entities.ChallengeModels{
			Title: "Updated Challenge",
		}

		repo.On("GetChallengeById", uint64(1)).Return(existingChallenge, nil).Once()
		repo.On("UpdateChallenge", uint64(1), mock.AnythingOfType("*entities.ChallengeModels")).Return(existingChallenge, nil).Once()

		updatedChallenge, err := service.UpdateChallenge(uint64(1), updateData)
		assert.Nil(t, err)
		assert.Equal(t, updateData.Title, updatedChallenge.Title)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case: GetChallengeById Error", func(t *testing.T) {
		expectedErr := errors.New("GetChallengeById failed")
		repo.On("GetChallengeById", uint64(1)).Return(nil, expectedErr).Once()

		updateData := &entities.ChallengeModels{}

		result, err := service.UpdateChallenge(uint64(1), updateData)

		assert.Error(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedErr, err)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case: UpdateChallenge Error", func(t *testing.T) {
		expectedUpdateErr := errors.New("UpdateChallenge failed")
		existingChallenge := &entities.ChallengeModels{
			ID:    1,
			Title: "Existing Challenge",
		}
		repo.On("GetChallengeById", uint64(1)).Return(existingChallenge, nil).Once()
		repo.On("UpdateChallenge", uint64(1), mock.AnythingOfType("*entities.ChallengeModels")).Return(nil, expectedUpdateErr).Once()

		updateData := &entities.ChallengeModels{}

		result, err := service.UpdateChallenge(uint64(1), updateData)

		assert.Error(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedUpdateErr, err)
		repo.AssertExpectations(t)
	})
}

func TestChallengeService_DeleteChallenge(t *testing.T) {
	repo := mocks.NewRepositoryChallengeInterface(t)
	repoUser := user_mock.NewRepositoryUserInterface(t)
	userService := user_service.NewUserService(repoUser, utils.NewHash())
	service := NewChallengeService(repo, userService)

	existingChallenge := &entities.ChallengeModels{
		ID:          1,
		Title:       "Existing Challenge",
		StartDate:   time.Now().AddDate(0, 0, -14),
		EndDate:     time.Now().AddDate(0, 0, -7),
		Description: "Existing Description",
		Exp:         100,
		Status:      "Belum Kadaluwarsa",
	}

	t.Run("Success Case", func(t *testing.T) {
		repo.On("GetChallengeById", existingChallenge.ID).Return(existingChallenge, nil).Once()
		repo.On("DeleteChallenge", existingChallenge.ID).Return(nil).Once()

		err := service.DeleteChallenge(existingChallenge.ID)

		assert.Nil(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Challenge Not Found", func(t *testing.T) {
		expectedErr := errors.New("tantangan tidak ditemukan")
		repo.On("GetChallengeById", existingChallenge.ID).Return(nil, expectedErr).Once()

		result := service.DeleteChallenge(existingChallenge.ID)

		assert.Error(t, result)
		assert.Equal(t, expectedErr, result)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Challenge Found But Delete Failed", func(t *testing.T) {
		expectedDeleteErr := errors.New("gagal menghapus tantangan")
		repo.On("GetChallengeById", existingChallenge.ID).Return(existingChallenge, nil).Once()
		repo.On("DeleteChallenge", existingChallenge.ID).Return(expectedDeleteErr).Once()

		result := service.DeleteChallenge(existingChallenge.ID)

		assert.Error(t, result)
		assert.Equal(t, expectedDeleteErr, result)
		repo.AssertExpectations(t)
	})
}

func TestChallengeService_CreateSubmitChallengeForm(t *testing.T) {
	repo := mocks.NewRepositoryChallengeInterface(t)
	repoUser := user_mock.NewRepositoryUserInterface(t)
	userService := user_service.NewUserService(repoUser, utils.NewHash())
	service := NewChallengeService(repo, userService)

	existingChallenge := &entities.ChallengeModels{
		ID:          1,
		Title:       "Existing Challenge",
		StartDate:   time.Now().AddDate(0, 0, -14),
		EndDate:     time.Now().AddDate(0, 0, -7),
		Description: "Existing Description",
		Exp:         100,
		Status:      "Belum Kadaluwarsa",
	}

	form := &entities.ChallengeFormModels{
		UserID:      1,
		ChallengeID: 1,
		Username:    "user123",
		Photo:       "user123.jpg",
		Status:      "menunggu validasi",
		Exp:         100,
		CreatedAt:   time.Now(),
	}

	t.Run("Succes Case", func(t *testing.T) {
		repo.On("GetSubmitChallengeFormByUserAndChallenge", form.UserID).Return([]*entities.ChallengeFormModels{}, nil).Once()
		repo.On("GetChallengeById", form.ChallengeID).Return(existingChallenge, nil).Once()

		repo.On("CreateSubmitChallengeForm", mock.AnythingOfType("*entities.ChallengeFormModels")).Return(form, nil).Once()
		result, err := service.CreateSubmitChallengeForm(form)

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, form, result)

		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - CreateSubmitChallengeForm Error", func(t *testing.T) {
		repo.On("GetSubmitChallengeFormByUserAndChallenge", form.UserID).Return([]*entities.ChallengeFormModels{}, nil).Once()
		repo.On("GetChallengeById", form.ChallengeID).Return(existingChallenge, nil).Once()

		expectedCreateErr := errors.New("Gagal submit challenge")
		repo.On("CreateSubmitChallengeForm", mock.AnythingOfType("*entities.ChallengeFormModels")).Return(nil, expectedCreateErr).Once()
		result, err := service.CreateSubmitChallengeForm(form)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedCreateErr, err)

		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - GetChallengeById Error", func(t *testing.T) {
		repo.On("GetSubmitChallengeFormByUserAndChallenge", form.UserID).Return([]*entities.ChallengeFormModels{}, nil).Once()

		expectedGetChallengeErr := errors.New("Gagal mendapatkan challenge by ID")
		repo.On("GetChallengeById", form.ChallengeID).Return(nil, expectedGetChallengeErr).Once()
		result, err := service.CreateSubmitChallengeForm(form)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedGetChallengeErr, err)

		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Already Submitted and Error", func(t *testing.T) {
		repo.On("GetSubmitChallengeFormByUserAndChallenge", form.UserID).Return([]*entities.ChallengeFormModels{form}, nil).Once()
		result, err := service.CreateSubmitChallengeForm(form)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "Anda sudah submit challenge ini sebelumnya", err.Error())

		expectedErr := errors.New("Anda sudah submit challenge ini sebelumnya")
		repo.On("GetSubmitChallengeFormByUserAndChallenge", form.UserID).Return(nil, expectedErr).Once()
		result, err = service.CreateSubmitChallengeForm(form)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedErr, err)

		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Expired Challenge", func(t *testing.T) {
		expiredChallenge := &entities.ChallengeModels{
			ID:          2,
			Title:       "Expired Challenge",
			StartDate:   time.Now().AddDate(0, 0, -21),
			EndDate:     time.Now().AddDate(0, 0, -14),
			Description: "Expired Challenge Description",
			Exp:         150,
			Status:      "Kadaluwarsa",
		}

		repo.On("GetSubmitChallengeFormByUserAndChallenge", form.UserID).Return([]*entities.ChallengeFormModels{}, nil).Once()
		repo.On("GetChallengeById", form.ChallengeID).Return(expiredChallenge, nil).Once()

		result, err := service.CreateSubmitChallengeForm(form)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "Tantangan sudah kadaluwarsa, tidak dapat submit", err.Error())
		repo.AssertExpectations(t)
	})
}

func TestChallengeService_GetAllForm(t *testing.T) {
	repo := mocks.NewRepositoryChallengeInterface(t)
	repoUser := user_mock.NewRepositoryUserInterface(t)
	userService := user_service.NewUserService(repoUser, utils.NewHash())

	service := NewChallengeService(repo, userService)

	formChallenge := []*entities.ChallengeFormModels{
		{ID: 1, UserID: 1, ChallengeID: 1, Username: "user123", Photo: "user123.jpg", Status: "menunggu validasi", Exp: 100, CreatedAt: time.Now()},
		{ID: 2, UserID: 1, ChallengeID: 2, Username: "user123", Photo: "user123.jpg", Status: "menunggu validasi", Exp: 50, CreatedAt: time.Now()},
	}

	t.Run("Success Case - Form Challenge Found", func(t *testing.T) {
		expectedTotalItems := int64(10)
		repo.On("GetAllSubmitChallengeForm", 1, 10).Return(formChallenge, nil).Once()
		repo.On("GetTotalSubmitChallengeFormCount").Return(expectedTotalItems, nil).Once()

		result, totalItems, err := service.GetAllSubmitChallengeForm(1, 10)

		assert.NoError(t, err)
		assert.Equal(t, len(formChallenge), len(result))
		assert.Equal(t, expectedTotalItems, totalItems)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - GetTotalSubmitChallengeFormCount Error", func(t *testing.T) {
		expectedErr := errors.New("GetTotalSubmitChallengeFormCount Error")

		repo.On("GetAllSubmitChallengeForm", 1, 10).Return(formChallenge, nil).Once()
		repo.On("GetTotalSubmitChallengeFormCount").Return(int64(0), expectedErr).Once()

		result, totalItems, err := service.GetAllSubmitChallengeForm(1, 10)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, int64(0), totalItems)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - GetAllSubmitChallengeForm Error", func(t *testing.T) {
		expectedErr := errors.New("GetAllSubmitChallengeForm Error")
		repo.On("GetAllSubmitChallengeForm", 1, 10).Return(nil, expectedErr).Once()

		result, totalItems, err := service.GetAllSubmitChallengeForm(1, 10)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, int64(0), totalItems)
		repo.AssertExpectations(t)
	})
}

func TestChallengeService_GetChallengeFormByStatus(t *testing.T) {
	repo := mocks.NewRepositoryChallengeInterface(t)
	repoUser := user_mock.NewRepositoryUserInterface(t)
	userService := user_service.NewUserService(repoUser, utils.NewHash())

	service := NewChallengeService(repo, userService)

	formChallenge := []*entities.ChallengeFormModels{
		{ID: 1, UserID: 1, ChallengeID: 1, Username: "user123", Photo: "user123.jpg", Status: "menunggu validasi", Exp: 100, CreatedAt: time.Now()},
		{ID: 2, UserID: 1, ChallengeID: 2, Username: "user123", Photo: "user123.jpg", Status: "menunggu validasi", Exp: 50, CreatedAt: time.Now()},
	}

	status := "Test"

	t.Run("Success Case - Challenge Form Found by Status", func(t *testing.T) {
		expectedTotalItems := int64(10)
		repo.On("GetSubmitChallengeFormByStatus", 1, 10, status).Return(formChallenge, nil).Once()
		repo.On("GetTotalSubmitChallengeFormCountByStatus", status).Return(expectedTotalItems, nil).Once()

		result, totalItems, err := service.GetSubmitChallengeFormByStatus(1, 10, status)

		assert.NoError(t, err)
		assert.Equal(t, len(formChallenge), len(result))
		assert.Equal(t, expectedTotalItems, totalItems)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Error Finding Challenge Form by Status", func(t *testing.T) {
		expectedErr := errors.New("failed to find challenge form by status")
		repo.On("GetSubmitChallengeFormByStatus", 1, 10, status).Return(nil, expectedErr).Once()

		result, totalItems, err := service.GetSubmitChallengeFormByStatus(1, 10, status)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Zero(t, totalItems)
		assert.Equal(t, expectedErr, err)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Error Getting Total Challenge Form Count by Status", func(t *testing.T) {
		expectedErrSubstring := "failed to get total challenge form count by status"
		repo.On("GetSubmitChallengeFormByStatus", 1, 10, status).Return(formChallenge, nil).Once()
		repo.On("GetTotalSubmitChallengeFormCountByStatus", status).Return(int64(0), errors.New(expectedErrSubstring)).Once()

		result, totalItems, err := service.GetSubmitChallengeFormByStatus(1, 10, status)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, int64(0), totalItems)
		repo.AssertExpectations(t)
	})
}
func TestChallengeService_UpdateSubmitChallengeFormm(t *testing.T) {
	repo := mocks.NewRepositoryChallengeInterface(t)
	repoUser := user_mock.NewRepositoryUserInterface(t)
	userService := user_service.NewUserService(repoUser, utils.NewHash())
	service := NewChallengeService(repo, userService)

	form := &entities.ChallengeFormModels{
		ID:          1,
		UserID:      1,
		ChallengeID: 1,
		Username:    "user123",
		Photo:       "user123.jpg",
		Status:      "menunggu validasi",
		Exp:         100,
		CreatedAt:   time.Now().AddDate(0, 0, 7),
	}

	t.Run("Gagal Mendapatkan Formulir", func(t *testing.T) {
		expectedErr := errors.New("Form tidak ditemukan")
		repo.On("GetSubmitChallengeFormById", form.ID).Return(nil, expectedErr).Once()

		result, err := service.UpdateSubmitChallengeForm(form.ID, dto.UpdateChallengeFormStatusRequest{
			Status: "valid",
		})

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedErr, err)

		repo.AssertExpectations(t)
		repoUser.AssertExpectations(t)
	})

	t.Run("Gagal Mendapatkan User", func(t *testing.T) {
		repo.On("GetSubmitChallengeFormById", form.ID).Return(form, nil).Once()
		expectedErr := errors.New("Gagal mendapatkan data user")
		repoUser.On("GetUsersById", form.UserID).Return(nil, expectedErr).Once()

		result, err := service.UpdateSubmitChallengeForm(form.ID, dto.UpdateChallengeFormStatusRequest{
			Status: "valid",
		})

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedErr, err)

		repo.AssertExpectations(t)
		repoUser.AssertExpectations(t)
	})
}

func TestChallengeService_GetChallengeFormById(t *testing.T) {
	repo := mocks.NewRepositoryChallengeInterface(t)
	repoUser := user_mock.NewRepositoryUserInterface(t)
	userService := user_service.NewUserService(repoUser, utils.NewHash())
	service := NewChallengeService(repo, userService)

	t.Run("Success Case - Found", func(t *testing.T) {
		expectedForm := &entities.ChallengeFormModels{
			ID:          1,
			UserID:      1,
			ChallengeID: 1,
			Username:    "user123",
			Photo:       "user123.jpg",
			Status:      "menunggu validasi",
			Exp:         100,
			CreatedAt:   time.Now().AddDate(0, 0, 7),
		}

		repo.On("GetSubmitChallengeFormById", uint64(1)).Return(expectedForm, nil).Once()

		resultFound, errFound := service.GetSubmitChallengeFormById(uint64(1))

		assert.Nil(t, errFound)
		assert.NotNil(t, resultFound)
		assert.Equal(t, expectedForm, resultFound)

		repo.AssertExpectations(t)
	})

	t.Run("Failed Case - Not Found", func(t *testing.T) {
		repo.On("GetSubmitChallengeFormById", uint64(2)).Return(nil, errors.New("not found")).Once()

		resultNotFound, errNotFound := service.GetSubmitChallengeFormById(uint64(2))

		assert.Error(t, errNotFound)
		assert.Nil(t, resultNotFound)
		assert.Equal(t, errors.New("form tidak ditemukan"), errNotFound)

		repo.AssertExpectations(t)
	})
}
