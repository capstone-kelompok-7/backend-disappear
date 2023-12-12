package handler

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/chatbot"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/chatbot/dto"
	"github.com/capstone-kelompok-7/backend-disappear/utils/response"
	"github.com/labstack/echo/v4"
)

type ChatbotHandler struct {
	service chatbot.ServicChatbotInterface
}

func NewChatbotHandler(service chatbot.ServicChatbotInterface) chatbot.HandlerChatbotInterface {
	return &ChatbotHandler{
		service: service,
	}
}

func (h *ChatbotHandler) CreateQuestion() echo.HandlerFunc {
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

func (h *ChatbotHandler) CreateAnswer() echo.HandlerFunc {
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

func (h *ChatbotHandler) GetChatByIdUser() echo.HandlerFunc {
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

func (h *ChatbotHandler) GenerateArtikelAi() echo.HandlerFunc {
	return func(c echo.Context) error {
		judulRequest := new(dto.GenerateArticleAiRequest)
		if err := c.Bind(judulRequest); err != nil {
			return response.SendBadRequestResponse(c, "Format input yang Anda masukkan tidak sesuai.")
		}

		chat, err := h.service.GenerateArtikelAi(judulRequest.Text)
		if err != nil {
			return response.SendStatusInternalServerResponse(c, "Gagal generate artikel: "+err.Error())
		}

		return response.SendSuccessResponse(c, "Berhasil mendapatkan jawaban", chat)
	}
}
