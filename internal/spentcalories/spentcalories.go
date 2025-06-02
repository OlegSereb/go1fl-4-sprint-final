package spentcalories

import (
	"errors"
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

func parseTraining(data string) (int, string, time.Duration, error) {

	//Разделить строку на слайс строк.
	parts := strings.Split(data, ",")

	//Проверить, чтобы длина слайса была равна 3, так как в строке данных у нас количество шагов, вид активности и продолжительность.
	if len(parts) != 3 {
		return 0, "", 0, errors.New("некорректный формат данных: ожидается 'шаги,вид активности,длительность'")
	}

	//Преобразовать первый элемент слайса (количество шагов) в тип int. Обработать возможные ошибки. При их возникновении из функции вернуть 0 шагов, 0 продолжительность и ошибку.
	steps, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return 0, "", 0, fmt.Errorf("неверный формат количества шагов: %w", err)
	}
	if steps <= 0 {
		return 0, "", 0, errors.New("количество шагов должно быть положительным")
	}

	//Если всё прошло без ошибок, верните количество шагов, вид активности, продолжительность и nil (для ошибки).
	activity := strings.TrimSpace(parts[1])
	if activity != "Ходьба" && activity != "Бег" {
		return 0, "", 0, errors.New("неподдерживаемый вид активности: доступны 'Ходьба' или 'Бег'")
	}

	//Преобразовать третий элемент слайса в time.Duration. В пакете time есть метод для парсинга строки в time.Duration. Обработать возможные ошибки. При их возникновении из функции вернуть 0 шагов, 0 продолжительность и ошибку.
	duration, err := time.ParseDuration(strings.TrimSpace(parts[2]))
	if err != nil {
		return 0, "", 0, fmt.Errorf("неверный формат длительности: %w", err)
	}
	if duration <= 0 {
		return 0, "", 0, errors.New("длительность должна быть положительной")
	}

	return steps, activity, duration, nil
}

func distance(steps int, height float64) float64 {
	if steps <= 0 || height <= 0 {
		return 0
	}

	//рассчитайте длину шага. Для этого умножьте высоту пользователя на коэффициент длины шага stepLengthCoefficient. Соответствующая константа уже определена в пакете.
	stepLength := height * stepLengthCoefficient

	//умножьте пройденное количество шагов на длину шага.
	distanceMeters := float64(steps) * stepLength

	//разделите полученное значение на число метров в километре (mInKm, константа определена в пакете).
	distanceKm := distanceMeters / mInKm
	return distanceKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {

	//Проверить, что продолжительность duration больше 0. Если это не так, вернуть 0.
	if duration <= 0 || steps <= 0 || height <= 0 {
		return 0
	}
	return distance(steps, height) / duration.Hours()
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {

	//Проверить входные параметры на корректность. Если параметры некорректны, вернуть 0 калорий и соответствующую ошибку.
	if steps <= 0 {
		return 0, errors.New("количество шагов должно быть положительным")
	}
	if weight <= 0 {
		return 0, errors.New("вес должен быть положительным")
	}
	if height <= 0 {
		return 0, errors.New("рост должен быть положительным")
	}
	if duration <= 0 {
		return 0, errors.New("длительность должна быть положительной")
	}

	//Рассчитать среднюю скорость с помощью meanSpeed().
	speed := meanSpeed(steps, height, duration)
	return (weight * speed * duration.Minutes()) / minInH, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {

	//Проверить входные параметры на корректность. Если параметры некорректны, вернуть 0 калорий и соответствующую ошибку.
	if steps <= 0 {
		return 0, errors.New("количество шагов должно быть положительным")
	}
	if weight <= 0 {
		return 0, errors.New("вес должен быть положительным")
	}
	if height <= 0 {
		return 0, errors.New("рост должен быть положительным")
	}
	if duration <= 0 {
		return 0, errors.New("длительность должна быть положительной")
	}

	//Рассчитать среднюю скорость с помощью meanSpeed().
	speed := meanSpeed(steps, height, duration)

	//Рассчитать и количество калорий.
	calories := (weight * speed * duration.Minutes()) / minInH

	//Умножить полученное число калорий на корректирующий коэффициент walkingCaloriesCoefficient. Соответствующая константа объявляена в пакете. Вернуть полученное значение.
	return calories * walkingCaloriesCoefficient, nil
}

func TrainingInfo(data string, weight, height float64) (string, error) {

	if weight <= 0 || height <= 0 {
		return "", errors.New("вес и рост должны быть положительными")
	}
	// Получаем данные о тренировке с помощью функции parseTraining
	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		return "", fmt.Errorf("ошибка обработки данных: %w", err)
	}

	// Проверяем, какой вид тренировки был передан

	var calories float64
	var calcErr error

	switch activity {
	case "Ходьба":
		calories, calcErr = WalkingSpentCalories(steps, weight, height, duration)
	case "Бег":
		calories, calcErr = RunningSpentCalories(steps, weight, height, duration)
	default:
		return "", errors.New("неподдерживаемый тип тренировки")
	}

	if calcErr != nil {
		return "", fmt.Errorf("ошибка расчета калорий: %w", calcErr)
	}

	dist := distance(steps, height)
	speed := meanSpeed(steps, height, duration)

	return fmt.Sprintf(
		"Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f",
		activity,
		duration.Hours(),
		dist,
		speed,
		calories,
	), nil
}
