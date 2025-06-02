package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Yandex-Practicum/tracker/internal/daysteps"
	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

func main() {
	const (
		weight = 84.6
		height = 1.87
	)

	// Настройка логгера
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetOutput(os.Stdout)

	// Дневная активность
	dailyActivities := []string{
		"678,0h50m",
		"792,1h14m",
		"1078,1h30m",
		"7830,2h40m",
		",3456",
		"12:40:00, 3456",
		"something is wrong",
	}

	fmt.Println("Активность в течение дня")
	fmt.Println("========================")

	for i, activity := range dailyActivities {
		log.Printf("Обработка активности #%d: %s", i+1, activity)
		info := daysteps.DayActionInfo(activity, weight, height)
		if info == "" {
			log.Printf("Пропущена некорректная активность #%d: %s", i+1, activity)
			continue
		}
		fmt.Println(info)
		fmt.Println("-----------------------")
	}

	// Тренировки
	trainings := []string{
		"3456,Ходьба,3h00m",
		"something is wrong",
		"678,Бег,0h5m",
		"1078,Бег,0h10m",
		",3456 Ходьба",
		"7892,Ходьба,3h10m",
		"15392,Бег,0h45m",
	}

	fmt.Println("\nЖурнал тренировок")
	fmt.Println("================")

	for i, training := range trainings {
		log.Printf("Обработка тренировки #%d: %s", i+1, training)
		info, err := spentcalories.TrainingInfo(training, weight, height)
		if err != nil {
			log.Printf("Ошибка обработки тренировки #%d: %v", i+1, err)
			continue
		}
		fmt.Println(info)
		fmt.Println("----------------")
	}
}
