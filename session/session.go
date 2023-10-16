package session

import (
	"crypto/rand"
	"net/http"
	"talker/config"
	"time"
)

var (
	SC *Control
)

type Control struct {
	cookieName     string
	SessionDataMap map[string]SessionData
}

type SessionData struct {
	ExpireAt time.Time
	Data     any
}

func Create(cookieName string) {
	SC = &Control{
		cookieName:     cookieName,
		SessionDataMap: make(map[string]SessionData),
	}
}

func (c *Control) Get(r *http.Request) (string, *SessionData, bool) {
	cookies := r.Cookies()
	if len(cookies) == 0 {
		return "", nil, false
	}

	cookie, err := r.Cookie(c.cookieName)
	if err != nil {
		return "", nil, false
	}

	s, ok := c.SessionDataMap[cookie.Value]
	if !ok {
		return "", nil, false
	}

	if s.ExpireAt.Before(time.Now()) {
		delete(c.SessionDataMap, cookie.Value)
		return "", nil, false
	}

	return cookie.Value, &s, true
}

func (c *Control) Delete(w http.ResponseWriter, id string) {
	delete(c.SessionDataMap, id)
	cookie := http.Cookie{
		Name:   c.cookieName,
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(w, &cookie)
}

func (c *Control) Save(w http.ResponseWriter, id string, sessionData *SessionData) {
	expireAt := time.Now().Add(time.Duration(config.CFG.MaxAgeSession) * time.Second)
	cookie := &http.Cookie{
		Path:     "/",
		Name:     c.cookieName,
		Value:    id,
		Expires:  expireAt,
		Secure:   config.CFG.SecureCookie,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}

	sessionData.ExpireAt = expireAt
	c.SessionDataMap[id] = *sessionData

	http.SetCookie(w, cookie)
}

func (c *Control) Create() (string, *SessionData) {
	sessionData := &SessionData{
		ExpireAt: time.Now().Add(3 * time.Hour),
	}

	return RandomID(), sessionData
}

func (c *Control) RemoveExpired() {
	for k, v := range c.SessionDataMap {
		if v.ExpireAt.Before(time.Now()) {
			delete(c.SessionDataMap, k)
		}
	}
}

func RandomID() string {
	const (
		length  = 16
		charset = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	)
	lenCharset := byte(len(charset))
	b := make([]byte, length)
	rand.Read(b)
	for i := 0; i < length; i++ {
		b[i] = charset[b[i]%lenCharset]
	}
	return string(b)
}
