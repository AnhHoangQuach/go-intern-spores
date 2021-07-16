package models

import (
	"fmt"
	"time"
)

type Auction struct {
	ID           uint32    `gorm:"primary_key;auto_increment" json:"id"`
	ItemID       uint32    `json:"item_id" binding:"required"`
	InitialPrice float64   `json:"initial_price" binding:"required"`
	FinalPrice   float64   `json:"final_price" binding:"required"`
	Status       string    `gorm:"default:Pending" json:"status" binding:"required"`
	CreatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	EndAt        time.Time `json:"end_at" binding:"required"`
}

type AuctionModel struct{}

func (a *AuctionModel) Save(auction *Auction) error {
	var err error
	done := make(chan bool)
	go func(ch chan<- bool) {
		defer close(ch)
		err = DB.Model(&Auction{}).Create(&auction).Error
		if err != nil {
			ch <- false
			return
		}
		ch <- true
	}(done)
	if OK(done) {
		return nil
	}
	return err
}

func (a *AuctionModel) Update(auction *Auction) error {
	if err := DB.Model(&auction).Where("id = ?", auction.ID).Save(&auction).Error; err != nil {
		return fmt.Errorf("Save auction failed")
	}
	return nil
}

func (a *AuctionModel) Create(item_id uint32, initial_price float64, end_at int) (*Auction, error) {
	end_at_time := time.Now().AddDate(0, 0, end_at)

	var auction = &Auction{
		ItemID:       item_id,
		InitialPrice: initial_price,
		EndAt:        end_at_time,
	}

	err := a.Save(auction)
	if err != nil {
		return nil, fmt.Errorf("Create auction failed")
	}

	return auction, nil
}
