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
		assert.Equal(t, "anda sudah submit tantangan ini sebelumnya", err.Error())

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
		assert.Equal(t, "tantangan sudah kadaluwarsa, tidak dapat submit", err.Error())
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

	challenge := &entities.ChallengeModels{
		ID:  1,
		Exp: 10,
	}

	user := &entities.UserModels{
		ID:             2,
		TotalChallenge: 20,
		Exp:            100,
		Level:          "Bronze",
	}
	form := &entities.ChallengeFormModels{
		ID:          1,
		UserID:      user.ID,
		ChallengeID: challenge.ID,
		Status:      "menunggu validasi",
		Exp:         challenge.Exp,
	}

	updatedData := dto.UpdateChallengeFormStatusRequest{
		Status: "valid",
	}

	var changeTotalChallenge int64

	switch {
	case form.Status == "valid" && updatedData.Status == "tidak valid":
		user.Exp -= form.Exp
		changeTotalChallenge = -1
	case form.Status == "tidak valid" && updatedData.Status == "valid":
		user.Exp += form.Exp
		changeTotalChallenge = 1
	case form.Status == "menunggu validasi" && updatedData.Status == "valid":
		user.Exp += form.Exp
		changeTotalChallenge = 1
	case form.Status == "valid" && updatedData.Status == "menunggu validasi":
		user.Exp -= form.Exp
		changeTotalChallenge = -1
	}

	user.TotalChallenge += uint64(changeTotalChallenge)

	t.Run("Success Case - Waiting Validation to Valid", func(t *testing.T) {
		updatedData := dto.UpdateChallengeFormStatusRequest{
			Status: "valid",
		}

		repo.On("GetSubmitChallengeFormById", form.ID).Return(form, nil)
		repoUser.On("GetUsersById", user.ID).Return(user, nil)
		user, err := userService.GetUsersById(form.UserID)
		assert.NoError(t, err)

		expectedExp := user.Exp + form.Exp
		expectedTotalChallenge := user.TotalChallenge + uint64(changeTotalChallenge)

		repo.On("UpdateSubmitChallengeForm", form.ID, updatedData).Return(form, nil)
		repoUser.On("UpdateUserExp", user.ID, user.Exp+form.Exp).Return(user, nil)
		repoUser.On("UpdateUserChallengeFollow", user.ID, user.TotalChallenge+uint64(changeTotalChallenge)).Return(user, nil)
		result, err := service.UpdateSubmitChallengeForm(form.ID, updatedData)
		assert.NoError(t, err)
		_, err = userService.UpdateUserExp(user.ID, user.Exp)
		assert.NoError(t, err)
		_, err = userService.UpdateUserChallengeFollow(user.ID, user.TotalChallenge)
		assert.NoError(t, err)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedExp, user.Exp)
		assert.Equal(t, expectedTotalChallenge, user.TotalChallenge)

		repo.AssertExpectations(t)
		repoUser.AssertExpectations(t)
	})

	t.Run("Success Case Valid to Not Valid", func(t *testing.T) {
		form2 := &entities.ChallengeFormModels{
			ID:          3,
			UserID:      user.ID,
			ChallengeID: 3,
			Status:      "valid",
			Exp:         20,
		}

		updatedData2 := dto.UpdateChallengeFormStatusRequest{
			Status: "tidak valid",
		}

		repo.On("GetSubmitChallengeFormById", form2.ID).Return(form2, nil)
		repoUser.On("GetUsersById", user.ID).Return(user, nil)
		user, err := userService.GetUsersById(form.UserID)
		assert.NoError(t, err)

		expectedExp := user.Exp - form2.Exp
		expectedTotalChallenge := user.TotalChallenge - uint64(changeTotalChallenge)

		repo.On("UpdateSubmitChallengeForm", form2.ID, updatedData2).Return(form2, nil)
		repoUser.On("UpdateUserExp", user.ID, expectedExp).Return(user, nil)
		repoUser.On("UpdateUserChallengeFollow", user.ID, expectedTotalChallenge).Return(user, nil)
		result, err := service.UpdateSubmitChallengeForm(form2.ID, updatedData2)
		assert.NoError(t, err)
		_, err = userService.UpdateUserExp(user.ID, user.Exp)
		assert.NoError(t, err)
		_, err = userService.UpdateUserChallengeFollow(user.ID, user.TotalChallenge)
		assert.NoError(t, err)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedExp, user.Exp)
		assert.Equal(t, expectedTotalChallenge, user.TotalChallenge)

		repo.AssertExpectations(t)
		repoUser.AssertExpectations(t)
	})

	t.Run("Success Case - Not Valid to Valid", func(t *testing.T) {

		form3 := &entities.ChallengeFormModels{
			ID:          4,
			UserID:      user.ID,
			ChallengeID: 4,
			Status:      "tidak valid",
			Exp:         20,
		}

		updatedData3 := dto.UpdateChallengeFormStatusRequest{
			Status: "valid",
		}

		repo.On("GetSubmitChallengeFormById", form3.ID).Return(form3, nil)
		repoUser.On("GetUsersById", user.ID).Return(user, nil)
		user, err := userService.GetUsersById(form.UserID)
		assert.NoError(t, err)

		expectedExp := user.Exp + form3.Exp
		expectedTotalChallenge := user.TotalChallenge + uint64(changeTotalChallenge)

		repo.On("UpdateSubmitChallengeForm", form3.ID, updatedData3).Return(form3, nil)
		repoUser.On("UpdateUserExp", user.ID, user.Exp+form3.Exp).Return(user, nil)
		repoUser.On("UpdateUserChallengeFollow", user.ID, user.TotalChallenge+uint64(changeTotalChallenge)).Return(user, nil)
		result, err := service.UpdateSubmitChallengeForm(form3.ID, updatedData3)
		assert.NoError(t, err)
		_, err = userService.UpdateUserExp(user.ID, user.Exp)
		assert.NoError(t, err)
		_, err = userService.UpdateUserChallengeFollow(user.ID, user.TotalChallenge)
		assert.NoError(t, err)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedExp, user.Exp)
		assert.Equal(t, expectedTotalChallenge, user.TotalChallenge)

		repo.AssertExpectations(t)
		repoUser.AssertExpectations(t)
	})

	t.Run("Failed Case - UpdateSubmitChallengeForm", func(t *testing.T) {
		repo := mocks.NewRepositoryChallengeInterface(t)
		repoUser := user_mock.NewRepositoryUserInterface(t)
		userService := user_service.NewUserService(repoUser, utils.NewHash())
		service := NewChallengeService(repo, userService)
		repo.On("GetSubmitChallengeFormById", form.ID).Return(form, nil)
		repoUser.On("GetUsersById", user.ID).Return(user, nil)
		repo.On("UpdateSubmitChallengeForm", form.ID, updatedData).Return(nil, errors.New("gagal memperbarui formulir"))
		result, err := service.UpdateSubmitChallengeForm(form.ID, updatedData)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, errors.New("gagal memperbarui formulir"), err)

		repo.AssertExpectations(t)
		repoUser.AssertExpectations(t)
	})

	t.Run("Failed Case - UpdateUserExp", func(t *testing.T) {
		repo.On("GetSubmitChallengeFormById", form.ID).Return(form, nil)
		repoUser.On("GetUsersById", user.ID).Return(user, nil)
		repo.On("UpdateSubmitChallengeForm", form.ID, updatedData).Return(form, nil)
		repoUser.On("UpdateUserExp", user.ID, user.Exp+form.Exp).Return(nil, errors.New("gagal menyimpan perubahan exp user ke database"))

		result, err := service.UpdateSubmitChallengeForm(form.ID, updatedData)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, errors.New("gagal menyimpan perubahan exp user ke database"), err)
		repo.AssertExpectations(t)
		repoUser.AssertExpectations(t)
	})

	t.Run("Failed Case - UpdateUserChallengeFollow", func(t *testing.T) {
		repo.On("GetSubmitChallengeFormById", form.ID).Return(form, nil)
		repoUser.On("GetUsersById", user.ID).Return(user, nil)
		repo.On("UpdateSubmitChallengeForm", form.ID, updatedData).Return(form, nil)
		repoUser.On("UpdateUserExp", user.ID, user.Exp+form.Exp).Return(user, nil)
		repoUser.On("UpdateUserChallengeFollow", user.ID, user.TotalChallenge+uint64(changeTotalChallenge)).Return(nil, errors.New("gagal menyimpan perubahan total challenge user ke database"))

		result, err := service.UpdateSubmitChallengeForm(form.ID, updatedData)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, errors.New("gagal menyimpan perubahan total tantangan user ke database"), err)
		repo.AssertExpectations(t)
		repoUser.AssertExpectations(t)
	})

	t.Run("Failed Case - Failed to Retrieve Form", func(t *testing.T) {
		form6 := &entities.ChallengeFormModels{
			ID:     10,
			Status: "menunggu validasi",
		}
		expectedErr := errors.New("formulir tidak ditemukan")
		repo.On("GetSubmitChallengeFormById", form6.ID).Return(nil, expectedErr).Once()

		result, err := service.UpdateSubmitChallengeForm(form6.ID, dto.UpdateChallengeFormStatusRequest{
			Status: "valid",
		})

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedErr, err)

		repo.AssertExpectations(t)
		repoUser.AssertExpectations(t)
	})

	t.Run("Failure Case - User Not Found", func(t *testing.T) {
		form7 := &entities.ChallengeFormModels{
			ID:     10,
			Status: "menunggu validasi",
		}
		repo.On("GetSubmitChallengeFormById", form7.ID).Return(form7, nil).Once()
		expectedErr := errors.New("pengguna tidak ada")
		repoUser.On("GetUsersById", form7.UserID).Return(nil, expectedErr).Once()

		result, err := service.UpdateSubmitChallengeForm(form7.ID, dto.UpdateChallengeFormStatusRequest{
			Status: "valid",
		})

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedErr, err)

		repo.AssertExpectations(t)
		repoUser.AssertExpectations(t)
	})

	t.Run("Success Case - Waiting Validation to Valid", func(t *testing.T) {
		form4 := &entities.ChallengeFormModels{
			ID:          5,
			UserID:      user.ID,
			ChallengeID: 5,
			Status:      "valid",
			Exp:         20,
		}

		updatedData4 := dto.UpdateChallengeFormStatusRequest{
			Status: "menunggu validasi",
		}

		repo.On("GetSubmitChallengeFormById", form4.ID).Return(form4, nil)
		repoUser.On("GetUsersById", user.ID).Return(user, nil)
		user, err := userService.GetUsersById(form.UserID)
		assert.NoError(t, err)

		expectedExp := user.Exp - form4.Exp
		expectedTotalChallenge := user.TotalChallenge - uint64(changeTotalChallenge)

		repo.On("UpdateSubmitChallengeForm", form4.ID, updatedData4).Return(form4, nil)
		repoUser.On("UpdateUserExp", user.ID, expectedExp).Return(user, nil)
		repoUser.On("UpdateUserChallengeFollow", user.ID, expectedTotalChallenge).Return(user, nil)
		result, err := service.UpdateSubmitChallengeForm(form4.ID, updatedData4)
		assert.NoError(t, err)
		_, err = userService.UpdateUserExp(user.ID, user.Exp)
		assert.NoError(t, err)
		_, err = userService.UpdateUserChallengeFollow(user.ID, user.TotalChallenge)
		assert.NoError(t, err)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedExp, user.Exp)
		assert.Equal(t, expectedTotalChallenge, user.TotalChallenge)

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
			ID: 1,
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
		assert.Equal(t, errors.New("formulir tidak ditemukan"), errNotFound)

		repo.AssertExpectations(t)
	})
}

