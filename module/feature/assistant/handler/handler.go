package handler

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/assistant"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/assistant/dto"
	"github.com/capstone-kelompok-7/backend-disappear/utils/response"
	"github.com/labstack/echo/v4"
)

type AssistantHandler struct {
	service assistant.ServiceAssistantInterface
}

func NewAssistantHandler(service assistant.ServiceAssistantInterface) assistant.HandlerAssistantInterface {
	return &AssistantHandler{
		service: service,
	}
}

func (h *AssistantHandler) CreateQuestion() echo.HandlerFunc {
	return func(c echo.Context) error {
		currentUser := c.Get("CurrentUser").(*entities.UserModels)
		if currentUser.Role != "customer" {
			return response.SendStatusForbiddenResponse(c, "Tidak diizinkan: Anda tidak memiliki izin")
		}

		chatRequest := new(dto.CreateChatRequest)
		if err := c.Bind(chatRequest); err != nil {
			return response.SendBadRequestResponse(c, "Format input yang Anda masukkan tidak sesuai.")
		}

		newUser := &entities.ChatModel{
			Text: chatRequest.Text,
		}

		err := h.service.CreateQuestion(currentUser.ID, *newUser)
		if err != nil {
			return response.SendStatusInternalServerResponse(c, "Gagal membuat chat: "+err.Error())
		}

		return response.SendSuccessResponse(c, "Berhasil membuat chat", chatRequest.Text)
	}
}

func (h *AssistantHandler) CreateAnswer() echo.HandlerFunc {
	return func(c echo.Context) error {
		currentUser := c.Get("CurrentUser").(*entities.UserModels)
		if currentUser.Role != "customer" {
			return response.SendStatusForbiddenResponse(c, "Tidak diizinkan: Anda tidak memiliki izin")
		}

		chatRequest := new(dto.CreateChatRequest)
		if err := c.Bind(chatRequest); err != nil {
			return response.SendBadRequestResponse(c, "Format input yang Anda masukkan tidak sesuai.")
		}

		newUser := &entities.ChatModel{
			Text: chatRequest.Text,
		}

		chat, err := h.service.CreateAnswer(currentUser.ID, *newUser)
		if err != nil {
			return response.SendStatusInternalServerResponse(c, "Gagal membuat chat: "+err.Error())
		}

		return response.SendSuccessResponse(c, "Berhasil mendapatkan jawaban", chat)
	}
}

func (h *AssistantHandler) GetChatByIdUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		currentUser := c.Get("CurrentUser").(*entities.UserModels)
		if currentUser.Role != "customer" {
			return response.SendStatusForbiddenResponse(c, "Tidak diizinkan: Anda tidak memiliki izin")
		}

		chat, err := h.service.GetChatByIdUser(currentUser.ID)
		if err != nil {
			return response.SendStatusInternalServerResponse(c, "Gagal mendapatkan chat by id: "+err.Error())
		}

		return response.SendSuccessResponse(c, "Berhasil mendapatkan chat by id User", chat)
	}
}

func (h *AssistantHandler) GenerateArticle() echo.HandlerFunc {
	return func(c echo.Context) error {
		request := new(dto.GenerateArticleAiRequest)
		if err := c.Bind(request); err != nil {
			return response.SendBadRequestResponse(c, "Format input yang Anda masukkan tidak sesuai.")
		}

		chat, err := h.service.GenerateArticle(request.Text)
		if err != nil {
			return response.SendStatusInternalServerResponse(c, "Gagal generate artikel: "+err.Error())
		}

		return response.SendSuccessResponse(c, "Berhasil mendapatkan jawaban", chat)
	}
}

func (h *AssistantHandler) GetProductByIdUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		currentUser := c.Get("CurrentUser").(*entities.UserModels)
		if currentUser.Role != "customer" {
			return response.SendStatusForbiddenResponse(c, "Tidak diizinkan: Anda tidak memiliki izin")
		}

		chat, err := h.service.GenerateRecommendationProduct(currentUser.ID)
		if err != nil {
			return response.SendStatusInternalServerResponse(c, "Gagal rekomendasi "+err.Error())
		}

		return response.SendSuccessResponse(c, "Berhasil rekomendasi", chat)
	}
}
