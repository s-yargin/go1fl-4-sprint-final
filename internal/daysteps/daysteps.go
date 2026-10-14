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
	s := strings.Split(data, ",")
	if len(s) != 2 {
		return 0, 0, errors.New("input data must contain 2 fields separated by commas")
	}

	steps, err := strconv.Atoi(s[0])
	if err != nil {
		return 0, 0, err
	}
	if steps <= 0 {
		return 0, 0, errors.New("steps must be greater than 0")
	}

	duration, err := time.ParseDuration(s[1])
	if err != nil {
		return 0, 0, err
	}
	if duration <= 0 {
		return 0, 0, errors.New("duration must be greater than 0")
	}

	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)
	if err != nil {
		err = fmt.Errorf("не получилось получить информацию об активности: %v", err)
		fmt.Println(err)
		return ""
	}

	if weight <= 0 {
		err = errors.New("weight must be greater than 0")
		fmt.Println(err)
		return ""
	}

	if height <= 0 {
		err = errors.New("height must be greater than 0")
		fmt.Println(err)
		return ""
	}

	distance := float64(steps) * stepLength / mInKm
	walkingSpentCalories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		fmt.Println(err)
	}

	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", steps, distance, walkingSpentCalories)
}
