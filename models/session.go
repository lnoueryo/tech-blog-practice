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
	Id        string
	UserId    int
	Name      string
	Email     string
	CreatedAt time.Time
	CSRFToken string
}

func CreateSession(u User) (cryptext string, err error) {
	sessionId := string(u.Id) + timeToString(u.CreatedAt)
	hashedByteSessionId := sha256.Sum256([]byte(sessionId))
	hashedSessionId := fmt.Sprintf("%x", hashedByteSessionId)
	s := Session {
		Id: hashedSessionId,
		UserId: u.Id,
		Name: u.Name,
		Email: u.Email,
		CreatedAt: time.Now(),
	}
	filepath := fmt.Sprintf("./session/%v.txt", hashedSessionId)
    f, err := os.Create(filepath)
	defer f.Close()
    if err != nil {
		return hashedSessionId, err
    }
    enc := gob.NewEncoder(f)

    if err := enc.Encode(s); err != nil {
        return hashedSessionId, err
    }
	return hashedSessionId, nil
}

// Checks if the user is logged in and has a session, if not err is not nil
func CheckSession(r *http.Request) (Session, error) {
	cookie, err := r.Cookie("_cookie")
	s := Session{}
	if err == nil {
		filepath := fmt.Sprintf("./session/%v.txt", cookie.Value)
		isSession := IsSession(filepath)
		if isSession {
			s.readSession(filepath)
			return s, err
		} else {
			err = errors.New("invalid session")
		}
	}
	return s, err
}

func GetSession(r *http.Request) (Session) {
	cookie, _ := r.Cookie("_cookie")
	s := Session{}
	filepath := fmt.Sprintf("./session/%v.txt", cookie.Value)
	s.readSession(filepath)
	return s
}

func GenerateCSRFToken(r *http.Request) (string) {
	s, err := CheckSession(r); if err != nil {
		return ""
	}
	filepath := fmt.Sprintf("./session/%v.txt", s.Id)
	s.CSRFToken, _ = MakeRandomStr(32)
    f, err := os.Create(filepath)
    if err != nil {
        return ""
    }
    defer f.Close()
    enc := gob.NewEncoder(f)

    if err := enc.Encode(&s); err != nil {
        err = errors.New("failed encode")
		return ""
    }
	return s.CSRFToken
}

func (s *Session)CheckCSRFToken(r *http.Request) bool {
	err := r.ParseForm()
	if err != nil {
		log.Print(err, "Cannot find user")
	}
	if s.CSRFToken != r.Form.Get("csrf_token") {
		return false
	}
	return true
}

// Checks if the user is logged in and has a session, if not err is not nil
func (s *Session)DeleteSession(w http.ResponseWriter, r *http.Request) error {
	cookie, err := r.Cookie("_cookie")
	if err != nil {
		return err
	}
	filename := fmt.Sprintf("./session/%v.txt", cookie.Value)
	isSession := IsSession(filename)
	if isSession {
		err := os.Remove(filename)
		if err != nil {
			return err
		}
	}
	DeleteCookie(w, r)
	return nil
}

func DeleteCookie(w http.ResponseWriter, r *http.Request) error {
	cookie, err := r.Cookie("_cookie")
	if err != nil {
		return err
	}
	cookie.MaxAge = -1
    http.SetCookie(w, cookie)
	return nil
}

func IsSession(filename string) bool {
    _, err := os.Stat(filename)
    if err == nil {
        return true
    } else {
        return os.IsExist(err)
    }
}

func (s *Session)readSession(filename string) {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer f.Close()
	dec := gob.NewDecoder(f)
	if err := dec.Decode(&s); err != nil {
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

func Encrypt(plaintext string) string {
	cryptext := fmt.Sprintf("%x", sha256.Sum256([]byte(plaintext)))
	return cryptext
}