package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

const n = 8

const k = 4

var number = 0

const d = 2
const maxLifespan = 8

//Node
type Node struct {
	id      int
	next    []*Node
	inChan  []chan *packet
	outChan []chan *packet
	outFic  []chan *packet
	packet  []*packet
}

//Packet with value
type packet struct {
	value    int
	nodes    []int
	size     int
	response chan bool
	lifespan int
}

func run(arr []Node, trapsChannels []chan bool) {

	printer := make(chan string)
	logger := make(chan bool)

	var wg sync.WaitGroup

	go printer2(printer)

	go starter(&arr[0], printer)

	//go hunter(n, trapsChannels, printer)

	for i := 0; i < n; i++ {

		go NodeTransfer(&arr[i], printer, trapsChannels[i], logger)
	}
	wg.Add(1)
	go ender(arr, &wg, printer, logger)

	wg.Wait()

	packets(arr, printer)
	nodes(arr, printer)

}

func nodes(arr []Node, printer chan string) {

	for j := 0; j < n; j++ {
		printer <- "Wezel " + strconv.Itoa(arr[j].id) + ": "
		for y := 0; y < len(arr[j].packet); y++ {
			printer <- strconv.Itoa(arr[j].packet[y].value)
			printer <- " "
		}
		printer <- "\n"
	}
}
func packets(arr []Node, printer chan string) {

	for y := 0; y < len(arr[0].packet); y++ {
		printer <- "Pakiet " + strconv.Itoa(arr[0].packet[y].value) + ": "
		for x := 0; x < len(arr[0].packet[y].nodes); x++ {
			printer <- strconv.Itoa(arr[0].packet[y].nodes[x]) + " "
		}
		printer <- "\n"
	}
	printer <- "\n"

}

//Wysylam Pakiety
func starter(zero *Node, printer chan string) {

	var i = 0
	for i < k {

		packetID := packet{i + 1, make([]int, 0), 0, make(chan bool), 0}

		zero.inChan[0] <- &packetID
		printer <- "Wysylam Pakiet: " + strconv.Itoa(packetID.value) + "\n"
		//oczekiwanie na odpowiedz czy pakiet dotarl
		for {
			select {
			case <-packetID.response:
				{
					i++
				}

			default:
				{
					time.Sleep(time.Millisecond * time.Duration(rand.Intn(2000)))
					continue
				}
			}
			break
		}

	}

}

//Odbieram wyswietlam i wysylam dalej
func NodeTransfer(currentNode *Node, printer chan string, hunterChannel <-chan bool, logger chan bool) {

	trap := 0

	for {

		select {
		case receivedPacket := <-currentNode.inChan[0]:
			{

				printer <- "Pakiet dotarl"
				if trap == 1 {
					printer <- "Pakiet " + strconv.Itoa(receivedPacket.value) + " wpadl w pulapke w wierzcholku " + strconv.Itoa(currentNode.id) + "\n"
					trap = 0
					logger <- true
					continue
				}
				printer <- "Pakiet " + strconv.Itoa(receivedPacket.value) + " jest w wierzcholku " + strconv.Itoa(currentNode.id) + "\n"
				receivedPacket.nodes = append(receivedPacket.nodes, currentNode.id)
				currentNode.packet = append(currentNode.packet, receivedPacket)
				receivedPacket.size = receivedPacket.size + 1
				receivedPacket.lifespan++

				time.Sleep(time.Millisecond * time.Duration(rand.Intn(2000)))
				//pakiet umiera jesli przekroczyl czas zycia
				receivedPacket.response <- true
				if receivedPacket.lifespan > maxLifespan {
					printer <- "Pakiet " + strconv.Itoa(receivedPacket.value) + " zostal zniszczony" + "\n"
					logger <- true
					continue
				}

				for {
					nextChannel := rand.Intn(len(currentNode.outChan))
					currentNode.outChan[nextChannel] <- receivedPacket

					select {

					case <-receivedPacket.response:
						break

					default:
						time.Sleep(time.Millisecond * time.Duration(rand.Intn(2000)))
					}

				}

			}
		case <-hunterChannel:
			{
				trap = 1
			}
		default:
			{
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(2000)))
				continue
			}
		}

	}

}

