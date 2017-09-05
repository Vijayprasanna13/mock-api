package Users

import (
	"database/sql"
	_ "fmt"
	"log"
	"mock-api/Helpers"
	"mock-api/databases"
	"net/http"
	"regexp"
	"time"
)

/*
*
*Function autheticates the user using the provided credentials
*@param adhaar id, dob, bio-metric content (fingerprints)
*@return bool
 */
func Authenticate() httprouter.Handle {

	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

		response := Helpers.ConvertToJSON("500 Internal Server Error", map[string]interface{}{
			"message": "Hold on. Something's wrong",
		})
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, response)
		return
	}
}
