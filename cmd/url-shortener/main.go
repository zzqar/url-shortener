package main

import (
	"fmt"
	"url-shortener/internal/config"
)

func main() {
	/*
		Init config: cleanenv
		Что в себе он содержит ? // настройки сервера...
		Для чего ? // для чтения параметров настроек из файла
	*/
	cfg := config.MustLoadConfig()
	fmt.Println(cfg)

	/*
		TODO init logger: slog
		логировать ошибки ?
		Почему тут ?
	*/

	/*
		TODO init storage: sqlite
		что он делает ?
		почему тут ?
	*/

	/*
		TODO init router: chi, 'che render'
		Что делает ? // Маршрутизация (например, маршруты для API)
		Почему тут ? // Старт приложения, тут происходит распределение запросов на контроллеры
	*/

	/*
		TODO run server
		что значит запускаем ?

	*/

}
