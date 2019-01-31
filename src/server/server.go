package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
	"util/trick"

	"github.com/go-redis/redis"
	"github.com/levigross/grequests"
)

// Server Struct
type Server struct {
	Port     int
	Appid    string
	Secret   string
	Host     string
	Password string
	Database int
	mux      sync.RWMutex
	client   *redis.Client
	token    Token
	timer    *time.Timer
}

// Token Struct
type Token struct {
	AccessToken string `json:"access_token"`
	Expires     int    `json:"expires_in"`
	Timestamp   int64  `json:"timestamp"`
}

// Response Struct
type Response struct {
	AccessToken string `json:"access_token"`
	Timestamp   int64  `json:"timestamp"`
}

// Start token server
func (s *Server) Start() {
	s.init()
	if s.Port == 0 {
		s.Port = 3000
	}
	if err := http.ListenAndServe(":"+strconv.Itoa(s.Port), s.handle()); err != nil {
		log.Fatalln("Server: ", err)
	}
}

func (s *Server) init() {
	s.client = redis.NewClient(&redis.Options{
		Addr:     s.Host,
		Password: s.Password,
		DB:       s.Database,
	})
	if err := s.client.Ping().Err(); err != nil {
		log.Fatalln("Redis: ", err)
	}
	token, err := s.load()
	if err != redis.Nil {
		log.Fatalln("Load: ", err)
	}
	s.token = token
	if !s.valid() {
		newToken := s.fetch()
		s.save(newToken)
		s.timer = time.NewTimer(time.Second * time.Duration(s.token.Expires-100))
	} else {
		s.timer = time.NewTimer(time.Second * time.Duration(s.token.Timestamp+int64(s.token.Expires)-time.Now().Unix()-100))
	}
	s.schedule()
}

func (s *Server) schedule() {
	go func() {
		for {
			<-s.timer.C
			newToken := s.fetch()
			s.save(newToken)
			s.timer.Reset(time.Second * time.Duration(s.token.Expires-100))
		}
	}()
}

func (s *Server) fetch() Token {
	log.Println("fetch")
	var params = map[string]string{
		"appid":      s.Appid,
		"secret":     s.Secret,
		"grant_type": "client_credential",
	}
	ro := &grequests.RequestOptions{
		Params: params,
	}
	res, _ := grequests.Get("https://api.weixin.qq.com/cgi-bin/token", ro)
	var t Token
	if err := json.Unmarshal(res.Bytes(), &t); err != nil {
		return Token{}
	}
	log.Printf("%s\n", res)
	return Token{AccessToken: t.AccessToken, Expires: t.Expires, Timestamp: time.Now().Unix()}
}

func (s *Server) valid() bool {
	if s.token.AccessToken == "" {
		return false
	}
	curTime := time.Now().Unix()
	if curTime >= s.token.Timestamp+int64(s.token.Expires)-100 {
		return false
	}
	return true
}

func (s *Server) load() (Token, error) {
	token, err := s.client.Get("access-token").Result()
	if err != nil {
		return Token{}, err
	}
	var t Token
	err = json.Unmarshal(trick.String2Bytes(token), &t)
	if err != nil {
		return Token{}, err
	}
	return t, nil
}

func (s *Server) save(token Token) {
	b, err := json.Marshal(token)
	if err != nil {
		log.Println(err)
	}
	err = s.client.Set("access-token", trick.Bytes2String(b), time.Second*time.Duration(s.token.Expires)).Err()
	if err != nil {
		log.Println(err)
	}
	s.mux.Lock()
	defer s.mux.Unlock()
	s.token = token
}

func (s *Server) handle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.mux.RLock()
		defer s.mux.RUnlock()
		w.Header().Set("Content-Type", "application/json")
		res := Response{
			AccessToken: s.token.AccessToken,
			Timestamp:   s.token.Timestamp + int64(s.token.Expires) - 100,
		}
		json.NewEncoder(w).Encode(res)
	}
}
