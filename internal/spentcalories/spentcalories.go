package spentcalories

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

// parseTraining парсит строку формата "3456,Ходьба,3h00m".
func parseTraining(data string) (int, string, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 3 {
		return 0, "", 0, fmt.Errorf("неверный формат данных: требуется разделение на 3 части через запятую")
	}

	stepsStr := parts[0]
	activity := parts[1]
	durationStr := parts[2]

	steps, err := strconv.Atoi(stepsStr)
	if err != nil {
		return 0, "", 0, fmt.Errorf("ошибка при парсинге количества шагов: %w", err)
	}
	if steps <= 0 {
		return 0, "", 0, fmt.Errorf("количество шагов должно быть больше 0")
	}

	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return 0, "", 0, fmt.Errorf("ошибка при парсинге продолжительности: %w", err)
	}
	if duration <= 0 {
		return 0, "", 0, fmt.Errorf("продолжительность тренировки должна быть больше нуля")
	}

	return steps, activity, duration, nil
}

// distance рассчитывает дистанцию в километрах на основе количества шагов и роста пользователя.
func distance(steps int, height float64) float64 {
	if steps <= 0 || height <= 0 {
		return 0
	}
	stepLength := height * stepLengthCoefficient
	totalMeters := float64(steps) * stepLength
	return totalMeters / mInKm
}

// meanSpeed рассчитывает среднюю скорость в км/ч.
func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}
	dist := distance(steps, height)
	hours := duration.Hours()
	if hours <= 0 {
		return 0
	}
	return dist / hours
}

// RunningSpentCalories рассчитывает калории при беге.
func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, fmt.Errorf("некорректные входные параметры")
	}

	speed := meanSpeed(steps, height, duration)
	minutes := duration.Minutes()

	calories := (weight * speed * minutes) / minInH
	return calories, nil
}

// WalkingSpentCalories рассчитывает калории при ходьбе.
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, fmt.Errorf("некорректные входные параметры")
	}

	baseCalories, _ := RunningSpentCalories(steps, weight, height, duration)
	calories := baseCalories * walkingCaloriesCoefficient
	return calories, nil
}

// TrainingInfo возвращает полную информацию о тренировке.
func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		return "", err
	}

	var calories float64
	var errCalc error

	switch activity {
	case "Бег":
		calories, errCalc = RunningSpentCalories(steps, weight, height, duration)
	case "Ходьба":
		calories, errCalc = WalkingSpentCalories(steps, weight, height, duration)
	default:
		return "", fmt.Errorf("неизвестный тип тренировки: %s", activity)
	}

	if errCalc != nil {
		return "", errCalc
	}

	speed := meanSpeed(steps, height, duration)
	hours := duration.Hours()
	distanceKm := distance(steps, height)

	info := fmt.Sprintf(`Тип тренировки: %s
Длительность: %.2f ч.
Дистанция: %.2f км.
Скорость: %.2f км/ч
Сожгли калорий: %.2f
`, activity, hours, distanceKm, speed, calories)

	return info, nil
}
