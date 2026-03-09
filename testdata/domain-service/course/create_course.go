package course

import "net/http"

// @sequence post
// @model Course.Create
// @param Name request
// @result course Course

// @sequence response json
// @var course
func CreateCourse(w http.ResponseWriter, r *http.Request) {}
