package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/girigiribauer/transcribe_goblueprints/thesaurus"
)

func main() {
	apiKey := os.Getenv("BHT_APIKEY")
	thesaurus := &thesaurus.BigHuge{
		APIKey: apiKey,
	}
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		word := s.Text()
		syns, err := thesaurus.Synonyms(word)
		if err != nil {
			log.Fatalf("%qの類語検索に失敗しました: %v\n", word, err)
		}
		if len(syns) == 0 {
			log.Fatalf("%qに類語はありませんでした\n", syns)
		}
		for _, syn := range syns {
			fmt.Println(syn)
		}
	}
}
