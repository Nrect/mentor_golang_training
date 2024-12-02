package main

import (
	"flag"
	"fmt"
	"mentor-training/pkg/logger"
	"mentor-training/pkg/queue"
	"mentor-training/queue_broker_service/handlers"
	"net/http"
	_ "net/http/pprof"
)

func main() {
	log := logger.New(true)

	port := flag.Int("port", 8080, "Порт, на котором будет запущен сервер")
	pprofPort := flag.Int("pprof-port", 6060, "Порт для pprof")
	flag.Parse()

	qm := queue.NewQueueManager()
	handler := handlers.NewHandler(log, qm)

	go func() {
		address := fmt.Sprintf(":%d", *port)
		log.Infof("Сервер запущен по адресу http://127.0.0.1%s", address)
		log.Info("Доступные эндпоинты:")
		log.Info("  PUT  /{name}?v=... - Добавить сообщение в очередь")
		log.Info("  GET  /{name}?timeout=N - Получить сообщение из очереди")

		if err := http.ListenAndServe(address, handler); err != nil {
			log.Fatalf("Ошибка запуска сервера: %v", err)
		}
	}()

	go func() {
		address := fmt.Sprintf(":%d", *pprofPort)
		log.Infof("Pprof доступен по адресу http://127.0.0.1%s/debug/pprof/", address)
		if err := http.ListenAndServe(address, nil); err != nil {
			log.Fatalf("Ошибка запуска pprof: %v", err)
		}
	}()

	select {}
}
