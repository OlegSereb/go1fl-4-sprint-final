package daysteps

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	// Разделяем data по ','
	parts := strings.Split(data, ",")

	// Проверить, чтобы длина слайса была равна 2.
	if len(parts) != 2 {
		return 0, 0, errors.New("некорректный формат данных (ожидается: шаги,длительность)")
	}

	// Преобразовать первый элемент слайса (количество шагов) в тип int.
	steps, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return 0, 0, fmt.Errorf("неверный формат шагов: %w", err)
	}

	// Проверить: количество шагов должно быть больше 0.
	if steps <= 0 {
		return 0, 0, errors.New("количество шагов должно быть больше 0 (parsePackage)")
	}

	// Преобразовать второй элемент слайса в time.Duration.
	duration, err := time.ParseDuration(strings.TrimSpace(parts[1]))
	if err != nil {
		return 0, 0, fmt.Errorf("неверный формат длительности: %w", err)
	}

	// Проверка на нулевую продолжительность
	if duration <= 0 {
		return 0, 0, errors.New("длительность должна быть больше 0")
	}

	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	// Значения должны быть положительными.
	if weight <= 0 || height <= 0 {
		return "Ошибка: вес и рост должны быть положительными числами"
	}

	// Получить данные о количестве шагов и продолжительности прогулки
	steps, duration, err := parsePackage(data)
	if err != nil {
		return fmt.Sprintf("Ошибка: %v", err)
	}

	// Вычислить дистанцию в метрах
	distanceMeters := float64(steps) * stepLength

	// Перевести дистанцию в километры
	distanceKm := distanceMeters / mInKm

	// Вычислить количество калорий, потраченных на прогулке
	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		return fmt.Sprintf("Ошибка при расчете калорий: %v", err)
	}

	// Сформировать строку с результатами
	return fmt.Sprintf(
		"Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n",
		steps, distanceKm, calories,
	)
}
