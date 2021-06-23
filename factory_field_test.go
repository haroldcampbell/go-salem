// Copyright 2021 Harold Campbell. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package salem_test

import (
	"go-salem"
	"testing"

	"github.com/stretchr/testify/assert"
)

type person struct {
	Name string
}

func Test_OnField(t *testing.T) {
	names := []string{"Mary", "lilo", "Frank"}
	fieldHandler := func(itemIndex int) interface{} {
		return names[itemIndex]
	}

	f := salem.Mock(person{}).
		OnField("Name", fieldHandler).
		WithExactItems(3)

	results := f.ExecuteToType().([]person)

	assert.Equal(t, 3, len(results)) // Sanity check
	assert.Equal(t, names[0], results[0].Name)
	assert.Equal(t, names[1], results[1].Name)
	assert.Equal(t, names[2], results[2].Name)

}
