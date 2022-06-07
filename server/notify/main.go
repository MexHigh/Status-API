package notify

import (
	"fmt"
	"status-api/structs"
)

var servicesThatAreDown = make(map[string]*structs.CheckResultWithNameAndTime, 0)

func ReportUp(result *structs.CheckResultWithNameAndTime) {
	fmt.Println("reported up")
}

func ReportDown(result *structs.CheckResultWithNameAndTime) {
	fmt.Println("reported down")
}
