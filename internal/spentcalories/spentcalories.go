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
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {
	s := strings.Split(data, ",")
	if len(s) != 3 {
		return 0, "", 0, errors.New("input data must contain 3 fields separated by commas")
	}

	steps, err := strconv.Atoi(s[0])
	if err != nil {
		return 0, "", 0, fmt.Errorf("invalid steps format: %w", err)
	}
	if steps <= 0 {
		return 0, "", 0, errors.New("steps must be greater than 0")
	}

	duration, err := time.ParseDuration(s[2])
	if err != nil {
		return 0, "", 0, fmt.Errorf("invalid duration format: %w", err)
	}
	if duration <= 0 {
		return 0, "", 0, errors.New("duration must be greater than 0")
	}

	return steps, s[1], duration, nil
}

func distance(steps int, height float64) float64 {
	switch {
	case steps <= 0:
		return 0
	case height <= 0:
		return 0
	default:
		return height * stepLengthCoefficient * float64(steps) / mInKm
	}
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	switch {
	case steps <= 0:
		return 0
	case height <= 0:
		return 0
	case duration <= 0:
		return 0
	default:
		return distance(steps, height) / duration.Hours()
	}
}

func isInputDataCorrect(steps int, weight, height float64, duration time.Duration) error {
	switch {
	case steps <= 0:
		return errors.New("steps must be greater than 0")
	case weight <= 0:
		return errors.New("weight must be greater than 0")
	case height <= 0:
		return errors.New("height must be greater than 0")
	case duration <= 0:
		return errors.New("duration must be greater than 0")
	default:
		return nil
	}
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	err := isInputDataCorrect(steps, weight, height, duration)
	if err != nil {
		return 0, fmt.Errorf("invalid input data: %w", err)
	}

	return (weight * meanSpeed(steps, height, duration) * duration.Minutes()) / float64(minInH), nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	err := isInputDataCorrect(steps, weight, height, duration)
	if err != nil {
		return 0, fmt.Errorf("invalid input data: %w", err)
	}

	return ((weight * meanSpeed(steps, height, duration) * duration.Minutes()) / float64(minInH)) * walkingCaloriesCoefficient, nil
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, typeOfTrening, duration, err := parseTraining(data)
	if err != nil {
		log.Printf("не получилось получить информацию о тренировке: %v", err)
		return "", err
	}

	if weight <= 0 {
		return "", errors.New("weight must be greater than 0")
	}

	if height <= 0 {
		return "", errors.New("height must be greater than 0")
	}

	switch typeOfTrening {
	case "Бег":
		runningSpentCalories, err := RunningSpentCalories(steps, weight, height, duration)
		if err != nil {
			log.Printf("не получилось получить информацию о тренировке: %v", err)
			return "", err
		}
		return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
			typeOfTrening, duration.Hours(), distance(steps, height), meanSpeed(steps, height, duration), runningSpentCalories), nil
	case "Ходьба":
		walkingSpentCalories, err := WalkingSpentCalories(steps, weight, height, duration)
		if err != nil {
			log.Printf("не получилось получить информацию о тренировке: %v", err)
			return "", err
		}
		return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
			typeOfTrening, duration.Hours(), distance(steps, height), meanSpeed(steps, height, duration), walkingSpentCalories), nil
	default:
		return "", errors.New("неизвестный тип тренировки")
	}
}
