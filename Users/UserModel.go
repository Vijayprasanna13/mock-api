package Users

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io"
	_ "log"
	_ "mock-api/Helpers"
	"mock-api/databases"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"
)

type User struct {
	id         sql.NullInt64
	aadhaar_id sql.NullString
	name       sql.NullString
	dob        sql.NullString
	image_link sql.NullString
	created_at sql.NullString
	updated_at sql.NullString
}

/*
*
*Function autheticates the user using the provided credentials
*@param adhaar id, dob, bio-metric content (fingerprints)
*@return bool
 */
func Authenticate() httprouter.Handle {

	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

		// response := Helpers.ConvertToJSON("500 Internal Server Error", map[string]interface{}{
		// 	"message": "Hold on. Something's wrong",
		// })
		// w.WriteHeader(http.StatusInternalServerError)
		// fmt.Fprintf(w, response)
		// return
	}
}

func storeImageAndGetFileName(r *http.Request) string {
	r.ParseMultipartForm(32 << 20)

	// Open the file and store the details in the handler
	file, handler, err := r.FormFile("image")
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer file.Close()

	// Create a folder called images in the src directory if not already exists
	if _, err := os.Stat("../images/"); os.IsNotExist(err) {
		os.Mkdir("../images/", 0775)
	}

	filePath := "../images/" + strconv.FormatInt(time.Now().UnixNano(), 10) + filepath.Ext(handler.Filename)

	// Store the uploaded image with the timestamp as its name in order to not replace multiple images with name filename
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer f.Close()
	io.Copy(f, file)

	return filePath
}

func convertUserRequestToUserObject(r *http.Request) User {

	var user User

	user.aadhaar_id.String = r.FormValue("aadhaar_id")
	user.name.String = r.FormValue("name")
	user.dob.String = r.FormValue("dob")
	user.image_link.String = storeImageAndGetFileName(r)

	return user
}

func validateAddUserRequest(r *http.Request) (User, error) {
	user := convertUserRequestToUserObject(r)

	if m, _ := regexp.MatchString("^[0-9]{12}$", user.aadhaar_id.String); !m {
		return User{}, errors.New("Invalid aadhaar number " + user.aadhaar_id.String)
	}

	if m, _ := regexp.MatchString("^[a-zA-Z .]+$", user.name.String); !m {
		return User{}, errors.New("Invalid name")
	}

	if m, _ := regexp.MatchString("^[0-9]{4}-[0-9]{2}-[0-9]{2}$", user.dob.String); !m {
		return User{}, errors.New("Invalid dob")
	}

	return user, nil
}

func storeUserDetails(user User) (string, error) {

	_, err := databases.DB_CONN.Exec(`INSERT INTO users
                                        (
                                            aadhaar_id,
                                            name,
                                            dob,
                                            image_link,
                                            created_at,
                                            updated_at
                                        )
                                        VALUES (?, ?, ?, ?, ?, ?)
                                     `, user.aadhaar_id.String, user.name.String, user.dob.String, user.image_link.String, time.Now().Format("2006/01/02 15:04:05"), time.Now().Format("2006/01/02 15:04:05"))
	if err != nil {
		return "", err
	}

	return "user created", nil
}

func AddUser() httprouter.Handle {

	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		user, err := validateAddUserRequest(r)

		if err != nil {
			fmt.Println(err)
			return
		}

		result, err := storeUserDetails(user)

		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(result)
		return

		// if user_validation_result != "" {
		//     response := Helpers.ConvertToJSON("500 Internal Server Error", map[string]interface{}{
		//         "message": user_validation_result,
		//     })
		//     w.WriteHeader(http.StatusInternalServerError)
		//     fmt.Fprintf(w, response)
		//     return
		// }
	}
}
