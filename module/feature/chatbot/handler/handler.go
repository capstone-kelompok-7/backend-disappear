package handler

import (
	"github.com/capstone-kelompok-7/backend-disappear/module/entities"
	"github.com/capstone-kelompok-7/backend-disappear/module/feature/chatbot"
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
		chatRequest := new(entities.ChatModel)
		if err := c.Bind(chatRequest); err != nil {
			return response.SendBadRequestResponse(c, "Format input yang Anda masukkan tidak sesuai.")
		}

		err := h.service.CreateQuestion(*chatRequest)
		if err != nil {
			return response.SendStatusInternalServerResponse(c, "Gagal membuat chat: "+err.Error())
		}

		return response.SendSuccessResponse(c, "Berhasil membuat chat", chatRequest.Text)
	}
}
func (h *ChatbotHandler) CreateAnswer() echo.HandlerFunc {
	return func(c echo.Context) error {
		chatRequest := new(entities.ChatModel)
		if err := c.Bind(chatRequest); err != nil {
			return response.SendBadRequestResponse(c, "Format input yang Anda masukkan tidak sesuai.")
		}

		chat, err := h.service.CreateAnswer(*chatRequest)
		if err != nil {
			return response.SendStatusInternalServerResponse(c, "Gagal membuat chat: "+err.Error())
		}

		return response.SendSuccessResponse(c, "Berhasil mendapatkan jawaban", chat)
	}
}
func (h *ChatbotHandler) GetChatByIdUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		chat, err := h.service.GetChatByIdUser(id)
		if err != nil {
			return response.SendStatusInternalServerResponse(c, "Gagal mendapatkan chat by id: "+err.Error())
		}

		return response.SendSuccessResponse(c, "Berhasil mendapatkan chat by id User", chat)
	}
}
