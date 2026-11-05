package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CreateOrderRequest struct {
	ClientID    string `json:"client_id" binding:"required"`
	ServiceName string `json:"service_name" binding:"required"`
	Description string `json:"description"`
}

func CreateOrder(c *gin.Context) {
	var req CreateOrderRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Data pesanan tidak valid",
		})
		return
	}

	clientID, err := uuid.Parse(req.ClientID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Client ID tidak valid",
		})
		return
	}

	var client User

	if err := DB.First(&client, "id = ?", clientID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Pelanggan tidak ditemukan",
		})
		return
	}

	if client.Role != "CLIENT" {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Hanya pelanggan yang dapat membuat pesanan",
		})
		return
	}

	order := Order{
		ClientID:     clientID,
		ServiceName:  req.ServiceName,
		Description:  req.Description,
		Status:       "PENDING",
		Version:      1,
		LastSequence: 1,
	}

	err = DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&order).Error; err != nil {
			return err
		}

		event := OrderEvent{
			OrderID:   order.ID,
			Sequence:  1,
			NewStatus: "PENDING",
			ActorID:   clientID,
			Note:      "Pesanan berhasil dibuat",
		}

		return tx.Create(&event).Error
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal membuat pesanan",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Pesanan berhasil dibuat",
		"data":    order,
	})
}

func GetOrders(c *gin.Context) {
	var orders []Order

	if err := DB.Order("created_at DESC").Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal mengambil data pesanan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data pesanan berhasil diambil",
		"data":    orders,
	})
}

func GetOrderByID(c *gin.Context) {
	id := c.Param("id")

	orderID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID pesanan tidak valid",
		})
		return
	}

	var order Order

	if err := DB.First(&order, "id = ?", orderID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Pesanan tidak ditemukan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Detail pesanan berhasil diambil",
		"data":    order,
	})
}

func GetOrderEvents(c *gin.Context) {
	id := c.Param("id")

	orderID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID pesanan tidak valid",
		})
		return
	}

	var events []OrderEvent

	if err := DB.Where("order_id = ?", orderID).
		Order("sequence ASC").
		Find(&events).Error; err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal mengambil riwayat pesanan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Riwayat pesanan berhasil diambil",
		"data":    events,
	})
}

// UPDATE STATUS

type UpdateStatusRequest struct {
	NewStatus       string `json:"new_status" binding:"required"`
	ActorID         string `json:"actor_id" binding:"required"`
	ExpectedVersion int    `json:"expected_version" binding:"required"`
}

func UpdateOrderStatus(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID pesanan tidak valid",
		})
		return
	}

	var req UpdateStatusRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Data tidak valid",
		})
		return
	}

	actorID, err := uuid.Parse(req.ActorID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Actor ID tidak valid",
		})
		return
	}

	allowedStatuses := map[string]bool{
		"ASSIGNED":    true,
		"ON_THE_WAY":  true,
		"IN_PROGRESS": true,
		"COMPLETED":   true,
	}

	if !allowedStatuses[req.NewStatus] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Status tidak valid",
		})
		return
	}

	var order Order

	err = DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&order, "id = ?", id).Error; err != nil {
			return err
		}

		validTransitions := map[string]string{
			"PENDING":     "ASSIGNED",
			"ASSIGNED":    "ON_THE_WAY",
			"ON_THE_WAY":  "IN_PROGRESS",
			"IN_PROGRESS": "COMPLETED",
		}

		if validTransitions[order.Status] != req.NewStatus {
			return fmt.Errorf("perubahan status tidak diizinkan")
		}

		result := tx.Model(&Order{}).
			Where("id = ? AND version = ?", id, req.ExpectedVersion).
			Updates(map[string]interface{}{
				"status":        req.NewStatus,
				"version":       order.Version + 1,
				"last_sequence": order.LastSequence + 1,
			})

		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected == 0 {
			return fmt.Errorf("versi pesanan sudah berubah, silakan muat ulang")
		}

		event := OrderEvent{
			OrderID:   id,
			Sequence:  order.LastSequence + 1,
			OldStatus: order.Status,
			NewStatus: req.NewStatus,
			ActorID:   actorID,
			Note:      "Status pesanan diperbarui",
		}

		return tx.Create(&event).Error
	})

	if err != nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Broadcast event melalui WebSocket.
	var event OrderEvent

	if err := DB.Where("order_id = ?", id).
		Order("sequence DESC").
		First(&event).Error; err == nil {

		BroadcastOrderEvent(id, event)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Status pesanan berhasil diperbarui",
	})
}

// CANCEL ORDER

type CancelOrderRequest struct {
	RequestedBy string `json:"requested_by" binding:"required"`
	Reason      string `json:"reason" binding:"required"`
}

