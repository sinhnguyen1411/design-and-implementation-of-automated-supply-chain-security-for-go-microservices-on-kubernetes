package health

import (
	"fmt"
	"github.com/google/uuid"
)

func Status() string {
	return fmt.Sprintf("ok:%s", uuid.NewString())
}