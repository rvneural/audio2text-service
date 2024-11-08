package normalization

import "os"

// ADDR - Адрес сервиса по работе с текстом
var ADDR = os.Getenv("TEXT_2_TEXT_ADDR")

// TEXT_2_TEXT_KEY - Ключ для авторизации на сервисе по работе с текстом
var TEXT_2_TEXT_KEY = os.Getenv("TEXT_2_TEXT_KEY")
