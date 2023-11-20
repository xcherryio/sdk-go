package xc

type DBConverter interface {
	// ToDBValue converts global attribute value to database query value
	ToDBValue(val interface{}, hint *DBHint) (dbValue string, err error)
	// FromDBValue converts database query value to global attribute value
	FromDBValue(dbQueryValue string, hint *DBHint, resultPtr interface{}) error
}
