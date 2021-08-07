package models

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/AnhHoangQuach/go-intern-spores/utils"
)

type Transaction struct {
	Id        string    `gorm:"primary_key;size:255;not null" json:"id"`
	ItemId    string    `gorm:"size:255;not null" json:"item_id" binding:"required"`
	TxHash    string    `gorm:"size:255;not null" json:"tx_hash" binding:"required"`
	Buyer     string    `gorm:"size:255;not null" json:"buyer" binding:"required"`
	Seller    string    `gorm:"size:255;not null" json:"seller" binding:"required"`
	Price     float64   `json:"price" binding:"required"`
	Status    string    `gorm:"size:255;not null;default:Pending" json:"status" binding:"required"`
	Fee       float64   `gorm:"size:255;not null" json:"fee" binding:"required"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
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
	if err := DB.Model(&tx).Where("tx_hash = ?", tx.TxHash).Save(&tx).Error; err != nil {
		return fmt.Errorf("Transaction failed")
	}
	return nil
}

func (t *TxModel) Create(hash string, item_id string, buyer, seller string, price float64, fee float64) (*Transaction, error) {
	var tx = &Transaction{
		Id:     utils.NewGuuid().NewString(),
		ItemId: item_id,
		TxHash: hash,
		Buyer:  buyer,
		Seller: seller,
		Price:  price,
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

func (t *TxModel) TxPagination(tx *Transaction, pagination *Pagination, item_id string) (*[]Transaction, int64, int64, error) {
	var trans []Transaction
	var totalRows int64
	offset := (pagination.Page - 1) * pagination.Limit
	queryBuilder := DB.Limit(pagination.Limit).Offset(offset).Order(pagination.Sort)

	// generate where query
	searchs := pagination.Searchs

	if searchs != nil {
		for _, value := range searchs {
			column := value.Column
			action := value.Action
			query := value.Query

			switch action {
			case "equals":
				whereQuery := fmt.Sprintf("%s = ?", column)
				queryBuilder = queryBuilder.Where(whereQuery, query)
				break
			case "contains":
				whereQuery := fmt.Sprintf("%s LIKE ?", column)
				queryBuilder = queryBuilder.Where(whereQuery, "%"+query+"%")
				break
			case "in":
				whereQuery := fmt.Sprintf("%s IN (?)", column)
				queryArray := strings.Split(query, ",")
				queryBuilder = queryBuilder.Where(whereQuery, queryArray)
				break

			}
		}
	}

	result := queryBuilder.Model(&Transaction{}).Where("item_id = ?", item_id).Find(&trans)

	result.Model(&Transaction{}).Count(&totalRows)
	totalPages := int64(math.Ceil(float64(totalRows) / float64(pagination.Limit)))
	if result.Error != nil {
		msg := result.Error
		return nil, 0, 0, msg
	}
	return &trans, totalRows, totalPages, nil
}
