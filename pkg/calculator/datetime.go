package calculator

import (
	"fmt"
	"time"
)

type DateTimeCalc struct {
	StartDate time.Time
	EndDate   time.Time
	Result    string
}

func NewDateTimeCalc() *DateTimeCalc {
	now := time.Now()
	return &DateTimeCalc{
		StartDate: now,
		EndDate:   now,
	}
}

type DateDifference struct {
	Years        int
	Months       int
	Days         int
	TotalDays    int
	TotalWeeks   int
	TotalHours   int
	TotalMinutes int
	TotalSeconds int
}

func (d *DateTimeCalc) CalculateDifference() DateDifference {
	start := d.StartDate
	end := d.EndDate

	if end.Before(start) {
		start, end = end, start
	}

	duration := end.Sub(start)
	totalDays := int(duration.Hours() / 24)
	totalHours := int(duration.Hours())
	totalMinutes := int(duration.Minutes())
	totalSeconds := int(duration.Seconds())
	totalWeeks := totalDays / 7

	years := end.Year() - start.Year()
	months := int(end.Month()) - int(start.Month())
	days := end.Day() - start.Day()

	if days < 0 {
		months--
		lastMonth := end.AddDate(0, -1, 0)
		days += daysInMonth(lastMonth.Year(), lastMonth.Month())
	}
	if months < 0 {
		years--
		months += 12
	}

	return DateDifference{
		Years:        years,
		Months:       months,
		Days:         days,
		TotalDays:    totalDays,
		TotalWeeks:   totalWeeks,
		TotalHours:   totalHours,
		TotalMinutes: totalMinutes,
		TotalSeconds: totalSeconds,
	}
}

