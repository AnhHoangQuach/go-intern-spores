package models

import (
	"fmt"
	"time"

	"github.com/AnhHoangQuach/go-intern-spores/utils"
)

type Auction struct {
	Id           string    `json:"id" gorm:"primary_key;size:255;not null" binding:"required"`
	ItemId       string    `gorm:"size:255;not null" json:"item_id" binding:"required"`
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
	if err := DB.Model(&auction).Where("id = ?", auction.Id).Save(&auction).Error; err != nil {
		return fmt.Errorf("Save auction failed")
	}
	return nil
}

func (a *AuctionModel) Create(item_id string, initial_price float64, end_at int) (*Auction, error) {
	end_at_time := time.Now().AddDate(0, 0, end_at)

	var auction = &Auction{
		Id:           utils.NewGuuid().NewString(),
		ItemId:       item_id,
		InitialPrice: initial_price,
		FinalPrice:   initial_price,
		EndAt:        end_at_time,
	}

	err := a.Save(auction)
	if err != nil {
		return &Auction{}, fmt.Errorf("Create auction failed")
	}

	return auction, nil
}

func (a *AuctionModel) FindByID(id string) (*Auction, error) {
	var result Auction
	if err := DB.Where("id = ?", id).First(&result).Error; err != nil {
		return &Auction{}, err
	}
	return &result, nil
}

func (a *AuctionModel) FindByItemId(id string) (*Auction, error) {
	var result Auction
	if err := DB.Where("item_id = ?", id).First(&result).Error; err != nil {
		return nil, err
	}
	return &result, nil
}

func (a *AuctionModel) Delete(id string) error {
	var result Auction
	if err := DB.Where("id = ?", id).Delete(&result).Error; err != nil {
		return err
	}
	return nil
}

func (a *AuctionModel) Bid(id string, amount float64) (*Auction, error) {
	auction, err := a.FindByID(id)
	if err != nil {
		return &Auction{}, err
	}
	if amount <= auction.FinalPrice {
		return &Auction{}, fmt.Errorf("Please bid bigger than now price")
	}
	auction.FinalPrice = amount
	err = a.Update(auction)
	if err != nil {
		return &Auction{}, fmt.Errorf("Bid is error")
	}
	return auction, nil
}
