package entities

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Track struct {
	Prev     *Track
	Next     *Track
	ID       int
	Duration int
}

type Playlist struct {
	Head      *Track
	Tail      *Track
	Paused    bool
	Played    bool
	NextTrack bool
	PrevTrack bool
	Stoped    bool
	cond      *sync.Cond
}

func NewPlaylist() *Playlist {
	return &Playlist{
		Head: nil,
		Tail: nil,
		cond: sync.NewCond(&sync.Mutex{}),
	}
}

func (l *Playlist) Add(value int, duration int) {
	newTrack := &Track{ID: value, Duration: duration}

	if l.Tail == nil {
		l.Head = newTrack
		l.Tail = newTrack
	} else {
		newTrack.Prev = l.Tail
		l.Tail.Next = newTrack
		l.Tail = newTrack
	}
}

func (l *Playlist) Play() {
	go func() {
		l.Played = true
		for node := l.Head; node != nil; {
			flag := false
			log.Printf("now play track:%d Duration:%d", node.ID, node.Duration)

			for i := 0; i < node.Duration; i++ {

				if l.NextTrack {
					l.NextTrack = false
					break
				}
				if l.Stoped {
					l.Stoped = false
					l.Played = false

					return
				}

				l.cond.L.Lock()
				for l.Paused {
					l.Played = false
					l.cond.Wait()
				}
				l.cond.L.Unlock()
				time.Sleep(time.Second)

				if l.PrevTrack {
					if node == l.Head {
						log.Println("It's a fist track in playlist")
						l.PrevTrack = false
						//flag = true
					} else {
						node = node.Prev
						l.PrevTrack = false
						flag = true
						break
					}

				}

			}
			if !flag {
				node = node.Next
			}

		}
		log.Println("playlist has ended")
		l.Played = false
	}()
}

func (l *Playlist) AddSong(msg []string) string {
	if len(msg) != 3 {
		log.Println("wrong input (wait: add [id] [duration])")
		return "wrong input (wait: add [id] [duration])"
	} else {
		id, err := strconv.Atoi(msg[1])
		if err != nil {
			log.Println("wrong input (wait: add [id] [duration])")
			return "wrong input (wait: add [id] [duration])"
		}
		duration, err := strconv.Atoi(msg[2])
		if err != nil {
			log.Println("wrong input (wait: add [id] [duration])")
			return "wrong input (wait: add [id] [duration])"
		}
		fmt.Println(id, msg[2])
		l.Add(id, duration)
		return "song was added"
	}
}

func (l *Playlist) Start(ch chan string, pauseCh chan bool) {

	input := <-ch
	msg := strings.Split(input, " ")
	switch msg[0] {
	case "play":
		if !l.Played {
			if l.Paused {
				l.Paused = false
				log.Println("playlist resumed")
				l.cond.Broadcast()
			} else {
				l.Play()
			}

		} else {
			log.Println("playlist alredy playing")

		}
	case "add":
		l.AddSong(msg)
	case "pause":
		if l.Played {
			if l.Paused {
				log.Println("playlist alredy paused")

			} else {
				l.Paused = true
				log.Println("playlist paused")

			}
		} else {
			log.Println("Playlist isn't playing")
		}
	case "next":
		if l.Played {
			l.NextTrack = true
			l.Paused = false
			l.cond.Broadcast()
			log.Println("Next track")
		} else {
			log.Println("Playlist isn't playing")
		}

	case "prev":
		if l.Played {
			l.PrevTrack = true
			l.Paused = false
			l.cond.Broadcast()
			log.Println("Prev track")
		} else {
			log.Println("Playlist isn't playing")
		}
	case "stop":
		if l.Played {
			l.Stoped = true
			log.Println("Playlist has stoped")
			l.Paused = false
			l.cond.Broadcast()
		} else {
			log.Println("Playlist isn't playing")
		}

	default:
		log.Println("unkown comand")

	}

}
