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
	lenStep   = 0.65  // средняя длина шага.
	mInKm     = 1000  // количество метров в километре.
	minInH    = 60    // количество минут в часе.
	kmhInMsec = 0.278 // коэффициент для преобразования км/ч в м/с.
	cmInM     = 100   // количество сантиметров в метре.

)

func parseTraining(data string) (int, string, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 3 {
		return 0, "", 0, errors.New("неверный формат данных")
	}

	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, "", 0, err
	}

	duration, err := time.ParseDuration(parts[2])
	if err != nil {
		return 0, "", 0, err
	}

	return steps, parts[1], duration, nil
}

func distance(steps int) float64 {
	return float64(steps) * lenStep / float64(mInKm)
}

func meanSpeed(steps int, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}
	dist := distance(steps)
	return dist / duration.Hours()
}

func TrainingInfo(data string, weight, height float64) string {
	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	dist := distance(steps)
	speed := meanSpeed(steps, duration)

	switch activity {
	case "Бег":
		calories := RunningSpentCalories(steps, weight, duration)
		return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f",
			activity, duration.Hours(), dist, speed, calories)
	case "Ходьба":
		calories := WalkingSpentCalories(steps, weight, height, duration)
		return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f",
			activity, duration.Hours(), dist, speed, calories)
	default:
		return "неизвестный тип тренировки"
	}

}

const (
	runningCaloriesMeanSpeedMultiplier = 18.0 // множитель средней скорости.
	runningCaloriesMeanSpeedShift      = 20.0 // среднее количество сжигаемых калорий при беге.
)

func RunningSpentCalories(steps int, weight float64, duration time.Duration) float64 {

	speed := meanSpeed(steps, duration)
	return (runningCaloriesMeanSpeedMultiplier*speed - runningCaloriesMeanSpeedShift) * weight
}

const (
	walkingCaloriesWeightMultiplier = 0.035 // множитель массы тела.
	walkingSpeedHeightMultiplier    = 0.029 // множитель роста.
)

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) float64 {

	speed := meanSpeed(steps, duration)
	heightCm := height * float64(cmInM)
	return (walkingCaloriesWeightMultiplier*weight + (speed*speed/heightCm)*walkingSpeedHeightMultiplier) * duration.Hours() * float64(minInH)
}
