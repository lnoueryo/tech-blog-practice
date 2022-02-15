package models

import (
	"crypto/sha256"
	"encoding/gob"
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
	Image	  string
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
		Image: u.Image,
		CreatedAt: time.Now(),
	}
	filepath := fmt.Sprintf("./session/%v.txt", hashedSessionId)
    f, err := os.Create(filepath)
    if err != nil {
		return hashedSessionId, err
    }
	defer f.Close()
    enc := gob.NewEncoder(f)

    if err := enc.Encode(s); err != nil {
        return hashedSessionId, err
    }
	return hashedSessionId, nil
}

// Checks if the user is logged in and has a session, if not err is not nil
func CheckSession(r *http.Request) (Session, bool) {
	cookie, err := r.Cookie("_cookie")
	s := Session{}
	if err == nil {
		filepath := fmt.Sprintf("./session/%v.txt", cookie.Value)
		err = s.readSession(filepath)
		return s, IsSession(filepath)
	}
	return s, false
}

func GetSession(r *http.Request) (Session) {
	cookie, _ := r.Cookie("_cookie")
	s := Session{}
	filepath := fmt.Sprintf("./session/%v.txt", cookie.Value)
	s.readSession(filepath)
	return s
}

func DeliverSession(r *http.Request) (Session) {
	s := GetSession(r)
	err := s.GenerateCSRFToken(r); if err != nil {
		log.Print(err)
	}
	return s
}

func (s *Session)GenerateCSRFToken(r *http.Request) (error) {
	filepath := fmt.Sprintf("./session/%v.txt", s.Id)
	s.CSRFToken, _ = MakeRandomStr(32)
    f, err := os.Create(filepath)
    if err != nil {
        return err
    }
    enc := gob.NewEncoder(f)
    defer f.Close()

    if err := enc.Encode(&s); err != nil {
		return err
    }
	return nil
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

func (s *Session)readSession(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	dec := gob.NewDecoder(f)
	defer f.Close()
	if err := dec.Decode(&s); err != nil {
		return err
	}
	return nil
}
