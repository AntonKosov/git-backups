package internal_test

import (
	"testing"

	"github.com/AntonKosov/git-backups/internal"
	"github.com/AntonKosov/git-backups/internal/internalfakes"
)

func TestSum(t *testing.T) {
	fakeTransformer := internalfakes.FakeTransformer{}
	fakeTransformer.TransformStub = func(v int) int { return v * 10 }

	result := internal.Sum(2, 3, &fakeTransformer)
	if result != 50 {
		t.Errorf("Sum(2, 3, x10) = %v", result)
	}
}
