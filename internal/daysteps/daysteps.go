package daysteps

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"go1fl-4-sprint-final/internal/spentcalories"
)

var (
	StepLength = 0.65 // длина шага в метрах
)

func parsePackage(data string) (int, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 2 {
		return 0, 0, errors.New("неверный формат данных")
	}

	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, err
	}

	duration, err := time.ParseDuration(parts[1])
	if err != nil {
		return 0, 0, err
	}

	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	if steps <= 0 {
		return ""
	}

	distance := float64(steps) * StepLength / 1000.0
	calories := spentcalories.WalkingSpentCalories(steps, weight, height, duration)

	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.",
		steps, distance, calories)
}
