package main

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

const UPDATE_INTERVAL = time.Hour

type Message struct {
	Message  string `json:"message"`
	NextId   int    `json:"id,omitempty"`
	NextCode string `json:"code,omitempty"`
}

func index(rw http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(rw, "%v %s", g.Cs[0].Id, g.Cs[0].Code)
}

func ask(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json;charset=utf-8")
	req.ParseForm()
	id := req.Form.Get("id")
	code := req.Form.Get("code")
	// stuid := req.Form.Get("stuid")
	var m Message

	// match npc
	idint, _ := strconv.Atoi(id)
	curIndex, npc, err := g.GetNPC(idint)
	if err != nil {
		m = Message{
			Message: "You failed to find a NPC.",
		}
		bytes, _ := json.Marshal(m)
		fmt.Fprint(rw, string(bytes))
		return
	}
	fmt.Println(npc.Code, code, npc.Code == code)
	if npc.Code == code {
		if npc.IsLast {
			m = Message{
				Message: "Congratulations!",
			}
		} else {
			m = Message{
				Message:  "You'd better ask Miss." + strconv.Itoa(g.Cs[curIndex+1].Id) + ".",
				NextId:   g.Cs[curIndex+1].Id,
				NextCode: g.Cs[curIndex+1].Code,
			}
		}
	} else {
		m = Message{
			Message: "Sorry, the code you gave me is wrong, I can't tell you more.",
		}
	}
	bytes, _ := json.Marshal(m)
	fmt.Fprint(rw, string(bytes))
}

type NPC struct {
	Id int
	// LastId  int
	NextId  int
	Code    string
	IsFirst bool
	IsLast  bool
}

type Game struct {
	Cs []NPC
}

func (g *Game) Init() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	dup := map[int]bool{}
	for i := 0; i < 32+r.Intn(256); i++ {
		rnum := r.Intn(900) + 100
		if !dup[rnum] {
			h := sha1.New()
			io.WriteString(h, string(rnum)+"lalalalademacia")
			g.Cs = append(g.Cs, NPC{
				Id:   rnum,
				Code: base64.StdEncoding.EncodeToString([]byte(base64.StdEncoding.EncodeToString(h.Sum(nil)))),
			})
		}
		dup[rnum] = true
	}
	for i := 1; i < len(g.Cs)-1; i++ {
		g.Cs[i].NextId = g.Cs[i+1].Id
	}
	if len(g.Cs) > 1 {
		g.Cs[0].NextId = g.Cs[1].Id
	}
	if len(g.Cs) > 0 {
		g.Cs[0].IsFirst = true
		g.Cs[len(g.Cs)-1].IsLast = true
	} else {
		panic("No npc found.")
	}
}

func (g *Game) GetNPC(id int) (int, NPC, error) {
	for i, npc := range g.Cs {
		if id == npc.Id {
			return i, npc, nil
		}
	}
	return 0, NPC{}, errors.New("NPC not found")
}

var g Game

func Update() {
	time.AfterFunc(UPDATE_INTERVAL, func() {
		newGame := Game{}
		newGame.Init()
		g = newGame
		for i, npc := range g.Cs {
			fmt.Println(i, npc.Id, npc.NextId, npc.Code, npc.IsFirst, npc.IsLast)
		}
		Update()
	})
}

func main() {
	g.Init()
	for i, npc := range g.Cs {
		fmt.Println(i, npc.Id, npc.NextId, npc.Code, npc.IsFirst, npc.IsLast)
	}
	Update()
	http.HandleFunc("/", index)
	http.HandleFunc("/ask/", ask)
	log.Fatal(http.ListenAndServe(":3000", nil))
}
