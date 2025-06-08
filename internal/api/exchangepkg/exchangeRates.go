package exchangepkg

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/go-resty/resty/v2"
)

// Структура для получения курсов
type Money struct {
	Rub struct {
		Usd float64 `json:"usd"`
		Eur float64 `json:"eur"`
		Aed float64 `json:"aed"`
	} `json:"rub"`
}

// Создаем структуры для хранения данных, которые получим с сервера
func getExchange(client *resty.Client) (string, error) {

	// Переменная для полученного курса
	var money Money

	resp, err := client.R().
		Get("https://cdn.jsdelivr.net/npm/@fawazahmed0/currency-api@latest/v1/currencies/rub.json")
	if err != nil {
		return "", fmt.Errorf("не получили курс от API: %w", err)
	}

	// Декодируем ответ
	err = json.Unmarshal(resp.Body(), &money)
	if err != nil {
		return "", fmt.Errorf("не смогли декодировать ответ: %w", err)
	}

	year, month, day := time.Now().Date()

	return fmt.Sprintf("Курс валют на %v.%v.%v:\nДоллар: %.2f₽\nЕвро: %.2f₽\nДирхам: %.2f₽\n\nКурс обновится завтра", day, int(month), year,
		1.0/money.Rub.Usd, 1.0/money.Rub.Eur, 1.0/money.Rub.Aed), nil

}

func ExchangeRates() string {
	// Создаём новый клиент с таймаутом 5 секунд
	client := resty.New()
	client.SetTimeout(5 * time.Second)

	//Ходим за курсом
	result, err := getExchange(client)

	// показываем итог
	if err != nil {
		log.Println("Ошибка при получении курса:", err)
		return "Не получилость узнать курс :("
	} else {
		return result
	}

}
