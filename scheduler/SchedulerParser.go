package scheduler

import (
	"github.com/robfig/cron/v3"
	"strings"
	"time"
)

func (tr TimeRange) GetScheduleSpec(targetTime time.Time) (nextWindowEdge time.Time, isTimeBetween bool, err error) {
	err = tr.ValidateTimeRange()
	if err != nil {
		return nextWindowEdge, false, err
	}
	if tr.Frequency == FIXED {
		nextWindowEdge, isTimeBetween = getScheduleForFixedTime(targetTime, tr)
		return nextWindowEdge, isTimeBetween, err
	}
	month, year := tr.getMonthAndYear(targetTime)
	cronExp := tr.getCronExp(year, month)
	parser := cron.NewParser(CRON)
	schedule, err := parser.Parse(cronExp)
	if err != nil {
		return nextWindowEdge, false, err
	}
	duration, err := tr.getDuration(month, year)
	if err != nil {
		return nextWindowEdge, false, err
	}

	windowStart, windowEnd := tr.getWindowStartAndEndTime(targetTime, duration, schedule)
	if isTimeInBetween(targetTime, windowStart, windowEnd) {
		return windowEnd, true, err
	}
	return windowStart, false, err
}

func (tr TimeRange) getWindowStartAndEndTime(targetTime time.Time, duration time.Duration, schedule cron.Schedule) (time.Time, time.Time) {
	var windowEnd time.Time

	prevDuration := duration
	if tr.isCyclic() {
		diff := getLastDayOfMonth(targetTime.Year(), targetTime.Month()) - getLastDayOfMonth(targetTime.Year(), targetTime.Month()-1)
		prevDuration = duration - time.Duration(diff)*time.Hour*24
	}

	timeMinusDuration := targetTime.Add(-1 * prevDuration)
	windowStart := schedule.Next(timeMinusDuration)
	windowEnd = windowStart.Add(duration)
	if !tr.TimeFrom.IsZero() && windowStart.Before(tr.TimeFrom) {
		windowStart = tr.TimeFrom
	}
	if !tr.TimeTo.IsZero() && windowEnd.After(tr.TimeTo) {
		windowEnd = tr.TimeTo
	}
	return windowStart, windowEnd
}

func (tr TimeRange) getCronExp(year int, month time.Month) string {
	cronExp := tr.getCron()
	lastDayOfMonth := getLastDayOfMonth(year, month)
	if strings.Contains(cronExp, "L-2") {
		lastDayOfMonth = lastDayOfMonth - 2
		cronExp = strings.Replace(cronExp, "L-2", intToString(lastDayOfMonth), -1)
	} else if strings.Contains(cronExp, "L-1") {
		lastDayOfMonth = lastDayOfMonth - 1
		cronExp = strings.Replace(cronExp, "L-1", intToString(lastDayOfMonth), -1)
	} else {
		cronExp = strings.Replace(cronExp, "L", intToString(lastDayOfMonth), -1)
	}
	return cronExp
}

func (tr TimeRange) getMonthAndYear(targetTime time.Time) (time.Month, int) {
	month := targetTime.Month()
	year := targetTime.Year()
	day := targetTime.Day()

	isBeforeEndTime, err := isToHourMinuteBefore(tr, targetTime)
	if err != nil {
		return 0, 0
	}
	if day >= 1 && (day < tr.DayTo || (day == tr.DayTo && isBeforeEndTime)) && tr.isCyclic() {
		if month == 1 {
			month = 12
			year = year - 1
		} else {
			month = month - 1
		}
	}
	return month, year
}

func isTimeInBetween(timeCurrent, periodStart, periodEnd time.Time) bool {
	return (timeCurrent.After(periodStart) && timeCurrent.Before(periodEnd)) || timeCurrent.Equal(periodStart)
}
