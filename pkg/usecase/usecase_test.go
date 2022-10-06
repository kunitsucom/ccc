// nolint: testpackage
package usecase

import (
	"testing"
)

func TestNew(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		if actual := New(WithRepository(nil), WithDomain(nil), WithInfra(nil)); actual == nil {
			t.Errorf("actual == nil")
		}
	})
}
