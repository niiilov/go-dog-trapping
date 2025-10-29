package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func Encode(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func main() {
	users := map[string]string{
		"bolsh_otdel":   "bolsh_otdel48",
		"borns_otdel":   "borns_otdel48",
		"vaslvs_otdel":  "vaslvs_otdel48",
		"vedens_otdel":  "vedens_otdel48",
		"verbl_otdel":   "verbl_otdel48",
		"gryaz_otdel":   "gryaz_otdel48",
		"ivov_otdel":    "ivov_otdel48",
		"kosir_otdel":   "kosir_otdel48",
		"krut_otdel":    "krut_otdel48",
		"kuzot_otdel":   "kuzot_otdel48",
		"lenin_otdel":   "lenin_otdel48",
		"lubn_otdel":    "lubn_otdel48",
		"novder_otdel":  "novder_otdel48",
		"novdmtr_otdel": "novdmtr_otdel48",
		"padov_otdel":   "padov_otde48",
		"pruj_otdel":    "pruj_otdel48",
		"senc_otdel":    "senc_otdel48",
		"steb_otdel":    "steb_otdel48",
		"sirsk_otdel":   "sirsk_otdel48",
		"telej_otdel":   "telej_otdel48",
		"chastdu_otdel": "chastdu_otdel48",
		"admst_com":     "admst_com48",
		"fgbuz":         "fgbuz48",
		"ryaon_comm":    "ryaon_comm48",
	}

	fmt.Println("-- Хэшированные пароли:")
	for login, password := range users {
		hashedPassword, err := Encode(password)
		if err != nil {
			fmt.Printf("Ошибка при хэшировании пароля для %s: %v\n", login, err)
			continue
		}
		fmt.Printf("Login: %-15s | Hash: %s\n", login, hashedPassword)
	}

	fmt.Println("\n-- SQL INSERT для таблицы users:")
	fmt.Println("INSERT INTO users (login, password_hash) VALUES")
	count := 0
	total := len(users)

	for login, password := range users {
		hashedPassword, err := Encode(password)
		if err != nil {
			fmt.Printf("-- Ошибка для %s: %v\n", login, err)
			continue
		}
		count++
		if count == total {
			fmt.Printf("('%s', '%s');", login, hashedPassword)
		} else {
			fmt.Printf("('%s', '%s'),\n", login, hashedPassword)
		}
	}
}
