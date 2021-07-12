package models

import (
	"fmt"
	"time"
)

type Item struct {
	ID          uint32    `gorm:"primary_key;auto_increment" json:"id"`
	Name        string    `gorm:"size:255;not null" json:"name"`
	Description string    `gorm:"size:255;not null" json:"description"`
	Price       int64     `json:"price"`
	Currency    string    `gorm:"size:255;not null" json:"currency"`
	OwnerID     uint32    `gorm:"not null" json:"owner_id"`
	Owner       string    `gorm:"size:255;not null" json:"owner"`
	Creator     string    `gorm:"size:255;not null" json:"creator"`
	Metadata    string    `gorm:"size:255;not null" json:"metadata"`
	Status      string    `gorm:"size:255;not null;default:Pending" json:"status"`
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

type ItemModel struct{}

func (i *ItemModel) Save(item *Item) error {
	if err := DB.Create(&item).Error; err != nil {
		return fmt.Errorf("Save item failed")
	}
	return nil
}

func (i *ItemModel) Update(item *Item) error {
	if err := DB.Model(&item).Where("id = ?", item.ID).Save(&item).Error; err != nil {
		return fmt.Errorf("Save item failed")
	}
	return nil
}

func (i *ItemModel) Create(name, description, currency, owner, creator string, price int64, owner_id uint32) (*Item, error) {
	var item = &Item{
		Name:        name,
		Description: description,
		Price:       price,
		Currency:    currency,
		Owner:       owner,
		OwnerID:     uint32(owner_id),
		Creator:     creator,
	}

	err := i.Save(item)
	if err != nil {
		return nil, fmt.Errorf("Create item failed")
	}

	return item, nil
}

func (i *ItemModel) FindByID(id uint32) (*Item, error) {
	var result Item
	if err := DB.Where("id = ?", id).First(&result).Error; err != nil {
		return nil, err
	}
	return &result, nil
}

func (i *ItemModel) Delete(id uint32) error {
	var result Item
	if err := DB.Where("id = ?", id).Delete(&result).Error; err != nil {
		return err
	}
	return nil
}

func (i *ItemModel) AddMetadataLink(id uint32, metadata string) error {
	item, err := i.FindByID(id)
	if err != nil {
		return err
	}
	item.Metadata = metadata
	err = i.Update(item)
	if err != nil {
		return err
	}
	return nil
}