func TestGetDatesFromFilterType(t *testing.T) {
	t.Run("Valid Filter Type - Hari Ini", func(t *testing.T) {
		filterType := "Hari Ini"
		startDate, endDate, err := getDatesFromFilterType(filterType)

		assert.Nil(t, err)
		assert.NotNil(t, startDate)
		assert.NotNil(t, endDate)
		assert.True(t, startDate.Before(endDate))
	})

	t.Run("Valid Filter Type - Minggu Ini", func(t *testing.T) {
		filterType := "Minggu Ini"
		startDate, endDate, err := getDatesFromFilterType(filterType)

		assert.Nil(t, err)
		assert.NotNil(t, startDate)
		assert.NotNil(t, endDate)
		assert.True(t, startDate.Before(endDate))
	})

	t.Run("Valid Filter Type - Bulan Ini", func(t *testing.T) {
		filterType := "Bulan Ini"
		startDate, endDate, err := getDatesFromFilterType(filterType)

		assert.Nil(t, err)
		assert.NotNil(t, startDate)
		assert.NotNil(t, endDate)
		assert.True(t, startDate.Before(endDate))
	})

	t.Run("Valid Filter Type - Tahun Ini", func(t *testing.T) {
		filterType := "Tahun Ini"
		startDate, endDate, err := getDatesFromFilterType(filterType)

		assert.Nil(t, err)
		assert.NotNil(t, startDate)
		assert.NotNil(t, endDate)
		assert.True(t, startDate.Before(endDate))
	})

	t.Run("Invalid Filter Type", func(t *testing.T) {
		filterType := "Invalid Type"
		startDate, endDate, err := getDatesFromFilterType(filterType)

		assert.Error(t, err)
		assert.Equal(t, time.Time{}, startDate)
		assert.Equal(t, time.Time{}, endDate)
		assert.Equal(t, errors.New("tipe filter tidak valid"), err)
	})
}

