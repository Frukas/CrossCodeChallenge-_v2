package dataretriver

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func DataretriverClassTester() {
	fmt.Println("Dataretriver funcionado")
}

type NumberSet struct {
	Numbers []float32 `json:"numbers"`
}

var responseChannel = make(chan []byte)
var responseTotalPages = make(chan int)

func isValidPage(rs []byte) bool {
	if string(rs) != `{"numbers":[]}` {
		return true
	} else {
		return false
	}
}

func PageDataRetrive() ([]float32, int) {

	var responseNumberSet []float32

	i := 1

	go pagesCalls()

	pages := <-responseTotalPages

	for {
		select {
		case numberArray := <-responseChannel:
			if i >= pages {
				fmt.Println(numberArray)
				return responseNumberSet, pages
			}

			i++
			responseNumberSet = append(responseNumberSet, resultByteTofloatSlice(numberArray)...)
		}
	}
}

func getNumberFromPages(in int) []byte {

	pagesSTR := fmt.Sprintf("http://challenge.dienekes.com.br/api/numbers?page=%d", in)

	resp, err := http.Get(pagesSTR)
	if err != nil {
		log.Println("Err when get the page index: ", in, err)
		return getNumberFromPages(in)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln("Erro no corpo no indice: ", in, err)
	}

	defer resp.Body.Close()

	return body

}

func resultByteTofloatSlice(res []byte) []float32 {
	var numberRs NumberSet

	json.Unmarshal(res, &numberRs)

	return numberRs.Numbers

}

func findFistPart(first, last int) (int, int) {

	if isValidPage(getNumberFromPages(first)) == isValidPage(getNumberFromPages(last)) {
		return findFistPart(last, last*2)
	} else {
		return first, last
	}
}

func findIndex(first, last int) int {

	mid := (first + last) / 2

	if last == first+1 {
		return first
	} else if isValidPage(getNumberFromPages(first)) != isValidPage(getNumberFromPages(mid)) {
		return findIndex(first, mid)
	} else {
		return findIndex(mid, last)
	}
}

func findElement(first, last int) int {
	fst, lst := findFistPart(first, last)

	return findIndex(fst, lst)
}

func pagesCalls() {
	totalPagesNumber := findElement(1, 2)
	fmt.Println(totalPagesNumber)

	for i := 1; i <= totalPagesNumber; i++ {
		time.Sleep(10 * time.Millisecond)

		go channelSender(i)
	}

	responseTotalPages <- totalPagesNumber

}

func channelSender(i int) {
	responseChannel <- getNumberFromPages(i)
}
