package main

import (
	"bytes"
	"crypto"
	"crypto/sha256"
	"database/sql"
	"encoding/base32"
	"encoding/base64"
	"errors"
	"fmt"
	"image/png"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"github.com/opencoff/go-srp"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "alex"
	password = "1c2bfg"
	dbname   = "alexdb"
)

type user_init struct {
	Name    string
	Surname string
	Phone   string
}

func dbConnect() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	return db
}

/*
	func printUsers(db *sql.DB) {
		rows, err := db.Query("select * from users")
		if err != nil {
			panic(err)
		}
		defer rows.Close()
		users := []db_user{}
		for rows.Next() {
			u := db_user{}
			rows.Scan(&u.Id, &u.Name, &u.Password)
			users = append(users, u)
		}
		for _, u := range users {
			fmt.Println(u)
		}
	}

	func addUser(db *sql.DB, newUser db_user, passwd string) {
		h := sha256.New()
		h.Write([]byte(passwd))
		//fmt.Printf("insert into users (name,age,password) values ( '%s', %d, '\\x%x');\n", newUser.Name, newUser.Age, h.Sum(nil))

		r, err := db.Query(fmt.Sprintf("insert into users (name,password) values ( '%s', '\\x%x');\n", newUser.Name, h.Sum(nil)))

		if err != nil {
			panic(err)
		}

		defer r.Close()

		row := db.QueryRow("select id from users where name='" + newUser.Name + "';")
		user_id := 0
		row.Scan(&user_id)
		rr, err := db.Query(fmt.Sprintf("insert into tokens (user_id,token,date) values ( '%d', '\\x%x', '%s');\n", user_id, h.Sum(nil), time.Now().Format("01.02.2006")))
		if err != nil {
			panic(err)
		}
		defer rr.Close()
	}
*/
func getToken(db *sql.DB, username string) []byte {
	row := db.QueryRow("select token from tokens where user_id= (select id from users where name='" + username + "')")
	token := []byte{}
	row.Scan(&token)
	return token
}
func getUserfromToken(db *sql.DB, token string) (string, error) {
	dst := make([]byte, base64.StdEncoding.DecodedLen(len(token)))
	n, err := base64.StdEncoding.Decode(dst, []byte(token))
	if err != nil {
		fmt.Println("decode error:", err)
		return "", err
	}
	dst = dst[:n]
	row := db.QueryRow(fmt.Sprintf("select users.name from users join tokens on id=user_id where token='\\x%x';", dst))
	name := ""
	row.Scan(&name)
	if len(name) > 0 {
		return name, nil
	} else {
		return "", errors.New("no auth user")
	}
}
func loginUser_checkUserCred(db *sql.DB, user_cred string) (*srp.Server, string, error) {
	identifier := strings.Split(user_cred, ":")[0]
	fmt.Println(identifier)
	r := db.QueryRow(fmt.Sprintf("select verifier from users_verifier where identifier='%s';", identifier))
	var b []byte
	r.Scan(&b)
	_, B, err := srp.ServerBegin(user_cred)
	if err != nil {
		return nil, "", err
	}

	// Now, pretend to lookup the user db using "I" as the key and
	// fetch salt, verifier etc.
	s, v, err := srp.MakeSRPVerifier(string(b))
	if err != nil {
		return nil, "", err
	}
	srv, err := s.NewServer(v, B)
	if err != nil {
		return nil, "", err
	}
	serv_creds := srv.Credentials()
	return srv, serv_creds, nil
}
func loginUser_genSrpKey(srv *srp.Server, M1 string) (string, []byte, error) {
	proof, ok := srv.ClientOk(M1)

	if !ok {
		return "", nil, errors.New("client auth failde")
	}
	ks := srv.RawKey()
	return proof, ks, nil
}

func otpGenerate(db *sql.DB, user user_init) ([]byte, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "user",
		AccountName: user.Name + " " + user.Surname,
	})
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	img, err := key.Image(300, 300)
	if err != nil {
		return nil, err
	}
	png.Encode(&buf, img)
	h := sha256.New()
	fmt.Println(key.Secret())
	s := user.Name + " " + user.Surname + key.Secret()
	h.Write([]byte(s))
	fmt.Printf("%x\n", h.Sum(nil))
	//h.Write([]byte(GeneratePassCode(key.Secret())))
	r, err := db.Query(fmt.Sprintf("insert into tmp_users (fullname,data) values ( '%s', '\\x%x');\n", user.Name+" "+user.Surname, h.Sum(nil)))
	if err != nil {
		panic(err)
	}
	defer r.Close()
	return buf.Bytes(), nil
}

func GeneratePassCode(utf8string string) string {
	secret := base32.StdEncoding.EncodeToString([]byte(utf8string))
	passcode, err := totp.GenerateCodeCustom(secret, time.Now(), totp.ValidateOpts{
		Period:    30,
		Skew:      1,
		Digits:    otp.DigitsSix,
		Algorithm: otp.AlgorithmSHA512,
	})
	if err != nil {
		panic(err)
	}
	return passcode
}
func getTmpKey(db *sql.DB, id string) []byte {
	row := db.QueryRow("select data from tmp_users where fullname='" + id + "';")
	key := []byte{}
	row.Scan(&key)
	return key
}
func delTmpKey(db *sql.DB, id string) {
	db.QueryRow("delete from tmp_users where fullname='" + id + "';")
}

func createUser(db *sql.DB, id string, fullname string, password string) error {
	name, surname := strings.Split(fullname, " ")[0], strings.Split(fullname, " ")[1]
	s, _ := srp.NewWithHash(crypto.SHA256, 1024)
	fmt.Println(id)
	//s, _ := srp.New(1024)
	v, _ := s.Verifier([]byte(id), []byte(password))
	Id, verif := v.Encode()
	r, err := db.Query("insert into users (name,surname) values ('" + name + "','" + surname + "');")
	if err != nil {
		return err
	}
	r.Close()
	r, err = db.Query("insert into users_verifier values ((select id from users where name='" + name + "' and surname='" + surname + "'),'" + Id + "','" + verif + "');")
	if err != nil {
		return err
	}
	r.Close()
	r, err = db.Query("delete from tmp_users where fullname='" + fullname + "';")
	if err != nil {
		return err
	}
	r.Close()
	return nil
}
func checkUniquePhone(db *sql.DB, id string) bool {
	row := db.QueryRow("select count(identifier) from users_verifier where identifier='" + id + "';")
	count := []byte{}
	row.Scan(&count)
	return count[0] == 48
}

/*
	M1 := <-c2
	proof, ok := srv.ClientOk(M1)

	if !ok {
		panic("client auth failed")
	}
	c3 <- proof

	ks := srv.RawKey()

	fmt.Printf("Server Key: %x\n", ks)
	/*

		h := sha256.New()
		h.Write([]byte(passwd))
	 alert(client.t)	r := db.QueryRow(fmt.Sprintf("select password from users where name='%s';", u.Name))
		var b []byte
		r.Scan(&b)
		if len(b) == 0 {
			fmt.Println("auth invalid!")
			return false
		}
		if (fmt.Sprintf("%x", h.Sum(nil))) == (fmt.Sprintf("%x", b)) {
			return true
		} else {
			return false
		}*/

//func loginUser_genSrpKey(M1 string) (string, error) {
