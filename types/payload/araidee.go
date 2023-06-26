// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    welcome, err := UnmarshalWelcome(bytes)
//    bytes, err = welcome.Marshal()

package payload

import "encoding/json"

func UnmarshalWelcome(data []byte) (Welcome, error) {
	var r Welcome
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Welcome) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type Welcome struct {
	GradeInfo GradeInfo `json:"grade_info"`
}

type GradeInfo struct {
	ID          interface{}   `json:"id"`
	StudentCode string        `json:"student_code"`
	Grades      []Grade       `json:"grades"`
	Status      string        `json:"status"`
	ReasonLocks []interface{} `json:"reason_locks"`
	LastUpdated interface{}   `json:"last_updated"`
}

type Grade struct {
	AcadYear         int64       `json:"acad_year"`
	SemesterID       int64       `json:"semester_id"`
	Gpax             interface{} `json:"gpax"`
	AnnounceDate     string      `json:"announce_date"`
	AnnounceViewable bool        `json:"announce_viewable"`
	Courses          []Course    `json:"courses"`
}

type Course struct {
	CourseID        int64       `json:"course_id"`
	CourseCode      string      `json:"course_code"`
	CourseNameTh    string      `json:"course_name_th"`
	CourseNameEn    string      `json:"course_name_en"`
	CourseGetCredit int64       `json:"course_get_credit"`
	CourseGrade     string      `json:"course_grade"`
	Modules         interface{} `json:"modules"`
}
