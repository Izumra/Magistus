package storage

import "errors"

var (
	ErrUserNotFound       = errors.New("user with provided data was't found")
	ErrChartNotFound      = errors.New("chart with provided data was't found")
	ErrChartsNotFound     = errors.New("no charts for this profile")
	ErrUsersNotFound      = errors.New("no one of users of the system was't found")
	ErrUserChartsNotFound = errors.New("list of the user charts is empty")
	ErrUserExists         = errors.New("user with this id already exists")
)
