package quote

// Получить цитату
// добавить в параметры цитату и язык
// перевести цитату

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/go-resty/resty/v2"
)

//Создаем структуры для хранения данных, которые получим с сервера

// Структура для цитаты
type Data struct {
	Text   string `json:"q"`
	Author string `json:"a"`
}

// Структура для перевода цитаты
type Translated struct {
	ResponseData struct {
		TranslatedText string `json:"translatedText"`
	} `json:"responseData"`
}

func getQuote(client *resty.Client) ([]Data, error) {
	// Переменная для полученных данных
	var data []Data

	// Get запрос для получения цитаты
	resp, err := client.R().
		Get("https://zenquotes.io/api/random")

	if err != nil {
		return nil, errors.New("не получили цитату")
	}

	// Декодируем полученный ответ
	err = json.Unmarshal(resp.Body(), &data)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func translatedQuote(client *resty.Client, data []Data) (string, error) {
	// Переменная для полученного перевода
	var translated Translated

	// Get запрос на перевод
	resp, err := client.R().SetQueryParam("q", data[0].Text).SetQueryParam("langpair", "en|ru").
		Get("https://api.mymemory.translated.net/get")
	if err != nil {
		return "", errors.New("не смогли перевести")
	}

	// Декодируем ответ
	err = json.Unmarshal(resp.Body(), &translated)
	if err != nil {
		return "", err
	}

	return translated.ResponseData.TranslatedText, nil
}

func Quote() string {
	// Создаём новый клиент с таймаутом 5 секунд
	client := resty.New()
	client.SetTimeout(5 * time.Second)

	//Ходим за цитатой и переводом
	data, err := getQuote(client)
	if err != nil {
		return "Цитаты не будет. Знай, что у тебя все получится"
	}
	quoteRus, err := translatedQuote(client, data)

	if data[0].Author == "Unknown" {
		data[0].Author = "Неизвестно кто"
	}

	if err != nil {
		log.Println("Ошибка при получении цитаты:", err)
		return "Цитаты не будет. Знай, что у тебя все получится"
	}

	return fmt.Sprintf("%s\n ✍️ %s", quoteRus, data[0].Author)
}
