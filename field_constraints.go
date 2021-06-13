// Copyright 2021 Harold Campbell. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
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
