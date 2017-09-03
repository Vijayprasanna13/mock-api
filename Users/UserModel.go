package Users

import (
	"database/sql"
	_ "fmt"
	"log"
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
func IsRollValid(rollno string) bool {
	/*
	*
	 */
}
