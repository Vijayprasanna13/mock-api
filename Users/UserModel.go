package Users

import (
	"database/sql"
	"fmt"
    "github.com/julienschmidt/httprouter"
	"net/http"
	"regexp"
    "os"
    "io"
)

type User struct {
    id          sql.NullInt64
    aadhar_id   sql.NullString
    name        sql.NullString
    dob         sql.NullString
    image_link  sql.NullString
    created_at  sql.NullString
    updated_at  sql.NullString
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
    
    file, handler, err := r.FormFile("image")
    if err != nil {
        fmt.Println(err)
        return ""
    }
    defer file.Close()

    if _, err := os.Stat("../images/"); os.IsNotExist(err) {
        os.Mkdir("../images/", 0775)
    }
    f, err := os.OpenFile("../images/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
    if err != nil {
        fmt.Println(err)
        return ""
    }
    defer f.Close()
    io.Copy(f, file)
    
    return "./test/"+handler.Filename
}

func convertUserRequestToUserObject(r *http.Request) User {

    var user User
    
    user.aadhar_id.String       = r.FormValue("aadhar_id")
    user.name.String            = r.FormValue("name")
    user.dob.String             = r.FormValue("dob")
    user.image_link.String      = storeImageAndGetFileName(r)

    return user
}

func validateAddUserRequest(r *http.Request) string {
    user := convertUserRequestToUserObject(r)

    if m, _ := regexp.MatchString("^[0-9]{12}$", user.aadhar_id.String); !m {
        return "Invalid aadhar number"
    }

    if m, _ := regexp.MatchString("^[a-zA-Z .]+$", user.name.String); !m {
        return "Invalid name"
    }

    if m, _ := regexp.MatchString("^[0-9]{4}-[0-9]{2}-[0-9]{2}$", user.dob.String); !m {
        return "Invalid dob"
    }

    return ""
}

func AddUser() httprouter.Handle {

    return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

        user_validation_result := validateAddUserRequest(r)

        if user_validation_result != "" {
            fmt.Println(user_validation_result)
        } else {
            fmt.Println("Upload successful...")
        }
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