func TestChallengeService_GetSubmitChallengeFormByDateRange(t *testing.T) {
	repo := mocks.NewRepositoryChallengeInterface(t)
	userRepo := user_mock.NewRepositoryUserInterface(t)
	userService := user_service.NewUserService(userRepo, utils.NewHash())
	service := NewChallengeService(repo, userService)

	t.Run("Success Case", func(t *testing.T) {
		page := 1
		perPage := 10
		filterType := "Hari Ini"

		expectedForms := []*entities.ChallengeFormModels{}
		expectedTotalItems := int64(len(expectedForms))

		startDate, endDate, _ := getDatesFromFilterType(filterType)

		repo.On("GetSubmitChallengeFormByDateRange", page, perPage, startDate, endDate).Return(expectedForms, nil).Once()
		repo.On("GetTotalSubmitChallengeFormCountByDateRange", startDate, endDate).Return(expectedTotalItems, nil).Once()

		resultForms, resultTotalItems, err := service.GetSubmitChallengeFormByDateRange(page, perPage, filterType)

		assert.Nil(t, err)
		assert.Equal(t, expectedForms, resultForms)
		assert.Equal(t, expectedTotalItems, resultTotalItems)

		repo.AssertExpectations(t)
	})

	t.Run("Error Case - Failed to Get Dates from FilterType", func(t *testing.T) {
		page := 1
		perPage := 10
		filterType := "InvalidFilterType"
		_, _, err := service.GetSubmitChallengeFormByDateRange(page, perPage, filterType)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "tipe filter tidak valid")

	})

	t.Run("Error Case - GetSubmitChallengeFormByDateRange", func(t *testing.T) {
		page := 1
		perPage := 10
		filterType := "Hari Ini"

		startDate, endDate, _ := getDatesFromFilterType(filterType)

		repo.On("GetSubmitChallengeFormByDateRange", page, perPage, startDate, endDate).Return(nil, errors.New("database error")).Once()

		resultForms, resultTotalItems, err := service.GetSubmitChallengeFormByDateRange(page, perPage, filterType)

		assert.Error(t, err)
		assert.Nil(t, resultForms)
		assert.Zero(t, resultTotalItems)

		repo.AssertExpectations(t)
	})

	t.Run("Error Case - GetTotalSubmitChallengeFormCountByDateRange", func(t *testing.T) {
		page := 1
		perPage := 10
		filterType := "Hari Ini"

		expectedForms := []*entities.ChallengeFormModels{}

		startDate, endDate, _ := getDatesFromFilterType(filterType)

		repo.On("GetSubmitChallengeFormByDateRange", page, perPage, startDate, endDate).Return(expectedForms, nil).Once()
		repo.On("GetTotalSubmitChallengeFormCountByDateRange", startDate, endDate).Return(int64(0), errors.New("database error")).Once()

		resultForms, resultTotalItems, err := service.GetSubmitChallengeFormByDateRange(page, perPage, filterType)

		assert.Error(t, err)
		assert.Nil(t, resultForms)
		assert.Zero(t, resultTotalItems)

		repo.AssertExpectations(t)
	})
}

