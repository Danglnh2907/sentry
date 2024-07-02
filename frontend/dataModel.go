package main

type User struct {
	Username       string        `json:"username"`
	Password       string        `json:"password"`
	Fullname       string        `json:"fullname"`
	Budget         float64       `json:"budget"`
	PreferCurrency string        `json:"prefer-currency" default:"USD"`
	Transactions   []Transaction `json:"transaction"`
}

type Transaction struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Descripiton string  `json:"description"`
	Category    string  `json:"category"`
	Date        string  `json:"date"`
	Cost        float64 `json:"cost"`
}
