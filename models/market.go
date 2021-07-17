package models

import (
	"fmt"
	"strconv"
	"time"
)

type MarketModel struct{}

func HandleDate(input int) string {
	out := strconv.Itoa(input)
	zero := fmt.Sprintf("%c", '0')
	if input < 10 {
		return zero + out
	}
	return out
}

// Bug: tai sao 23:00 lai tinh la ngay moi, tuong tu voi month

func (m *MarketModel) CalculateRevenue(cal_type string, time_query int) []*Transaction {
	var txs []*Transaction

	yearNow := HandleDate(time.Now().Year())
	monthNow := HandleDate(int(time.Now().Month()))
	dateNow := HandleDate(time.Now().Day())
	tempChar := fmt.Sprintf("%c", '-')

	if cal_type == "date" {
		queryDay := yearNow + tempChar + monthNow + tempChar + HandleDate(time_query)
		fmt.Println(queryDay)
		DB.Model(&Transaction{}).Raw("SELECT * FROM transactions WHERE date_part('day', created_at) = date_part('day', ?::TIMESTAMP')").Find(&txs)
	} else if cal_type == "month" {
		queryMonth := yearNow + tempChar + HandleDate(time_query) + tempChar + dateNow
		DB.Model(&Transaction{}).Raw("SELECT * FROM transactions WHERE date_part('month', created_at) = date_part('month', ?::TIMESTAMP)", queryMonth).Find(&txs)
	} else if cal_type == "year" {
		queryYear := HandleDate(time_query) + tempChar + monthNow + tempChar + dateNow
		DB.Model(&Transaction{}).Raw("SELECT * FROM transactions WHERE date_part('year', created_at) = date_part('year', ?::TIMESTAMP)", queryYear).Find(&txs)
	}

	return txs
}
