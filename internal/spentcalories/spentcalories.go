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
		return 0, "", 0, errors.New("некорректный формат данных (parseTraining)")
	}

	//Преобразовать первый элемент слайса (количество шагов) в тип int. Обработать возможные ошибки. При их возникновении из функции вернуть 0 шагов, 0 продолжительность и ошибку.
	steps, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return 0, "", 0, err
	}

	//Преобразовать третий элемент слайса в time.Duration. В пакете time есть метод для парсинга строки в time.Duration. Обработать возможные ошибки. При их возникновении из функции вернуть 0 шагов, 0 продолжительность и ошибку.
	duration, err := time.ParseDuration(strings.TrimSpace(parts[2]))
	if err != nil {
		return 0, "", 0, err
	}

	//Если всё прошло без ошибок, верните количество шагов, вид активности, продолжительность и nil (для ошибки).
	activity := strings.TrimSpace(parts[1])
	if activity != "Ходьба" && activity != "Бег" {
		return 0, "", 0, errors.New("неизвестный тип активности")
	}

	return steps, activity, duration, nil
}

func distance(steps int, height float64) float64 {

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
	if duration <= 0 || steps <= 0 {
		return 0
	}
	//Вычислить дистанцию с помощью distance().
	dist := distance(steps, height)

	//Вычислить и вернуть среднюю скорость. Для этого разделите дистанцию на продолжительность в часах. Чтобы перевести продолжительность в часы, воспользуйтесь функцией из пакета time.
	hours := duration.Hours()
	meanSpeed := dist / hours
	return meanSpeed
}

func TrainingInfo(data string, weight, height float64) (string, error) {

	// Получаем данные о тренировке с помощью функции parseTraining
	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		return "", err
	}

	// Проверяем, какой вид тренировки был передан
	var distanceKm, meanSpd, calories float64

	switch activity {
	case "Ходьба":
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}
		distanceKm = distance(steps, height)
		meanSpd = meanSpeed(steps, height, duration)

	case "Бег":
		calories, err = RunningSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}
		distanceKm = distance(steps, height)
		meanSpd = meanSpeed(steps, height, duration)

	default:
		return "", errors.New("неизвестный тип тренировки")
	}

	// Формируем строку с информацией о тренировке
	result := fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f",
		activity, duration.Hours(), distanceKm, meanSpd, calories)

	return result, nil

}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {

	//Проверить входные параметры на корректность. Если параметры некорректны, вернуть 0 калорий и соответствующую ошибку.
	if steps <= 0 {
		return 0, errors.New("количество шагов должно быть больше 0")
	}
	if weight <= 0 {
		return 0, errors.New("вес пользователя должен быть больше 0")
	}
	if duration <= 0 {
		return 0, errors.New("продолжительность должна быть больше 0")
	}

	//Рассчитать среднюю скорость с помощью meanSpeed().
	avgSpeed := meanSpeed(steps, height, duration)
	durationInMinutes := duration.Minutes()

	//Рассчитать и вернуть количество калорий.
	calories := (weight * avgSpeed * durationInMinutes) / minInH

	return calories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {

	//Проверить входные параметры на корректность. Если параметры некорректны, вернуть 0 калорий и соответствующую ошибку.
	if steps <= 0 {
		return 0, errors.New("количество шагов должно быть больше 0")
	}
	if weight <= 0 {
		return 0, errors.New("вес пользователя должен быть больше 0")
	}
	if height <= 0 {
		return 0, errors.New("рост пользователя должен быть больше 0")
	}
	if duration <= 0 {
		return 0, errors.New("продолжительность должна быть больше 0")
	}

	//Рассчитать среднюю скорость с помощью meanSpeed().
	avgSpeed := meanSpeed(steps, height, duration)
	durationInMinutes := duration.Minutes()

	//Рассчитать и количество калорий.
	calories := (weight * avgSpeed * durationInMinutes) / minInH

	//Умножить полученное число калорий на корректирующий коэффициент walkingCaloriesCoefficient. Соответствующая константа объявляена в пакете. Вернуть полученное значение.
	calories *= walkingCaloriesCoefficient

	return calories, nil
}
