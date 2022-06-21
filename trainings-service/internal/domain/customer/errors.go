package customer

import "errors"

var (
	ErrEmptyWorkoutDayUUID = errors.New("empty customer workout day UUID")
	ErrEmptyGroupDate      = errors.New("empty workout date")
	ErrEmptyCustomerUUID   = errors.New("empty customer UUID")
	ErrEmptyCustomerName   = errors.New("empty customer name")
	ErrEmptyGroupUUID      = errors.New("empty trainer workout group UUID")
)
