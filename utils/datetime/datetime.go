package datetime

import (
	"errors"
	"gin-auth-mongo/utils/consts"
	"regexp"
	"strconv"
	"time"
)

// map
type TIME_UNIT time.Duration

var (
	MILLISECONDS TIME_UNIT = TIME_UNIT(time.Millisecond)
	SECONDS      TIME_UNIT = TIME_UNIT(time.Second)
	MINUTES      TIME_UNIT = TIME_UNIT(time.Minute)
	HOURS        TIME_UNIT = TIME_UNIT(time.Hour)
	DAYS         TIME_UNIT = TIME_UNIT(time.Hour * 24)

	WEEKS  TIME_UNIT = TIME_UNIT(time.Hour * 24 * 7)
	MONTHS TIME_UNIT = TIME_UNIT(time.Hour * 24 * 30)
	YEARS  TIME_UNIT = TIME_UNIT(time.Hour * 24 * 365)
)

// Get current time in datetime_nano_format
func GetCurrentTime() string {
	return time.Now().Format(consts.DATETIME_NANO_FORMAT)
}

func GetCurrentTimeUnix() int64 {
	return time.Now().Unix()
}

func GetCurrentTimeUnixString() string {
	return strconv.FormatInt(GetCurrentTimeUnix(), 10)
}

func ConvertUnixToDateTime(unix int64) string {
	return time.Unix(unix, 0).Format(consts.DATETIME_NANO_FORMAT)
}

// Get current date in date_format
func GetCurrentDate() string {
	return time.Now().Format(consts.DATE_FORMAT)
}

func IsTargetBeforeToday(targetDate string) (bool, error) {
	d, err := time.Parse(consts.DATE_FORMAT, targetDate)
	if err != nil {
		return false, errors.New("invalid date format")
	}
	return d.Before(time.Now()), nil
}

// check if current date is before date
func IsTodayBeforeTarget(targetDate string) (bool, error) {
	d, err := time.Parse(consts.DATE_FORMAT, targetDate)
	if err != nil {
		return false, errors.New("invalid date format")
	}
	return d.Before(time.Now()), nil
}

func IsTodayEqualTarget(targetDate string) (bool, error) {
	d, err := time.Parse(consts.DATE_FORMAT, targetDate)
	if err != nil {
		return false, errors.New("invalid date format")
	}
	return d.Equal(time.Now()), nil
}

// check if current date is after date
func IsTodayAfterTarget(targetDate string) (bool, error) {
	d, err := time.Parse(consts.DATE_FORMAT, targetDate)
	if err != nil {
		return false, errors.New("invalid date format")
	}
	return d.After(time.Now()), nil
}

func IsToday(date string) (bool, error) {
	d, err := time.Parse(consts.DATE_FORMAT, date)
	if err != nil {
		return false, errors.New("invalid date format")
	}
	return d.Equal(time.Now()), nil
}

// Date Format: YYYY-MM-DD
// check if date1 is before date2
func CompareDate(date1 string, date2 string) (bool, error) {
	d1, err1 := time.Parse(consts.DATE_FORMAT, date1)
	d2, err2 := time.Parse(consts.DATE_FORMAT, date2)
	if err1 != nil || err2 != nil {
		return false, errors.New("invalid date format")
	}
	return d1.Before(d2), nil
}

func IsValidRangeDate(startDate string, endDate string) (bool, error) {
	if startDate == "" {
		return false, errors.New("start date is empty")
	}

	if endDate == "" {
		endDate = "2099-12-31"
	}

	if isBefore, err := IsTargetBeforeToday(startDate); err != nil || isBefore {
		return false, errors.New("start date is before today")
	}

	if _, err := CompareDate(startDate, endDate); err != nil {
		return false, errors.New("start date is after end date")
	}
	return true, nil
}

func IsTodayInRange(startDate string, endDate string) (bool, error) {
	return InRangeDate(GetCurrentDate(), startDate, endDate)
}

func InRangeDate(date string, startDate string, endDate string) (bool, error) {

	if startDate == "" {
		return false, errors.New("start date is empty")
	}

	if endDate == "" {
		endDate = "2099-12-31"
	}

	if date == "" {
		date = GetCurrentDate()
	}

	d, err := time.Parse(consts.DATE_FORMAT, date)
	d1, err1 := time.Parse(consts.DATE_FORMAT, startDate)
	d2, err2 := time.Parse(consts.DATE_FORMAT, endDate)

	if err != nil || err1 != nil || err2 != nil {
		return false, errors.New("invalid date format")
	}

	return d.After(d1) && d.Before(d2), nil
}

func IsBeforeRangeDate(date string, startDate string, endDate string) (bool, error) {
	if startDate == "" {
		return false, errors.New("start date or end date is empty")
	}

	if endDate == "" {
		endDate = "2099-12-31"
	}

	if date == "" {
		date = GetCurrentDate()
	}

	d, err := time.Parse(consts.DATE_FORMAT, date)
	d1, err1 := time.Parse(consts.DATE_FORMAT, startDate)

	if err != nil || err1 != nil {
		return false, errors.New("invalid date format")
	}
	return d.Before(d1), nil
}

func IsAfterRangeDate(date string, startDate string, endDate string) (bool, error) {
	if startDate == "" || endDate == "" {
		return false, errors.New("start date or end date is empty")
	}

	if date == "" {
		date = GetCurrentDate()
	}

	d, err := time.Parse(consts.DATE_FORMAT, date)
	d2, err2 := time.Parse(consts.DATE_FORMAT, endDate)

	if err != nil || err2 != nil {
		return false, errors.New("invalid date format")
	}
	return d.After(d2), nil
}

func InRangeDateTime(date string, startDate string, endDate string) (bool, error) {
	if startDate == "" {
		return false, errors.New("start date is empty")
	}

	if endDate == "" {
		endDate = "2099-12-31T23:59:59Z"
	}

	if date == "" {
		date = GetCurrentTime()
	}

	d, err := time.Parse(consts.DATETIME_NANO_FORMAT, date)
	d1, err1 := time.Parse(consts.DATETIME_NANO_FORMAT, startDate)
	d2, err2 := time.Parse(consts.DATETIME_NANO_FORMAT, endDate)

	if err != nil || err1 != nil || err2 != nil {
		return false, errors.New("invalid date format")
	}

	return d.After(d1) && d.Before(d2), nil
}

// Count days between startDate and endDate
func CountDays(startDate string, endDate string) (int, error) {
	d1, err1 := time.Parse(consts.DATE_FORMAT, startDate)
	d2, err2 := time.Parse(consts.DATE_FORMAT, endDate)
	if err1 != nil || err2 != nil {
		return 0, errors.New("invalid date format")
	}
	return int((d2.Sub(d1).Hours() / 24) + 1), nil
}

func AddDays(date string, days int) (string, error) {
	d, err := time.Parse(consts.DATE_FORMAT, date)
	if err != nil {
		return "", errors.New("invalid date format")
	}
	return d.AddDate(0, 0, days).Format(consts.DATE_FORMAT), nil
}

func SubDays(date string, days int) (string, error) {
	d, err := time.Parse(consts.DATE_FORMAT, date)
	if err != nil {
		return "", errors.New("invalid date format")
	}
	return d.AddDate(0, 0, -days).Format(consts.DATE_FORMAT), nil
}

func CheckDateFormat(date string) bool {
	return regexp.MustCompile(consts.DATE_FORMAT_REGEX_PATTERN).MatchString(date)
}