func daysInMonth(year int, month time.Month) int {
	return time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

func (d *DateTimeCalc) AddDays(days int) time.Time {
	return d.StartDate.AddDate(0, 0, days)
}

func (d *DateTimeCalc) AddWeeks(weeks int) time.Time {
	return d.StartDate.AddDate(0, 0, weeks*7)
}

func (d *DateTimeCalc) AddMonths(months int) time.Time {
	return d.StartDate.AddDate(0, months, 0)
}

func (d *DateTimeCalc) AddYears(years int) time.Time {
	return d.StartDate.AddDate(years, 0, 0)
}

func (d *DateTimeCalc) AddTime(years, months, days, hours, minutes, seconds int) time.Time {
	result := d.StartDate.AddDate(years, months, days)
	result = result.Add(time.Duration(hours) * time.Hour)
	result = result.Add(time.Duration(minutes) * time.Minute)
	result = result.Add(time.Duration(seconds) * time.Second)
	return result
}

func (d *DateTimeCalc) SubtractDays(days int) time.Time {
	return d.StartDate.AddDate(0, 0, -days)
}

func (d *DateTimeCalc) SubtractWeeks(weeks int) time.Time {
	return d.StartDate.AddDate(0, 0, -weeks*7)
}

func (d *DateTimeCalc) SubtractMonths(months int) time.Time {
	return d.StartDate.AddDate(0, -months, 0)
}

func (d *DateTimeCalc) SubtractYears(years int) time.Time {
	return d.StartDate.AddDate(-years, 0, 0)
}

func (d *DateTimeCalc) GetWeekday() time.Weekday {
	return d.StartDate.Weekday()
}

func (d *DateTimeCalc) GetWeekNumber() int {
	_, week := d.StartDate.ISOWeek()
	return week
}

func (d *DateTimeCalc) GetDayOfYear() int {
	return d.StartDate.YearDay()
}

func (d *DateTimeCalc) IsLeapYear() bool {
	year := d.StartDate.Year()
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

func (d *DateTimeCalc) DaysUntilEndOfYear() int {
	endOfYear := time.Date(d.StartDate.Year(), 12, 31, 23, 59, 59, 0, d.StartDate.Location())
	return int(endOfYear.Sub(d.StartDate).Hours() / 24)
}

func (d *DateTimeCalc) DaysUntilEndOfMonth() int {
	year := d.StartDate.Year()
	month := d.StartDate.Month()
	lastDay := daysInMonth(year, month)
	return lastDay - d.StartDate.Day()
}

func (d *DateTimeCalc) GetAge(birthDate time.Time) (years, months, days int) {
	now := time.Now()
	years = now.Year() - birthDate.Year()
	months = int(now.Month()) - int(birthDate.Month())
	days = now.Day() - birthDate.Day()

	if days < 0 {
		months--
		lastMonth := now.AddDate(0, -1, 0)
		days += daysInMonth(lastMonth.Year(), lastMonth.Month())
	}
	if months < 0 {
		years--
		months += 12
	}
	return
}

func (d *DateTimeCalc) GetWorkingDays(excludeWeekends bool) int {
	start := d.StartDate
	end := d.EndDate

	if end.Before(start) {
		start, end = end, start
	}

	count := 0
	for current := start; !current.After(end); current = current.AddDate(0, 0, 1) {
		if excludeWeekends {
			weekday := current.Weekday()
			if weekday != time.Saturday && weekday != time.Sunday {
				count++
			}
		} else {
			count++
		}
	}
	return count
}

func (d *DateTimeCalc) GetNextWeekday(weekday time.Weekday) time.Time {
	current := d.StartDate
	for {
		current = current.AddDate(0, 0, 1)
		if current.Weekday() == weekday {
			return current
		}
	}
}

func (d *DateTimeCalc) GetPreviousWeekday(weekday time.Weekday) time.Time {
	current := d.StartDate
	for {
		current = current.AddDate(0, 0, -1)
		if current.Weekday() == weekday {
			return current
		}
	}
}

type TimeDifference struct {
	Hours        int
	Minutes      int
	Seconds      int
	TotalHours   float64
	TotalMinutes float64
	TotalSeconds float64
}

func (d *DateTimeCalc) CalculateTimeDifference(startTime, endTime time.Time) TimeDifference {
	if endTime.Before(startTime) {
		startTime, endTime = endTime, startTime
	}

	duration := endTime.Sub(startTime)
	totalSeconds := duration.Seconds()
	totalMinutes := duration.Minutes()
	totalHours := duration.Hours()

	hours := int(totalHours)
	minutes := int(totalMinutes) % 60
	seconds := int(totalSeconds) % 60

	return TimeDifference{
		Hours:        hours,
		Minutes:      minutes,
		Seconds:      seconds,
		TotalHours:   totalHours,
		TotalMinutes: totalMinutes,
		TotalSeconds: totalSeconds,
	}
}

func (d *DateTimeCalc) UnixTimestamp() int64 {
	return d.StartDate.Unix()
}

func (d *DateTimeCalc) FromUnixTimestamp(timestamp int64) time.Time {
	return time.Unix(timestamp, 0)
}

func (d *DateTimeCalc) FormatDate(format string) string {
	return d.StartDate.Format(format)
}

func FormatDifference(diff DateDifference) string {
	parts := []string{}
	if diff.Years > 0 {
		parts = append(parts, fmt.Sprintf("%d year(s)", diff.Years))
	}
	if diff.Months > 0 {
		parts = append(parts, fmt.Sprintf("%d month(s)", diff.Months))
	}
	if diff.Days > 0 {
		parts = append(parts, fmt.Sprintf("%d day(s)", diff.Days))
	}
	if len(parts) == 0 {
		return "0 days"
	}
	result := ""
	for i, part := range parts {
		if i > 0 {
			result += ", "
		}
		result += part
	}
	return result
}

func (d *DateTimeCalc) SetStartDate(year, month, day int) {
	d.StartDate = time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
}

func (d *DateTimeCalc) SetEndDate(year, month, day int) {
	d.EndDate = time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
}

func (d *DateTimeCalc) SetStartDateTime(year, month, day, hour, minute, second int) {
	d.StartDate = time.Date(year, time.Month(month), day, hour, minute, second, 0, time.Local)
}

func (d *DateTimeCalc) SetEndDateTime(year, month, day, hour, minute, second int) {
	d.EndDate = time.Date(year, time.Month(month), day, hour, minute, second, 0, time.Local)
}

func (d *DateTimeCalc) Today() {
	d.StartDate = time.Now()
}
