package api

import (
	"kv-cache/cache"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	cache *cache.Cache
}

func NewHandler(cache *cache.Cache) *Handler {
	return &Handler{cache: cache}
}

type PutRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Key     string `json:"key,omitempty"`
	Value   string `json:"value,omitempty"`
}

func (h *Handler) Put(c echo.Context) error {
	var req PutRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Status:  "ERROR",
			Message: "Invalid request format",
		})
	}

	existed, err := h.cache.Put(req.Key, req.Value)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Status:  "ERROR",
			Message: err.Error(),
		})
	}

	message := "Key inserted successfully"
	if existed {
		message = "Key updated successfully"
	}

	return c.JSON(http.StatusOK, Response{
		Status:  "OK",
		Message: message,
	})
}

func (h *Handler) Get(c echo.Context) error {
	key := c.QueryParam("key")
	if key == "" {
		return c.JSON(http.StatusBadRequest, Response{
			Status:  "ERROR",
			Message: "Key parameter is required",
		})
	}

	value, exists := h.cache.Get(key)
	if !exists {
		return c.JSON(http.StatusNotFound, Response{
			Status:  "ERROR",
			Message: "Key not found.",
		})
	}

	return c.JSON(http.StatusOK, Response{
		Status: "OK",
		Key:    key,
		Value:  value,
	})
} 