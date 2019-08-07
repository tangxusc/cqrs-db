package util

import (
	"fmt"
	"testing"
)

func TestGenerateUuid(t *testing.T) {
	uuid := GenerateUuid()
	fmt.Println(uuid)
}
