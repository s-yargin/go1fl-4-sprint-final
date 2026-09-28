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
		return 0, time.Duration(0), errors.New("data должно быть 2 поля, разделенных запятыми")
	}

	steps, err := strconv.Atoi(s[0])
	if err != nil {
		return 0, time.Duration(0), err
	}

	if steps <= 0 {
		return 0, time.Duration(0), errors.New("steps должно быть больше 0")
	}

	duration, err := time.ParseDuration(s[1])
	if err != nil {
		return 0, time.Duration(0), err
	}

	if duration <= 0 {
		return 0, 0, errors.New("duration должно быть больше 0")
	}

	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)
	if err != nil {
		err = fmt.Errorf("во время выполнения программы произошла ошибка: %v", err)
		fmt.Println(err)
		return ""
	}

	if steps <= 0 {
		return ""
	}

	kMet := float64(steps) * stepLength / float64(mInKm)
	kCal, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", steps, kMet, kCal)
	}

	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", steps, kMet, kCal)
}
