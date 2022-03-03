package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

const n = 8
const k = 4
const d = 2

//Node
type Node struct {
	id      int
	next    []*Node
	inChan  []chan *packet
	outChan []chan *packet
	packet  []*packet
}

//Packet with value
type packet struct {
	value int
	nodes []int
}

func run(arr []Node) {

	printer := make(chan string, 1)
	done := make(chan bool)
	done2 := make(chan bool)
	go printer2(printer, done2)

	go starter(&arr[0], printer)

	for i := 0; i < n; i++ {
		go NodeTransfer(&arr[i], printer)
	}
	go ender(&arr[n-1], done, printer)

	<-done
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
	for i := 0; i < k; i++ {
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(2000)))
		packetID := packet{i + 1, make([]int, 0)}
		printer <- "Wysylam Pakiet: " + strconv.Itoa(packetID.value) + "\n"
		zero.inChan[0] <- &packetID

	}
	close(zero.inChan[0])
}

//Odbieram wyswietlam i wysylam dalej
func NodeTransfer(currentNode *Node, printer chan string) {

	for packetID := range currentNode.inChan[0] {

		printer <- "Pakiet numer: " + strconv.Itoa(packetID.value) + " jest w Node: " + strconv.Itoa(currentNode.id) + " \n"
		packetID.nodes = append(packetID.nodes, currentNode.id)
		currentNode.packet = append(currentNode.packet, packetID)
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(2000)))

		nextRandom := rand.Intn(len(currentNode.outChan))
		currentNode.outChan[nextRandom] <- packetID
	}
	close(currentNode.outChan[0])
}

func ender(last *Node, done chan<- bool, printer chan string) {
	for packetID := range last.outChan[0] {
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(2000)))
		printer <- "Otrzymalem pakiet numer: " + strconv.Itoa(packetID.value) + "\n"
		//printer <- strconv.Itoa(packetID.value)
	}
	done <- true
}
func printer2(printer chan string, done2 chan<- bool) {
	for information := range printer {
		fmt.Print(information)
	}

}
func main() {
	//Tworz Tablice Wątkow / Node'ow
	arr := make([]Node, n)
	rand.NewSource(time.Now().UnixNano())
	//Dodaje Tutaj kanal do pierwszego Node aby otrzymywal pakiety
	arr[0].inChan = append(arr[0].inChan, make(chan *packet, 0))

	//Dodaje Tutaj nastepne Node' tzn 1->2->3 itp
	for j := 1; j < n; j++ {
		arr[j-1].id = j
		arr[j-1].next = append(arr[j-1].next, &arr[j])
	}

	//Dodaje Tutaj do kolejnych Node'ow aby mogly miedzy soba przesylac pakiety
	for x := 1; x < n-1; x++ {
		Channel := make(chan *packet, 0)
		arr[x-1].outChan = append(arr[x-1].outChan, Channel)
		arr[x].inChan = append(arr[x].inChan, Channel)
	}
	//Kalay do komunikacji z ostatnim Nodem
	Channel := make(chan *packet, 0)
	arr[n-2].outChan = append(arr[n-2].outChan, Channel)
	arr[n-1].inChan = append(arr[n-1].inChan, Channel)

	//Kanal wyjsciowy
	arr[n-1].outChan = append(arr[n-1].outChan, make(chan *packet, 0))
	arr[n-1].id = n

	//Losuje i lacze kanalami Node'y
	for i := 0; i < d; i++ {
		start := rand.Intn(n - 1)
		end := rand.Intn(n-start-1) + start + 1
		arr[start].next = append(arr[start].next, &arr[end])
		arr[start].outChan = append(arr[start].outChan, arr[end-1].outChan[0])
	}

	for i := 0; i < d; i++ {
		start := rand.Intn(n - 1)
		end := rand.Intn(n-start-1) + start + 1
		arr[end].next = append(arr[end].next, &arr[start])
		arr[start].outChan = append(arr[end].outChan, arr[start-1].outChan[0])
	}

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

	run(arr)
}
