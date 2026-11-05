package main

import (
	"log"

	"github.com/google/uuid"
)

func SeedData() {
	// ID tetap agar sesuai dengan frontend
	clientID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	adminID := uuid.MustParse("22222222-2222-2222-2222-222222222222")
	officerID := uuid.MustParse("33333333-3333-3333-3333-333333333333")

	// 1. SEED USERS

	users := []User{
		{
			ID:    clientID,
			Name:  "Pelanggan Demo",
			Email: "client@test.com",
			Role:  "CLIENT",
		},
		{
			ID:    adminID,
			Name:  "Admin Demo",
			Email: "admin@test.com",
			Role:  "ADMIN",
		},
		{
			ID:    officerID,
			Name:  "Teknisi Demo",
			Email: "officer@test.com",
			Role:  "OFFICER",
		},
	}

	for _, user := range users {
		if err := DB.Where("id = ?", user.ID).
			FirstOrCreate(&user).Error; err != nil {
			log.Fatal("Gagal membuat user:", err)
		}
	}

	// 2. SEED ORDERS

	orders := []Order{
		{
			ID:          uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"),
			ClientID:    clientID,
			ServiceName: "Perbaikan AC",
			Description: "AC tidak dingin",
			Status:      "PENDING",
			Version:     1,
			LastSequence: 1,
		},
		{
			ID:          uuid.MustParse("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"),
			ClientID:    clientID,
			OfficerID:   &officerID,
			ServiceName: "Perbaikan Mesin Cuci",
			Description: "Mesin cuci tidak berputar",
			Status:      "ASSIGNED",
			Version:     2,
			LastSequence: 2,
		},
		{
			ID:          uuid.MustParse("cccccccc-cccc-cccc-cccc-cccccccccccc"),
			ClientID:    clientID,
			OfficerID:   &officerID,
			ServiceName: "Perbaikan Kulkas",
			Description: "Kulkas tidak dingin",
			Status:      "IN_PROGRESS",
			Version:     4,
			LastSequence: 4,
		},
	}

	for _, order := range orders {
		if err := DB.Where("id = ?", order.ID).
			FirstOrCreate(&order).Error; err != nil {
			log.Fatal("Gagal membuat order:", err)
		}
	}

	// 3. SEED ORDER EVENTS

	events := []OrderEvent{
		// Order PENDING
		{
			OrderID: uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"),
			Sequence: 1,
			NewStatus: "PENDING",
			ActorID: clientID,
			Note: "Pesanan berhasil dibuat",
		},

		// Order ASSIGNED
		{
			OrderID: uuid.MustParse("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"),
			Sequence: 1,
			NewStatus: "PENDING",
			ActorID: clientID,
			Note: "Pesanan berhasil dibuat",
		},
		{
			OrderID: uuid.MustParse("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"),
			Sequence: 2,
			OldStatus: "PENDING",
			NewStatus: "ASSIGNED",
			ActorID: adminID,
			Note: "Teknisi berhasil ditugaskan",
		},

		// Order IN_PROGRESS
		{
			OrderID: uuid.MustParse("cccccccc-cccc-cccc-cccc-cccccccccccc"),
			Sequence: 1,
			NewStatus: "PENDING",
			ActorID: clientID,
			Note: "Pesanan berhasil dibuat",
		},
		{
			OrderID: uuid.MustParse("cccccccc-cccc-cccc-cccc-cccccccccccc"),
			Sequence: 2,
			OldStatus: "PENDING",
			NewStatus: "ASSIGNED",
			ActorID: adminID,
			Note: "Teknisi berhasil ditugaskan",
		},
		{
			OrderID: uuid.MustParse("cccccccc-cccc-cccc-cccc-cccccccccccc"),
			Sequence: 3,
			OldStatus: "ASSIGNED",
			NewStatus: "ON_THE_WAY",
			ActorID: officerID,
			Note: "Teknisi sedang menuju lokasi",
		},
		{
			OrderID: uuid.MustParse("cccccccc-cccc-cccc-cccc-cccccccccccc"),
			Sequence: 4,
			OldStatus: "ON_THE_WAY",
			NewStatus: "IN_PROGRESS",
			ActorID: officerID,
			Note: "Pengerjaan sedang berlangsung",
		},
	}

	for _, event := range events {
		if err := DB.Where(
			"order_id = ? AND sequence = ?",
			event.OrderID,
			event.Sequence,
		).FirstOrCreate(&event).Error; err != nil {
			log.Fatal("Gagal membuat event:", err)
		}
	}

	log.Println("Seed data berhasil dibuat!")
}