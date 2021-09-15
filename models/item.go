package models

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/AnhHoangQuach/go-intern-spores/utils"
)

type Item struct {
	Id          string    `gorm:"primary_key;size:255;not null"     json:"id"`
	Name        string    `gorm:"size:255;not null"                 json:"name"`
	Description string    `gorm:"size:255;not null"                 json:"description"`
	Price       float64   `                                         json:"price"`
	Currency    string    `gorm:"size:255;not null"                 json:"currency"`
	Owner       string    `gorm:"size:255;not null"                 json:"owner"`
	Creator     string    `gorm:"size:255;not null"                 json:"creator"`
	Metadata    string    `gorm:"size:255;not null"                 json:"metadata"`
	Status      string    `gorm:"size:255;not null;default:Private" json:"status"`
	Type        string    `gorm:"size:255;not null"                 json:"type"`
	Image       string    `                                         json:"image"`
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP"         json:"created_at"`
	UpdatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP"         json:"updated_at"`
	// OwnerID     uint32    `gorm:"not null" json:"owner_id"`
}

type Pagination struct {
	Limit      int      `json:"limit"`
	Page       int      `json:"page"`
	Sort       string   `json:"sort"`
	TotalRows  int64    `json:"total_rows"`
	TotalPages int64    `json:"total_pages"`
	Searchs    []Search `json:"searchs"`
}

type Search struct {
	Column string `json:"column"`
	Action string `json:"action"`
	Query  string `json:"query"`
}

type ItemModel struct{}

func (i *ItemModel) Save(item *Item) error {
	if err := DB.Create(&item).Error; err != nil {
		return fmt.Errorf("Save item failed")
	}
	return nil
}

func (i *ItemModel) Update(item *Item) error {
	if err := DB.Model(&item).Where("id = ?", item.Id).Save(&item).Error; err != nil {
		return fmt.Errorf("Save item failed")
	}
	return nil
}

func (i *ItemModel) Create(
	name, description, currency, owner, creator, type_item, image string,
	price float64,
) (*Item, error) {
	var item = &Item{
		Id:          utils.NewGuuid().NewString(),
		Name:        name,
		Description: description,
		Price:       price,
		Currency:    currency,
		Owner:       owner,
		Creator:     creator,
		Type:        type_item,
		Image:       image,
	}

	err := i.Save(item)
	if err != nil {
		return nil, fmt.Errorf("Create item failed")
	}

	return item, nil
}

func (i *ItemModel) FindByID(id string) (*Item, error) {
	var result Item
	if err := DB.Where("id = ?", id).First(&result).Error; err != nil {
		return nil, err
	}
	return &result, nil
}

func (i *ItemModel) Delete(id string) error {
	var result Item
	if err := DB.Where("id = ?", id).Delete(&result).Error; err != nil {
		return err
	}
	return nil
}

func (i *ItemModel) AddMetadataLink(id string, metadata string) error {
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

func (i *ItemModel) Pagination(
	item *Item,
	pagination *Pagination,
	params ...string,
) (*[]Item, int64, int64, error) {
	var items []Item
	var totalRows int64
	offset := (pagination.Page - 1) * pagination.Limit

	queryBuilder := DB.Model(&Item{})

	if len(params) == 2 {
		queryBuilder = DB.Model(&Item{}).Where(&Item{Status: params[0], Owner: params[1]})
	} else if len(params) == 1 {
		queryBuilder = DB.Model(&Item{}).Where(&Item{Status: params[0]})
	}

	queryBuilder.Count(&totalRows)

	queryBuilder = queryBuilder.Limit(pagination.Limit).Offset(offset).Order(pagination.Sort)

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

	queryBuilder.Model(&Item{}).Find(&items)

	totalPages := int64(math.Ceil(float64(totalRows) / float64(pagination.Limit)))
	return &items, totalRows, totalPages, nil
}
