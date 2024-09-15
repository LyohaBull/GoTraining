package main

import (
	"fmt"
	"strconv"
	"time"
)

func getWeekDay(date time.Time) string {
	switch date.Weekday().String() {
	case "Monday":
		{
			return "ПН"
		}
	case "Tuesday":
		{
			return "ВТ"
		}
	case "Wednesday":
		{
			return "СР"
		}
	case "Thursday":
		{
			return "ЧТ"
		}
	case "Friday":
		{
			return "ПТ"
		}
	case "Saturday":
		{
			return "СБ"
		}
	case "Sunday":
		{
			return "ВС"
		}
	default:
		{
			return "nil"
		}
	}
}
func getDayAgo(date time.Time, dif int) time.Time {
	return date.AddDate(0, 0, -dif)
}
func getLastDayOfMonth(date time.Time) int {
	t := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, date.UTC().Location())
	s := t.AddDate(0, 1, 0)
	difference := (s.Sub(t)) / (24 * 1000 * 60 * 60 * 1000000)
	return int(difference)
}

func getSecondsToday() int {
	now := time.Now()
	begin := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	return int(now.Sub(begin).Seconds())
}
func getSecondsTomorrow() int {
	now := time.Now()
	end := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
	return int(end.Sub(now).Seconds())
}
func formatDate(date time.Time) string {
	now := time.Now()
	duration := now.Sub(date)
	if duration.Milliseconds() < 1000 {
		return "прямо сейчас"
	}
	if duration.Seconds() < 60 {
		return strconv.Itoa(int(duration.Seconds())) + " сек. назад"
	}
	if duration.Minutes() < 60 {
		return strconv.Itoa(int(duration.Minutes())) + " мин. назад"
	}
	return date.Format("01.02.2006")
}
func main() {
	t, _ := time.Parse(time.RFC3339, "2012-02-20T03:12:00+03:00")
	fmt.Println(t)
	fmt.Println(getWeekDay(t))
	date := time.Date(2015, 1, 2, 0, 0, 0, 0, time.UTC)
	fmt.Println(date)
	fmt.Println(getDayAgo(date, 365))
	fmt.Println(getLastDayOfMonth(time.Date(2023, 2, 3, 0, 0, 0, 0, time.UTC)))
	fmt.Println(getSecondsToday())
	fmt.Println(getSecondsTomorrow())
	fmt.Println(formatDate(time.Now().AddDate(-1, -2, -45)))
}
