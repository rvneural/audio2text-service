package main

import (
	"log"
	"math/rand"
	"os"
)

func getRandomFileName(leng int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	rune_name := make([]rune, leng)
	for i := range rune_name {
		rune_name[i] = letters[rand.Intn(len(letters))]
	}
	return string(rune_name)
}

func getTemFileDir() string {

	var temp_files string = "./temp_files/" + getRandomFileName(15) + "/"

	// Проверяем наличие папки для хранения временных файлов
	if _, err := os.Stat(temp_files); err != nil {
		// Если папка отсутствует
		if os.IsNotExist(err) {
			// Пытаемся создать папку
			err2 := os.MkdirAll(temp_files, 0777)
			if err2 != nil {
				// Если не удается — сохраняем в корень проекта
				log.Println("Error while creating directory: ", err2)
				temp_files = "./"
			}
		}
	}

	return temp_files
}
