package functions

import (
	"kmuttBot/types/payload"
	"kmuttBot/types/response"
	"kmuttBot/utils/config"
	"net/http"
)

func GetGrade() (*payload.Welcome, *response.ErrorInstance) {

	var gradeInfo *payload.Welcome

	if err := DoRequest(nil, "GET", "https://api.kmutt.ac.th/kss/grade/v2/grades?academic_semester=2/2565", nil, func(r *http.Request) {
		r.Header.Set("Authorization", "Bearer "+config.C.AuthKey)
		r.Header.Set("Language", "EN")
	}, &gradeInfo); err != nil {
		return nil, err
	}

	return gradeInfo, nil
}
