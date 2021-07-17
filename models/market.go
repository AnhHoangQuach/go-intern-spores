package models

type MarketModel struct{}

func (m *MarketModel) CalculateRevenue(day, month, year, queryType int, started, to string) float64 {
	var price float64
	if queryType == 1 {
		DB.Model(&Transaction{}).Raw("SELECT SUM(price) FROM transactions WHERE date_part('day', created_at) = ? AND date_part('month', created_at) = ? AND date_part('year', created_at) = ?", day, month, year).Scan(&price)
	}
	if queryType == 2 {
		DB.Model(&Transaction{}).Raw("SELECT SUM(price) FROM transactions WHERE date(created_at) BETWEEN ? AND ?", started, to).Scan(&price)
	}
	return price
}

func (m *MarketModel) CountUserInDay() int64 {
	var sum int64
	DB.Model(&Transaction{}).Raw("SELECT COUNT(id) FROM users WHERE date(created_at) = CURRENT_DATE").Scan(&sum)
	return sum
}

func (m *MarketModel) ListItemsNew() []*Item {
	var items []*Item
	DB.Model(&Transaction{}).Raw("SELECT * FROM items ORDER BY created_at DESC LIMIT 5").Scan(&items)
	return items
}

func (m *MarketModel) ListAuctionsNew() []*Auction {
	var auctions []*Auction
	DB.Model(&Transaction{}).Raw("SELECT * FROM auctions ORDER BY created_at DESC LIMIT 5").Scan(&auctions)
	return auctions
}

func (m *MarketModel) SellestItems() []*Item {
	var items []*Item

	DB.Model(&Transaction{}).Raw("").Scan(&items)
	return items
}

func (m *MarketModel) Bighestitems() []*Item {
	var items []*Item

	DB.Model(&Transaction{}).Raw("").Scan(&items)
	return items
}
