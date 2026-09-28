package spentcalories

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	// lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {
	s := strings.Split(data, ",")
	if len(s) != 3 {
		return 0, "", 0, errors.New("data должно быть 3 поля, разделенных запятыми")
	}

	steps, err := strconv.Atoi(s[0])
	if err != nil {
		return 0, "", 0, err
	}

	if steps <= 0 {
		return 0, "", 0, errors.New("steps должно быть больше 0")
	}

	duration, err := time.ParseDuration(s[2])
	if err != nil {
		return 0, "", 0, err
	}

	if duration <= 0 {
		return 0, "", 0, errors.New("duration должно быть больше 0")
	}

	return steps, s[1], duration, nil
}

func distance(steps int, height float64) float64 {
	stepLength := height * stepLengthCoefficient
	dist := float64(steps) * stepLength
	km := dist / float64(mInKm)
	return km
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}

	dist := distance(steps, height)
	speed := dist / float64(duration.Hours())

	return speed
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, typeOfTrening, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}

	dist := distance(steps, height)
	speed := meanSpeed(steps, height, duration)
	caloriesRun, err := RunningSpentCalories(steps, weight, height, duration)
	if err != nil {
		log.Println(err)
		return "", err
	}

	caloriesWalk, err := WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		log.Println(err)
		return "", err
	}

	hours := duration.Hours()
	switch typeOfTrening {
	case "Бег":
		return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n", typeOfTrening, hours, dist, speed, caloriesRun), nil
	case "Ходьба":
		return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n", typeOfTrening, hours, dist, speed, caloriesWalk), nil
	default:
		return "", errors.New("неизвестный тип тренировки")
	}
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, errors.New("steps can't be less 1")
	}

	if weight <= 0 {
		return 0, errors.New("weight can't be less 1")
	}

	if height <= 0 {
		return 0, errors.New("height can't be less 1")
	}

	if duration <= 0 {
		return 0, errors.New("duration can't be less 1")
	}

	speed := meanSpeed(steps, height, duration)
	durationInMinutes := duration.Minutes()
	result := (weight * speed * durationInMinutes) / float64(minInH)

	return result, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, errors.New("steps can't be less 1")
	}

	if weight <= 0 {
		return 0, errors.New("weight can't be less 1")
	}

	if height <= 0 {
		return 0, errors.New("height can't be less 1")
	}

	if duration <= 0 {
		return 0, errors.New("duration can't be less 1")
	}

	speed := meanSpeed(steps, height, duration)
	durationInMinutes := duration.Minutes()
	result := (weight * speed * durationInMinutes) / float64(minInH)
	walkingCalories := result * walkingCaloriesCoefficient

	return walkingCalories, nil
}