func TestChallengeService_GetSubmitChallengeFormByStatusAndDate(t *testing.T) {
	repo := mocks.NewRepositoryChallengeInterface(t)
	userRepo := user_mock.NewRepositoryUserInterface(t)
	userService := user_service.NewUserService(userRepo, utils.NewHash())
	service := NewChallengeService(repo, userService)

	t.Run("Success Case", func(t *testing.T) {
		page := 1
		perPage := 10
		filterStatus := "approved"
		filterType := "Hari Ini"

		expectedForms := []*entities.ChallengeFormModels{}
		expectedTotalItems := int64(len(expectedForms))

		startDate, endDate, _ := getDatesFromFilterType(filterType)

		repo.On("GetSubmitChallengeFormByStatusAndDate", page, perPage, filterStatus, startDate, endDate).Return(expectedForms, nil).Once()
		repo.On("GetTotalSubmitChallengeFormCountByStatusAndDate", filterStatus, startDate, endDate).Return(expectedTotalItems, nil).Once()

		resultForms, resultTotalItems, err := service.GetSubmitChallengeFormByStatusAndDate(page, perPage, filterStatus, filterType)

		assert.Nil(t, err)
		assert.Equal(t, expectedForms, resultForms)
		assert.Equal(t, expectedTotalItems, resultTotalItems)

		repo.AssertExpectations(t)
	})

	t.Run("Error Case - Failed to Get Dates from FilterType", func(t *testing.T) {
		page := 1
		perPage := 10
		filterStatus := "status"
		filterType := "kadaluwarsa"

		_, _, err := service.GetSubmitChallengeFormByStatusAndDate(page, perPage, filterStatus, filterType)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "tipe filter tidak valid")
	})

	t.Run("Error Case - GetSubmitChallengeFormByStatusAndDate", func(t *testing.T) {
		page := 1
		perPage := 10
		filterStatus := "approved"
		filterType := "Hari Ini"

		startDate, endDate, _ := getDatesFromFilterType(filterType)

		repo.On("GetSubmitChallengeFormByStatusAndDate", page, perPage, filterStatus, startDate, endDate).Return(nil, errors.New("database error")).Once()

		resultForms, resultTotalItems, err := service.GetSubmitChallengeFormByStatusAndDate(page, perPage, filterStatus, filterType)

		assert.Error(t, err)
		assert.Nil(t, resultForms)
		assert.Zero(t, resultTotalItems)

		repo.AssertExpectations(t)
	})

	t.Run("Error Case - GetTotalSubmitChallengeFormCountByStatusAndDate", func(t *testing.T) {
		page := 1
		perPage := 10
		filterStatus := "approved"
		filterType := "Hari Ini"

		expectedForms := []*entities.ChallengeFormModels{}

		startDate, endDate, _ := getDatesFromFilterType(filterType)

		repo.On("GetSubmitChallengeFormByStatusAndDate", page, perPage, filterStatus, startDate, endDate).Return(expectedForms, nil).Once()
		repo.On("GetTotalSubmitChallengeFormCountByStatusAndDate", filterStatus, startDate, endDate).Return(int64(0), errors.New("database error")).Once()

		resultForms, resultTotalItems, err := service.GetSubmitChallengeFormByStatusAndDate(page, perPage, filterStatus, filterType)

		assert.Error(t, err)
		assert.Nil(t, resultForms)
		assert.Zero(t, resultTotalItems)

		repo.AssertExpectations(t)
	})
}

