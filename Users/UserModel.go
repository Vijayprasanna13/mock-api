package Users

import (
    "database/sql"
    "fmt"
    "github.com/julienschmidt/httprouter"
    "net/http"
    "regexp"
    "os"
    "io"
    "time"
    "strconv"
    "path/filepath"
    "errors"
    "mock-api/databases"
    "mime/multipart"
    "image/jpeg"
    "image/png"
    "github.com/jteeuwen/imghash"
    "mock-api/Helpers"
    "path"
    "image"
)

type User struct {
    id          sql.NullInt64
    aadhaar_id  sql.NullString
    name        sql.NullString
    dob         sql.NullString
    image_link  sql.NullString
    created_at  sql.NullString
    updated_at  sql.NullString
}

func validateAuthenticateRequest(r *http.Request) error {

    if m, _ := regexp.MatchString("^[0-9]{12}$", r.FormValue("aadhaar_id")); !m {
        return errors.New("Invalid aadhaar number " + r.FormValue("aadhaar_id"))
    }

    if m, _ := regexp.MatchString("^[0-9]{4}-[0-9]{2}-[0-9]{2}$", r.FormValue("dob")); !m {
        return errors.New("Invalid dob")
    }

    return nil
}

func getUserByAadhaarIdAndDob(aadhaar_id string, dob string) (User, error) {
    var user User
    var name, image_link string

    err := databases.DB_CONN.QueryRow(`SELECT name, image_link
                                         FROM users
                                         WHERE
                                            aadhaar_id=? AND
                                            dob=?
                                        `, aadhaar_id, dob).Scan(&name, &image_link)
    if err != nil {
        return User{}, err
    }

    user.aadhaar_id.String  = aadhaar_id
    user.name.String        = name
    user.dob.String         = dob
    user.image_link.String  = image_link

    return user, nil
}

func getHammingDistance(r *http.Request) (uint64, error) {

    user, err := getUserByAadhaarIdAndDob(r.FormValue("aadhaar_id"), r.FormValue("dob"))

    if err != nil {
        return 0, err
    }

    r.ParseMultipartForm(32 << 20)
    inputImage, inputImageHandler, err := r.FormFile("image")
    if err != nil {
        return 0, err
    }

    inputImageHash, err := getAverageHashOfImageFile(inputImage, inputImageHandler.Filename)

    userImageFile, err := os.Open(user.image_link.String)
    if err != nil {
        return 0, err
    }

    var userImage image.Image

    if path.Ext(user.image_link.String) == ".jpg" || path.Ext(user.image_link.String) == ".jpeg" {
        userImage, err = jpeg.Decode(userImageFile)
    } else if path.Ext(user.image_link.String) == ".png" {
        userImage, err = png.Decode(userImageFile)
    } else {
        return 0, errors.New("Unsupported image format. Supported formats : png, jpg, jpeg")
    }
    if err != nil {
        return 0, err
    }

    userImageHash := imghash.Average(userImage)

    hammingDistance := imghash.Distance(inputImageHash, userImageHash)

    return hammingDistance, nil
}

/*
*
*Function autheticates the user using the provided credentials
*@param adhaar id, dob, bio-metric content (fingerprints)
*@return bool
 */
func Authenticate() httprouter.Handle {

    return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

        err := validateAuthenticateRequest(r)
        if err != nil {
            response := Helpers.ConvertToJSON("500 Internal Server Error", map[string]interface{}{
                "message": err.Error(),
                "auth_status": false,
            })
            w.WriteHeader(http.StatusInternalServerError)
            fmt.Fprintf(w, response)
            return
        }

        hammingDistance, err := getHammingDistance(r)

        if err != nil {
            response := Helpers.ConvertToJSON("500 Internal Server Error", map[string]interface{}{
                "message": err.Error(),
                "auth_status": false,
            })
            w.WriteHeader(http.StatusInternalServerError)
            fmt.Fprintf(w, response)
            return
        }

        if hammingDistance > 5 {
            response := Helpers.ConvertToJSON("200 Successful", map[string]interface{}{
                "message": "unsuccessful",
                "auth_status": false,
            })

            w.WriteHeader(http.StatusInternalServerError)
            fmt.Fprintf(w, response)
            return
        }
        
        response := Helpers.ConvertToJSON("200 Successful", map[string]interface{}{
            "message": "successful",
            "auth_status": true,
        })

        w.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintf(w, response)
        return
    }
}

func getAverageHashOfImageFile(file multipart.File, filePath string) (uint64, error) {
    file.Seek(0, 0)

    var img image.Image
    var err error
    if path.Ext(filePath) == ".jpg" || path.Ext(filePath) == ".jpeg" {
        img, err = jpeg.Decode(file)
    } else if path.Ext(filePath) == ".png" {
        img, err = png.Decode(file)
    } else {
        return 0, errors.New("Unsupported image format. Supported formats : png, jpg, jpeg")
    }

    if err != nil {
        return 0, err
    }
    avg := imghash.Average(img)

    return avg, nil
}

func storeImageAndGetFileName(r *http.Request) (string, error) {
    r.ParseMultipartForm(32 << 20)
    
    // Open the file and store the details in the handler
    file, handler, err := r.FormFile("image")   // file is of type multipart.File
    if err != nil {
        return "", err
    }

    // Create a folder called images in the src directory if not already exists
    if _, err := os.Stat("images/"); os.IsNotExist(err) {
        os.Mkdir("images/", 0775)
    }

    filePath := "images/" + strconv.FormatInt(time.Now().UnixNano(), 10) + filepath.Ext(handler.Filename)
    
    // Store the uploaded image with the timestamp as its name in order to not replace multiple images with name filename
    f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
    if err != nil {
        return "", err
    }
    defer f.Close()
    io.Copy(f, file)

    if err != nil {
        return "", err
    }

    defer file.Close()

    return filePath, nil
}

func convertUserRequestToUserObject(r *http.Request) (User, error) {

    var user User
    
    user.aadhaar_id.String      = r.FormValue("aadhaar_id")
    user.name.String            = r.FormValue("name")
    user.dob.String             = r.FormValue("dob")
    filePath, err              := storeImageAndGetFileName(r)
    if err != nil {
        return User{}, err
    }

    user.image_link.String      = filePath

    return user, nil
}

func validateAddUserRequest(r *http.Request) (User, error) {
    user, err := convertUserRequestToUserObject(r)

    if err != nil {
        return User{}, err
    }

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
            response := Helpers.ConvertToJSON("500 Internal Server Error", map[string]interface{}{
                "message": err.Error(),
            })
            w.WriteHeader(http.StatusInternalServerError)
            fmt.Fprintf(w, response)
            return
        }

        result, err := storeUserDetails(user)
        if err != nil {
            response := Helpers.ConvertToJSON("500 Internal Server Error", map[string]interface{}{
                "message": err.Error(),
            })
            w.WriteHeader(http.StatusInternalServerError)
            fmt.Fprintf(w, response)
            return
        }

        response := Helpers.ConvertToJSON("200 Successful", map[string]interface{}{
            "message": result,
        })
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintf(w, response)
        return
    }
}
