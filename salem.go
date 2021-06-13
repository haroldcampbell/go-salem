// Copyright 2021 Harold Campbell. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package salem

// By default Mock is configured to generate 1 mock.
// This can be changed by using the factory.WithXXXItems(...) functions.
func Mock(t interface{}) *Factory {

	f := Factory{rootType: t}
	f.plan = NewPlan()
	f.WithExactItems(1) // Default to 1 item

	return &f
}
