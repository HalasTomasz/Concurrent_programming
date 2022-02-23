package main

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

const n = 7
const d = 2

var number = 0

type Node struct {
	id      int
	next    []*Node
	reciver chan []*PacketInfo
	sender  chan []*PacketInfo
	info    []*PacketInfo
}

type PacketInfo struct {
	id     int
	next   int
	cost   int
	change int
	owner  int
}

var mutex = &sync.Mutex{}

func Sender(currentNode *Node, printer chan<- string) {

	for {
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(500)))
		Tmp := make([]*PacketInfo, 0)
		mutex.Lock()
		for x := 0; x < len(currentNode.info); x++ {
			if currentNode.info[x].change == 1 {
				currentNode.info[x].change = 0
				Tmp = append(Tmp, currentNode.info[x])
			}
		}
		mutex.Unlock()

		if len(Tmp) > 0 {
			for i := 0; i < len(currentNode.next); i++ {
				printer <- "Moje " + strconv.Itoa(currentNode.id) + " mowi do " + strconv.Itoa(currentNode.next[i].id) + "\n"
				currentNode.next[i].reciver <- Tmp
			}
		}

	}
}

func Receiver(currentNode *Node, printer chan<- string) {

	for {
		select {
		case packetID := <-currentNode.reciver:
			{

				sender := packetID[0].owner
				printer <- "Dostalem " + strconv.Itoa(currentNode.id) + " od " + strconv.Itoa(sender) + "\n"
				for x := 0; x < len(packetID); x++ {
					for y := 0; y < len(currentNode.info); y++ {
						if packetID[x].id == currentNode.info[y].id && packetID[x].cost+1 < currentNode.info[y].cost {
							mutex.Lock()
							printer <- "W nodzie " + strconv.Itoa(currentNode.id) + " aktualizuje " + strconv.Itoa(packetID[x].id) + "\n"
							currentNode.info[y].cost = packetID[x].cost + 1
							currentNode.info[y].next = sender
							currentNode.info[y].change = 1
							mutex.Unlock()
						}
					}
				}
			}
		}
	}
}

func printer2(printer chan string, done2 chan<- bool) {
	for information := range printer {
		fmt.Print(information)
	}
	done2 <- true
}
func main() {

	arr := make([]Node, n)

	for j := 0; j < n-1; j++ {
		arr[j].id = j
		arr[j].next = append(arr[j].next, &arr[j+1])
		arr[j+1].next = append(arr[j+1].next, &arr[j])
	}
	arr[n-1].id = n - 1
	arr[n-1].next = append(arr[n-1].next, &arr[0])
	arr[0].next = append(arr[0].next, &arr[n-1])

	for i := 0; i < n; i++ {
		arr[i].reciver = make(chan []*PacketInfo)
		arr[i].sender = make(chan []*PacketInfo)
	}

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	for i := 0; i < d; i++ {
		start := r1.Intn(n - 1)
		end := r1.Intn(n-start-1) + start + 1
		arr[start].next = append(arr[start].next, &arr[end])
		arr[end].next = append(arr[end].next, &arr[start])
	}

	for j := 0; j < n; j++ {
		fmt.Print(arr[j].id)
		i := 0
		for i, item := range arr[j].next {
			if i == 0 {
				fmt.Print(" -> ", item.id)
			} else {
				fmt.Print(" , ", item.id)
			}
		}
		i++

		fmt.Println()
	}

	for j := 0; j < n; j++ {
		for i := 0; i < n; i++ {
			packetID := PacketInfo{}
			packetID.id = i
			packetID.owner = j
			packetID.change = 1
			if j < i {
				packetID.next = j + 1
				packetID.cost = i - j
			} else if j > i {
				packetID.next = j - 1
				packetID.cost = j - i
			} else {
				packetID.next = j
				packetID.cost = 0
			}
			arr[j].info = append(arr[j].info, &packetID)
		}
	}

	for j := 0; j < n; j++ {
		for i := 0; i < len(arr[j].next); i++ {
			if math.Abs(float64(j-arr[j].next[i].id)) > 1 {
				arr[j].info[arr[j].next[i].id].next = arr[j].next[i].id
				arr[j].info[arr[j].next[i].id].cost = i
			}
		}
	}

	for x := 0; x < n; x++ {
		for j := 0; j < len(arr[x].info); j++ {
			fmt.Println(arr[x].info[j])
		}
		fmt.Println()
	}
	printer := make(chan string, 10)
	done2 := make(chan bool)

	var wg sync.WaitGroup

	go printer2(printer, done2)

	for x := 0; x < n; x++ {
		go Sender(&arr[x], printer)
		go Receiver(&arr[x], printer)
	}
	wg.Wait()
	//<-done2
	time.Sleep(time.Millisecond * time.Duration(7000))
	for x := 0; x < n; x++ {
		for j := 0; j < len(arr[x].info); j++ {
			fmt.Println(arr[x].info[j])
		}
		fmt.Println()
	}
}
