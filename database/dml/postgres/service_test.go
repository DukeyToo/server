// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package postgres

import (
	"reflect"
	"testing"
)

func TestPostgres_createServiceService(t *testing.T) {
	// setup types
	want := &Service{
		List: map[string]string{
			"all":   ListServices,
			"build": ListBuildServices,
		},
		Select: map[string]string{
			"build":        SelectBuildService,
			"count":        SelectBuildServicesCount,
			"count-images": SelectServiceImagesCount,
		},
		Delete: DeleteService,
	}

	// run test
	got := createServiceService()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("createServiceService is %v, want %v", got, want)
	}
}