func CancelOrder(c *gin.Context) {
	orderID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID pesanan tidak valid",
		})
		return
	}

	var req CancelOrderRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Data pembatalan tidak valid",
		})
		return
	}

	userID, err := uuid.Parse(req.RequestedBy)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User ID tidak valid",
		})
		return
	}

	statusChanged := false

	err = DB.Transaction(func(tx *gorm.DB) error {
		var order Order

		if err := tx.First(&order, "id = ?", orderID).Error; err != nil {
			return err
		}

		if order.Status == "COMPLETED" || order.Status == "CANCELLED" {
			return fmt.Errorf("pesanan sudah tidak dapat dibatalkan")
		}

		if order.Status == "PENDING" || order.Status == "ASSIGNED" {
			result := tx.Model(&Order{}).
				Where("id = ? AND version = ?", orderID, order.Version).
				Updates(map[string]interface{}{
					"status":        "CANCELLED",
					"version":       order.Version + 1,
					"last_sequence": order.LastSequence + 1,
				})

			if result.Error != nil {
				return result.Error
			}

			if result.RowsAffected == 0 {
				return fmt.Errorf("pesanan sudah berubah, silakan coba lagi")
			}

			event := OrderEvent{
				OrderID:   orderID,
				Sequence:  order.LastSequence + 1,
				OldStatus: order.Status,
				NewStatus: "CANCELLED",
				ActorID:   userID,
				Note:      req.Reason,
			}

			if err := tx.Create(&event).Error; err != nil {
				return err
			}

			statusChanged = true
			return nil
		}

		if order.Status == "ON_THE_WAY" || order.Status == "IN_PROGRESS" {
			cancellation := CancellationRequest{
				OrderID:     orderID,
				RequestedBy: userID,
				Reason:      req.Reason,
				Status:      "PENDING",
			}

			return tx.Create(&cancellation).Error
		}

		return fmt.Errorf("status pesanan tidak valid")
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Broadcast jika status langsung berubah menjadi CANCELLED.
	if statusChanged {
		var event OrderEvent

		if err := DB.Where("order_id = ?", orderID).
			Order("sequence DESC").
			First(&event).Error; err == nil {

			BroadcastOrderEvent(orderID, event)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Permintaan pembatalan berhasil diproses",
	})
}

// REVIEW CANCELLATION

type ReviewCancellationRequest struct {
	ReviewerID string `json:"reviewer_id" binding:"required"`
	Decision   string `json:"decision" binding:"required"`
}

func ReviewCancellation(c *gin.Context) {
	cancellationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID pembatalan tidak valid",
		})
		return
	}

	var req ReviewCancellationRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Data tidak valid",
		})
		return
	}

	reviewerID, err := uuid.Parse(req.ReviewerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Reviewer ID tidak valid",
		})
		return
	}

	if req.Decision != "APPROVED" && req.Decision != "REJECTED" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Keputusan tidak valid",
		})
		return
	}

	var orderID uuid.UUID

	err = DB.Transaction(func(tx *gorm.DB) error {
		var admin User

		if err := tx.First(&admin, "id = ?", reviewerID).Error; err != nil {
			return err
		}

		if admin.Role != "ADMIN" {
			return fmt.Errorf("hanya admin yang dapat meninjau pembatalan")
		}

		var cancellation CancellationRequest

		if err := tx.First(&cancellation, "id = ?", cancellationID).Error; err != nil {
			return err
		}

		if cancellation.Status != "PENDING" {
			return fmt.Errorf("permintaan pembatalan sudah diproses")
		}

		var order Order

		if err := tx.First(&order, "id = ?", cancellation.OrderID).Error; err != nil {
			return err
		}

		orderID = order.ID

		if order.Status == "COMPLETED" || order.Status == "CANCELLED" {
			return fmt.Errorf("status pesanan sudah tidak dapat diubah")
		}

		cancellation.Status = req.Decision
		cancellation.ReviewedBy = &reviewerID

		if err := tx.Save(&cancellation).Error; err != nil {
			return err
		}

		if req.Decision == "APPROVED" {
			result := tx.Model(&Order{}).
				Where("id = ? AND version = ?", order.ID, order.Version).
				Updates(map[string]interface{}{
					"status":        "CANCELLED",
					"version":       order.Version + 1,
					"last_sequence": order.LastSequence + 1,
				})

			if result.Error != nil {
				return result.Error
			}

			if result.RowsAffected == 0 {
				return fmt.Errorf("pesanan sudah berubah")
			}

			event := OrderEvent{
				OrderID:   order.ID,
				Sequence:  order.LastSequence + 1,
				OldStatus: order.Status,
				NewStatus: "CANCELLED",
				ActorID:   reviewerID,
				Note:      "Pembatalan disetujui admin",
			}

			return tx.Create(&event).Error
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Broadcast jika admin menyetujui pembatalan.
	if req.Decision == "APPROVED" {
		var event OrderEvent

		if err := DB.Where("order_id = ?", orderID).
			Order("sequence DESC").
			First(&event).Error; err == nil {

			BroadcastOrderEvent(orderID, event)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Permintaan pembatalan berhasil ditinjau",
		"decision": req.Decision,
	})
}