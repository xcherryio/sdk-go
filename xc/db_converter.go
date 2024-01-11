package xc

type DBConverter interface {
	// ToDBValue converts global attribute value to database query value
	ToDBValue(val interface{}, hint *DatabaseHint) (dbValue string, err error)
	// FromDBValue converts database query value to global attribute value
	FromDBValue(dbQueryValue string, hint *DatabaseHint, resultPtr interface{}) error
}