func ender(arr []Node, wg *sync.WaitGroup, printer chan string, logger <-chan bool) {
	defer wg.Done()
	var counter = 0
	messages := make([]*packet, 0)

	for {
		select {
		case <-logger:
			// count every death as ‘received’ for simplicity’s sake
			counter++
		case receivedPacket := <-arr[n-1].inChan[0]:
			counter++
			printer <- "Pakiet " + strconv.Itoa(receivedPacket.value) + " zostal odebrany" + "\n"
			messages = append(messages, receivedPacket)
			// count the messages
		default:
			// finish after all the messages have arrived
			if counter == k {
				// close the logger
				close(arr[n-1].inChan[0])
			}
		}
	}
}

//watek klusownika
func hunter(sizeOfNodes int, channelMessage []chan bool, printer chan string) { //dostaje liste channeli, zrob liste channeli trap dla wierzcholkow
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	r := rand.New(rand.NewSource(r1.Int63()))

	for {
		a := r1.Intn(sizeOfNodes-2) + 1
		printer <- "Zastawiam pulapke w wierzcholku " + strconv.Itoa(a+1) + "\n"

		channelMessage[a] <- true
		time.Sleep(time.Duration(r.Float64() * float64(time.Second) * 9))
	}

}

func printer2(printer chan string) {
	for information := range printer {
		fmt.Print(information)
	}
}
func main() {
	//Tworz Tablice Wątkow / Node'ow
	arr := make([]Node, n)
	rand.Seed(time.Now().UTC().UnixNano())
	//Dodaje Tutaj kanal do pierwszego Node aby otrzymywal pakiety
	arr[0].inChan = append(arr[0].inChan, make(chan *packet, 0))
	arr[0].outFic = append(arr[0].outFic, make(chan *packet, 0))
	//Dodaje Tutaj nastepne Node' tzn 1->2->3 itp
	for j := 1; j < n; j++ {
		arr[j-1].id = j
		arr[j-1].next = append(arr[j-1].next, &arr[j])

	}

	//Dodaje Tutaj do kolejnych Node'ow aby mogly miedzy soba przesylac pakiety
	for x := 1; x < n-1; x++ {
		Channel := make(chan *packet)
		arr[x-1].outChan = append(arr[x-1].outChan, Channel)
		arr[x].inChan = append(arr[x].inChan, Channel)
		arr[x].outFic = append(arr[x].outFic, make(chan *packet))
	}
	//Kalay do komunikacji z ostatnim Nodem
	Channel := make(chan *packet, 0)
	arr[n-2].outChan = append(arr[n-2].outChan, Channel)
	arr[n-1].inChan = append(arr[n-1].inChan, Channel)
	arr[n-2].outFic = append(arr[n-2].outFic, make(chan *packet))

	//Kanal wyjsciowy
	arr[n-1].outChan = append(arr[n-1].outChan, make(chan *packet))
	arr[n-1].outFic = append(arr[n-1].outFic, make(chan *packet))
	arr[n-1].id = n

	//Losuje i lacze kanalami Node'y
	for i := 0; i < d; i++ {
		start := rand.Intn(n - 1)
		end := rand.Intn(n-start-1) + start + 1
		arr[start].next = append(arr[start].next, &arr[end])
		arr[start].outChan = append(arr[start].outChan, arr[end-1].outChan[0])
	}

	//for i := 0; i < 6; i++ {
	//	start := rand.Intn(n - 1)
	//	end := rand.Intn(n-3) + 1
	//	arr[start].next = append(arr[start].next, &arr[end])
	///	arr[start].outChan = append(arr[start].outChan, arr[end-1].inChan[0])
	//}

	//for i := 0; i < d; i++ {
	//	starter := rand.Intn(n - 1)
	//	ender := rand.Intn(n-starter-1) + starter + 1
	//	arr[starter].next = append(arr[starter].next, &arr[ender])
	//	arr[starter].outChan = append(arr[starter].outChan, arr[ender-1].outChan[0])
	//	}

	//Wypisuje Graf
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

	trapsChannels := make([]chan bool, 0)
	for i := 0; i < n; i++ {
		trapsChannels = append(trapsChannels, make(chan bool))
	}

	run(arr, trapsChannels)
}
