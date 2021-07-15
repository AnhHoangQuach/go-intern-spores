package models

import (
	"fmt"
	"time"
)

type Transaction struct {
	ID        uint32    `gorm:"primary_key;auto_increment" json:"id"`
	ItemID    uint32    `json:"item_id" binding:"required"`
	Hash      string    `gorm:"size:255;not null" json:"hash" binding:"required"`
	From      string    `gorm:"size:255;not null" json:"from" binding:"required"`
	To        string    `gorm:"size:255;not null" json:"to" binding:"required"`
	Amount    uint64    `json:"amount" binding:"required"`
	Status    string    `gorm:"size:255;not null;default:Pending" json:"status" binding:"required"`
	Fee       float64   `gorm:"size:255;not null" json:"fee" binding:"required"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
}

type TxModel struct{}

func (t *TxModel) Save(tx *Transaction) error {
	var err error
	done := make(chan bool)
	go func(ch chan<- bool) {
		defer close(ch)
		err = DB.Model(&Transaction{}).Create(&tx).Error
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

func (t *TxModel) Update(tx *Transaction) error {
	if err := DB.Model(&tx).Where("hash = ?", tx.Hash).Save(&tx).Error; err != nil {
		return fmt.Errorf("Transaction failed")
	}
	return nil
}

func (t *TxModel) Create(hash string, item_id uint32, from, to string, amount uint64, fee float64) (*Transaction, error) {
	var tx = &Transaction{
		ItemID: item_id,
		Hash:   hash,
		From:   from,
		To:     to,
		Amount: amount,
		Fee:    fee,
	}

	err := t.Save(tx)
	if err != nil {
		return nil, fmt.Errorf("Transaction save failed")
	}

	tx.Status = "Success"
	err = t.Update(tx)
	if err != nil {
		return nil, fmt.Errorf("Transaction update failed")
	}
	return tx, nil
}
