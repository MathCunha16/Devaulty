package handler

import (
	"errors"
	"log"
	"net/http"

	"devaulty-backend/internal/dto"
	"devaulty-backend/internal/usecase"

	"github.com/gin-gonic/gin"
)

type SecurityHandler struct {
	vaultUseCase *usecase.VaultUseCase
}

func NewSecurityHandler(vaultUseCase *usecase.VaultUseCase) *SecurityHandler {
	return &SecurityHandler{vaultUseCase: vaultUseCase}
}

func (h *SecurityHandler) SetupMasterPassword(c *gin.Context) {
	var masterPassword dto.MasterPassword
	err := c.ShouldBindJSON(&masterPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	password := []byte(masterPassword.MasterPassword)
	masterPassword.MasterPassword = "" // this will only remove the reference to the string, GC will handle eventually
	defer clear(password)

	err = h.vaultUseCase.SetupMasterPassword(c.Request.Context(), password)
	if err != nil {
		if errors.Is(err, usecase.ErrMasterPasswordAlreadyConfigured) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		log.Printf("[SecurityHandler.SetupMasterPassword] %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *SecurityHandler) CheckMasterPasswordSetup(g *gin.Context) {
	isRequired, err := h.vaultUseCase.CheckIfMasterSetupIsRequired(g.Request.Context())
	if err != nil {
		log.Printf("[SecurityHandler.CheckMasterPasswordSetup] %v", err)
		g.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	g.JSON(http.StatusOK, gin.H{"isRequired": isRequired})
}

func (h *SecurityHandler) GetSessionStatus(c *gin.Context) {
	vaultStatus := h.vaultUseCase.GetSessionStatus()
	c.JSON(http.StatusOK, vaultStatus)
}

func (h *SecurityHandler) UnlockVault(c *gin.Context) {
	var masterPassword dto.MasterPassword
	err := c.ShouldBindJSON(&masterPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	password := []byte(masterPassword.MasterPassword)
	masterPassword.MasterPassword = "" // this will only remove the reference to the string, GC will handle eventually
	defer clear(password)

	result, err := h.vaultUseCase.UnlockVault(c.Request.Context(), password)
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidMasterPassword) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, usecase.ErrMasterPasswordNotConfigured) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		log.Printf("[SecurityHandler.UnlockVault] %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *SecurityHandler) LockVault(c *gin.Context) {
	h.vaultUseCase.LockVault()
	c.Status(http.StatusNoContent)
}
