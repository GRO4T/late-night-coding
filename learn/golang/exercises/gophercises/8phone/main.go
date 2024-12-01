package main

import (
	"fmt"
	"log"
	"strings"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	PhoneNumber string
}

func normalize(phoneNumber string) string {
	r := strings.NewReplacer(
		"-", "",
		" ", "",
		"(", "",
		")", "",
	)
	return r.Replace(phoneNumber)
}

func main() {
	db, err := gorm.Open(sqlite.Open("phone.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database: ", err)
	}

	db.AutoMigrate(&User{})

	// Poopulate phone numbers
	phoneNumbers := []string{
		"1234567890",
		"123 456 7891",
		"(123) 456 7892",
		"(123) 456-7893",
		"123-456-7894",
		"123-456-7890",
		"1234567892",
		"(123)456-7892",
	}
	for _, p := range phoneNumbers {
		db.Create(&User{PhoneNumber: p})
	}

	// Normalize phone numbers
	uniquePhoneNumbers := make(map[string]struct{})
	var users []User
	db.Find(&users)

	for _, user := range users {
		normalizedPhoneNumber := normalize(user.PhoneNumber)
		if _, ok := uniquePhoneNumbers[normalizedPhoneNumber]; ok {
			db.Delete(&user)
		} else {
			uniquePhoneNumbers[normalizedPhoneNumber] = struct{}{}
		}
	}

	// Print unique normalized phone numbers
	for k := range uniquePhoneNumbers {
		fmt.Println(k)
	}
}
