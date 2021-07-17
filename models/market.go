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