func TestChallengeService_GetChallengesBySearchAndStatus(t *testing.T) {
	repo := mocks.NewRepositoryChallengeInterface(t)
	userRepo := user_mock.NewRepositoryUserInterface(t)
	userService := user_service.NewUserService(userRepo, utils.NewHash())
	service := NewChallengeService(repo, userService)

	t.Run("Success Case", func(t *testing.T) {
		page := 1
		perPage := 10
		search := "keyword"
		status := "active"

		expectedChallenges := []*entities.ChallengeModels{}
		expectedTotalItems := int64(len(expectedChallenges))

		repo.On("GetChallengesBySearchAndStatus", page, perPage, search, status).Return(expectedChallenges, expectedTotalItems, nil).Once()

		resultChallenges, resultTotalItems, err := service.GetChallengesBySearchAndStatus(page, perPage, search, status)

		assert.Nil(t, err)
		assert.Equal(t, expectedChallenges, resultChallenges)
		assert.Equal(t, expectedTotalItems, resultTotalItems)

		repo.AssertExpectations(t)
	})
	t.Run("Error Case - Failed to Get Challenges", func(t *testing.T) {
		page := 1
		perPage := 10
		search := "yourSearchQuery"
		status := "yourStatus"

		repo.On("GetChallengesBySearchAndStatus", page, perPage, search, status).Return(nil, int64(0), errors.New("database error")).Once()

		resultChallenges, resultTotalItems, err := service.GetChallengesBySearchAndStatus(page, perPage, search, status)

		assert.Error(t, err)
		assert.Nil(t, resultChallenges)
		assert.Zero(t, resultTotalItems)

		repo.AssertExpectations(t)
	})
}
