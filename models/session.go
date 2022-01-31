package models

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type Session struct {
	Id        int
	Uuid      string
	Name      string
	Email     string
	CreatedAt time.Time
	CSRFToken string
}

func CreateSession(user User) (cryptext string, err error) {
	s := Session {
		Id: user.Id,
		Name: user.Name,
		Email: user.Email,
		CreatedAt: time.Now(),
	}
	sessionId := string(user.Id) + timeToString(user.CreatedAt)
	hashedSessionId := sha256.Sum256([]byte(sessionId))
	cryptext = fmt.Sprintf("%x", hashedSessionId)
	filepath := fmt.Sprintf("./session/%v.txt", cryptext)
    f, err := os.Create(filepath)
    if err != nil {
        log.Fatal(err)
    }
    defer f.Close()
    enc := gob.NewEncoder(f)

    if err := enc.Encode(s); err != nil {
        log.Fatal(err)
    }
	return
}

// Checks if the user is logged in and has a session, if not err is not nil
func CheckSession(w http.ResponseWriter, r *http.Request) (Session, error) {
	cookie, err := r.Cookie("_cookie")
	session := Session{}
	if err == nil {
		filepath := fmt.Sprintf("./session/%v.txt", cookie.Value)
		isSession := IsSession(filepath)
		if isSession {
			session.readSession(filepath)
			session.AddCSRFToken(filepath)
			return session, err
		} else {
			err = errors.New("invalid session")
		}
	}
	return session, err
}

func (session *Session)AddCSRFToken(filepath string) {
	session.CSRFToken, _ = MakeRandomStr(32)
    f, err := os.Create(filepath)
    if err != nil {
        log.Fatal(err)
    }
    defer f.Close()
    enc := gob.NewEncoder(f)

    if err := enc.Encode(&session); err != nil {
        log.Fatal(err)
    }
	return
}

func Auth(w http.ResponseWriter, r *http.Request) (bool, error) {
	err := r.ParseForm()
	if err != nil {
		log.Print(err, "Cannot find user")
	}
	cookie, err := r.Cookie("_cookie")
	if err == nil {
		filename := fmt.Sprintf("./session/%v.txt", cookie.Value)
		isSession := IsSession(filename)
		session := Session{}
		if isSession {
			session.readSession(filename)
			if session.CSRFToken == r.Form.Get("csrf_token") {
				return true, nil
			}
			err = errors.New("invalid csrf_token")
		}
	}
	return false, err
}

// Checks if the user is logged in and has a session, if not err is not nil
func DeleteSession(w http.ResponseWriter, r *http.Request) (session Session, err error) {
	cookie, err := r.Cookie("_cookie")
	if err == nil {
		filename := fmt.Sprintf("./session/%v.txt", cookie.Value)
		isSession := IsSession(filename)
		if isSession {
			err := os.Remove(filename)
			if err != nil {
				log.Println(err)
			}
		}
	}
	cookie.MaxAge = -1
    http.SetCookie(w, cookie)
	return
}

func IsSession(filename string) bool {
    _, err := os.Stat(filename)
    if err == nil {
        return true
    } else {
        return os.IsExist(err)
    }
}

func (session *Session)readSession(filename string) {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	dec := gob.NewDecoder(f)
	if err := dec.Decode(&session); err != nil {
		log.Fatal("decode error:", err)
	}
}


func MakeRandomStr(n uint32) (string, error) {
    const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

    // 乱数を生成
    b := make([]byte, n)
    if _, err := rand.Read(b)
	err != nil {
        return "", err
    }

    // letters からランダムに取り出して文字列を生成
    var result string
    for _, v := range b {
        // index が letters の長さに収まるように調整
        result += string(letters[int(v)%len(letters)])
    }
    return result, nil
}

func timeToString(t time.Time) string {
    str := t.Format("20060102150405")
    return str
}

func Encrypt(plaintext string) (cryptext string) {
	cryptext = fmt.Sprintf("%x", sha256.Sum256([]byte(plaintext)))
	return
}