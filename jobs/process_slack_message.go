package jobs

import (
	"fmt"
)

func ProcessSlackMessage(queue string, args ...interface{}) error {
	fmt.Printf("From %s, %v\n", queue, args)
	return nil
}
