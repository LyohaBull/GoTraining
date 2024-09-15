package main

import (
	"crypto/tls"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/albrow/forms"
	"github.com/opencoff/go-srp"
)

var db *sql.DB

func no_auth(w http.ResponseWriter, r *http.Request) {
	tmp, err := template.ParseFiles("./templates/FlexStart/no_auth.html", "./templates/FlexStart/base.html")
	if err != nil {
		panic(err)
	}
	tmp.ExecuteTemplate(w, "no_auth", nil)
}

func handleMain(w http.ResponseWriter, r *http.Request) {
	//	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./templates/FlexStart/assets/"))))
	tmp, err := template.ParseFiles("./templates/FlexStart/index.html", "./templates/FlexStart/base.html")
	if err != nil {
		panic(err)
	}
	tmp.ExecuteTemplate(w, "index", nil)
}
func srpAuth() func(http.ResponseWriter, *http.Request) {
	srp_servs := map[string]*srp.Server{}
	return func(w http.ResponseWriter, r *http.Request) {
		type SRP_Message struct {
			Status  string
			Message string
		}
		req, err := io.ReadAll(r.Body)
		if err != nil {
			errWrite(w, "На сервер поступили некорректные данные")
			return
		}

		s := SRP_Message{}
		err = json.Unmarshal(req, &s)
		if err != nil {
			errWrite(w, "На сервер поступили некорректные данные")
			return
		}
		if s.Status == "sendA" {
			srp_serv, response, err := loginUser_checkUserCred(db, s.Message)
			if err != nil {
				errWrite(w, "Проверка логина и пароля завершилась неудачно")
				return
			}
			srp_servs[strings.Split(s.Message, ":")[0]] = srp_serv
			s_r := SRP_Message{"getM1", response}
			json_data, err := json.Marshal(s_r)
			if err != nil {
				errWrite(w, "Не удалось отправить данные с сервера")
				return
			}
			w.Write(json_data)
			return
		}
		if s.Status == "sendM1" {
			proof, session_key, err := loginUser_genSrpKey(srp_servs[strings.Split(s.Message, ":")[0]], strings.Split(s.Message, ":")[1])
			if err != nil {
				errWrite(w, "Неверные логин или пароль")
				return
			}
			s_r := SRP_Message{"sendM2", proof}
			json_data, err := json.Marshal(s_r)
			if err != nil {
				errWrite(w, "Не удалось отправить данные с сервера")
				return
			}
			w.Write(json_data)
			fmt.Printf("%x\n", session_key)
			return
		}
	}
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	tmp, err := template.ParseFiles("./templates/FlexStart/login.html", "./templates/FlexStart/base.html")
	if err != nil {
		panic(err)
	}
	tmp.ExecuteTemplate(w, "login", nil)
}
func handleRegAuth(w http.ResponseWriter, r *http.Request) {
	formdata, err := forms.Parse(r)
	if err != nil {
		errWrite(w, "На сервер поступили некорректные данные")
		return
	}
	name := formdata.Get("name")
	surname := formdata.Get("surname")
	phone := formdata.Get("phone")
	name_test, _ := regexp.MatchString("^\\p{Cyrillic}{2,15}$", name)
	surname_test, _ := regexp.MatchString("^\\p{Cyrillic}{2,15}$", surname)
	if !name_test || !surname_test {
		errWrite(w, "На сервер поступили некорректные данные")
		return
	}
	png, err := otpGenerate(db, user_init{name, surname, phone})
	if err != nil {
		errWrite(w, "Ошибка генерации qr-code")
		return
	}
	w.Header().Set("Content-Type", "image/png")
	w.Write(png)
}
func errWrite(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(message))
}
func handleRegAuthUser(w http.ResponseWriter, r *http.Request) {
	req, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	type RegUser struct {
		User string
		Id   string
		Data string
	}
	u := RegUser{}
	err = json.Unmarshal(req, &u)
	if err != nil {
		errWrite(w, "На сервер поступили некорректные данные")
		delTmpKey(db, u.User)
		return
	}
	phone_test, _ := regexp.MatchString("^8\\d{10}$", u.Id)
	if !phone_test {
		errWrite(w, "На сервер поступили некорректные данные")
		delTmpKey(db, u.User)
		return
	}
	k := getTmpKey(db, u.User)
	decoded, err := hex.DecodeString(u.Data)
	if err != nil {
		errWrite(w, "Ошибка при декодировании данных из 16-го формата")
		delTmpKey(db, u.User)
		return
	}
	s, err := AesDeCrypt(decoded, k)
	if err != nil {
		errWrite(w, "Был введен неправильный код!")
		delTmpKey(db, u.User)
		return
	}
	fmt.Println(u.Id, u.User, s)
	err = createUser(db, u.Id, u.User, s)
	if err != nil {
		errWrite(w, "Ошибка при создании пользователя")
	}
	w.Write([]byte("UserCreated"))
}
func handleReg(w http.ResponseWriter, r *http.Request) {
	//	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./templates/FlexStart/assets/"))))
	tmp, err := template.ParseFiles("./templates/FlexStart/registrations.html", "./templates/FlexStart/base.html")
	if err != nil {
		panic(err)
	}
	tmp.ExecuteTemplate(w, "reg", nil)
}
func handleRegPhone(w http.ResponseWriter, r *http.Request) {
	req, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	uniq := checkUniquePhone(db, string(req))
	if uniq {
		w.Write([]byte("ok"))
	} else {
		w.Write([]byte("created"))
	}
}
func handleBlog(w http.ResponseWriter, r *http.Request) {
	//http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./templates/FlexStart/assets/"))))
	tmp, err := template.ParseFiles("./templates/FlexStart/blog.html")
	if err != nil {
		panic(err)
	}
	tmp.Execute(w, nil)
}
func readTLS() *http.Server {
	certPem, err := os.ReadFile("server-cert.pem")
	if err != nil {
		log.Fatal(err)
	}
	keyPem, err := os.ReadFile("server-key.pem")
	if err != nil {
		log.Fatal(err)
	}
	cert, err := tls.X509KeyPair(certPem, keyPem)
	if err != nil {
		log.Fatal(err)
	}
	cfg := &tls.Config{Certificates: []tls.Certificate{cert}}
	srv := &http.Server{
		TLSConfig:    cfg,
		ReadTimeout:  time.Minute,
		WriteTimeout: time.Minute,
	}
	return srv
}
func main() {
	db = dbConnect()
	srv := readTLS()
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./templates/FlexStart/assets/"))))
	http.HandleFunc("/", handleMain)
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/login/auth", srpAuth())
	http.HandleFunc("/registration/auth", handleRegAuth)
	http.HandleFunc("/registration/auth/user", handleRegAuthUser)
	http.HandleFunc("/registration", handleReg)
	http.HandleFunc("/registration/phone", handleRegPhone)
	http.HandleFunc("/blog", handleBlog)
	log.Fatal(srv.ListenAndServeTLS("", ""))
	db.Close()
}
