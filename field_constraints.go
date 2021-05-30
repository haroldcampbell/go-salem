package salem

// SuggestedConstraintRetryAttempts is the default number of times to try generating a new mock before failing
const SuggestedConstraintRetryAttempts = 40

type FieldConstraint interface {
	IsValid(field interface{}) bool
}

type stringFieldConstraint struct {
	min int
	max int
}

func (s *stringFieldConstraint) IsValid(field interface{}) bool {
	str := field.(string)

	return len(str) >= s.min && len(str) <= s.max
}

func ConstrainStringLength(min int, max int) FieldConstraint {
	return &stringFieldConstraint{min: min, max: max}
}
