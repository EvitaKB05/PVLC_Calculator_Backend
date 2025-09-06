package repository

import (
	"fmt"
)

type Repository struct {
}

func NewRepository() (*Repository, error) {
	return &Repository{}, nil
}

type Order struct { // вот наша новая структура
	ID    int    // поля структур, которые передаются в шаблон
	Title string // ОБЯЗАТЕЛЬНО должны быть написаны с заглавной буквы (то есть публичными)
}

func (r *Repository) GetOrders() ([]Order, error) {
	// имитируем работу с БД. Типа мы выполнили sql запрос и получили эти строки из БД
	orders := []Order{ // массив элементов из наших структур
		{
			ID:    1,
			Title: "first order",
		},
		{
			ID:    2,
			Title: "second order",
		},
		{
			ID:    3,
			Title: "third order",
		},
	}
	// обязательно проверяем ошибки, и если они появились - передаем выше, то есть хендлеру
	// тут я снова искусственно обработаю "ошибку" чисто чтобы показать вам как их передавать выше
	if len(orders) == 0 {
		return nil, fmt.Errorf("массив пустой")
	}

	return orders, nil
}
