package nextdate

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func NextDate(now time.Time, dstart string, repeat string) (string, error) {

	/*
	   now — время, от которого ищется ближайшая дата
	   dstart — исходное время в формате 20060102, от которого начинается отсчёт повторений
	   repeat — правило повторения
	*/

	pattern := `^\d+$` // Только цифры от начала до конца строки
	re := regexp.MustCompile(pattern)
	if len(repeat) == 0 {
		return "", errors.New("не задано повторение")
	}

	if len(dstart) != 0 && re.MatchString(dstart) == false {
		return "", errors.New("не верный формат даты")

	}

	date, err := time.Parse("20060102", dstart)

	if err != nil {
		return "", err
	}

	mass_interval := strings.Split(strings.TrimSpace(repeat), " ")

	if mass_interval[0] == "d" || mass_interval[0] == "y" {

		if mass_interval[0] == "d" {

			if len(strings.TrimSpace(repeat)) == 1 {
				return "", errors.New("не указан интервал в днях")
			}

			interval, err := strconv.Atoi(mass_interval[1])

			if err != nil {
				return "", err
			}

			if interval > 400 {
				return "", errors.New("превышен максимально допустимый интервал")
			}

			if interval == 0 {
				return "", errors.New("не указано количество дней повторения")
			}

			for {
				date = date.AddDate(0, 0, interval)
				if AfterNow(date, now) {
					break
				}
			}
		}

		if mass_interval[0] == "y" {

			for {
				date = date.AddDate(1, 0, 0)
				if AfterNow(date, now) {
					break
				}
			}
		}
	} else {
		return "", errors.New("указан неверный формат")

	}

	return date.Format("20060102"), nil

}

func AfterNow(date, now time.Time) bool {

	if now.Compare(date) < 0 {

		return true
	}
	return false
}
