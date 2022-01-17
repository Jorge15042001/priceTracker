package main

import (
	_ "fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Product struct {
	ProductID        uint            `gorm:"primaryKey;autoIncrement" json:"ProductID"`
	ProductName      string          `gorm:"not null" json:"ProductName"`
	TrackingInterval uint            `gorm:"not null" json:"TrackingInterval"`
	Options          []ProductOption `gorm:"foreignKey:ProductID;references:ProductID;constraint:OnUpdate:CASCADE,OnSave:CASCADE,OnDelete:CASCADE"`
}

type ProductOption struct {
	ProductID   uint         `json:"ProductID"`
	OptionID    uint         `gorm:"primaryKey;autoIncrement" json:"OptionID"`
	Opinion     float32      `gorm:"not null" json:"Opinion"`
	ArrivalDate time.Time    `gorm:"not null" json:"ArrivalDate"`
	ImageSrc    string       `gorm:"not null" json:"ImageSrc"`
	Url         string       `gorm:"not null" json:"Url"`
	Prices      []PricePoint `gorm:"ForeignKey:OptionID;references:OptionID;constraint:OnUpdate:CASCADE,OnSave:CASCADE,OnDelete:CASCADE"`
}

type PricePoint struct {
	OptionID     uint      `json:"OptionID"`
	PricePointID uint      `gorm:"primaryKey;autoIncrement" json:"PricePointID"`
	CheckDate    time.Time `gorm:"not null" json:"CheckDate"`
	Price        float32   `gorm:"not null" json:"Price"`
}

func getDataBaseConnection() *gorm.DB {
	dsn := "root:123456@tcp(127.0.0.1:3306)/pricetracker?parseTime=true"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&Product{})
	db.AutoMigrate(&ProductOption{})
	db.AutoMigrate(&PricePoint{})

	// 	product := Product{ProductName: "monitor"}
	// 	db.Create(&product)

	// 	option := ProductOption{Opinion: 1.0, ArrivalDate: time.Now(), imageSrc: "hola"}
	// 	option2 := ProductOption{Opinion: 2.0, ArrivalDate: time.Now(), imageSrc: "hola"}

	// 	product.Options = append(product.Options, option)
	// 	product.Options = append(product.Options, option2)

	// 	price := PricePoint{CheckDate: time.Now(), Price: 100.00}
	// 	price2 := PricePoint{CheckDate: time.Now(), Price: 101.00}
	// 	price3 := PricePoint{CheckDate: time.Now(), Price: 102.00}

	// 	product.Options[0].Prices = append(product.Options[0].Prices, price)
	// 	product.Options[0].Prices = append(product.Options[0].Prices, price2)
	// 	product.Options[0].Prices = append(product.Options[0].Prices, price3)
	// 	fmt.Print(product.Options[0].Prices)

	// 	fmt.Print("\n\n")
	// 	db.Save(product)
	// 	fmt.Print(product.Options[0].Prices)

	return db
}
