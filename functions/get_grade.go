package functions

import (
	"kmuttBot/types/payload"
	"kmuttBot/types/response"
)

func GetGrade() (*payload.Welcome, *response.ErrorInstance) {

	var gradeInfo *payload.Welcome
	if err := DoRequest(nil, "GET", "https://bot.ryyyyyy.com/kmutt-grade.json", nil, nil, &gradeInfo); err != nil {
		return nil, err
	}

	return gradeInfo, nil
}
