package main

import (
	"bufio"
	"flag"
	babble "github.com/danmarg/babble/lib"
	"log"
	"os"
	"strings"
)

var (
	corpus    string
	chain     string
	verbosity int
)

func init() {
	flag.StringVar(&corpus, "add_corpus", "", "Add words from this corpus.")
	flag.StringVar(&chain, "chain_file", "", "File to load/save chain from.")
}

func main() {
	flag.Parse()
	// Check args are specified.
	if chain == "" {
		log.Fatalln("missing required --chain")
	}
	// If chain file exists, try to load from it.
	var c babble.Chain
	loaded := false
	if f, err := os.Open(chain); !os.IsNotExist(err) {
		defer f.Close()
		if err != nil {
			log.Fatalln(err)
		}
		c, err = babble.ReadChain(f)
		if err != nil {
			log.Fatalln(err)
		}
		loaded = true
	}
	// Try to add from corpus.
	if corpus != "" {
		f, err := os.Open(corpus)
		defer f.Close()
		if err != nil {
			log.Fatalln(err)
		}
		scnr := bufio.NewScanner(f)
		for scnr.Scan() {
			i := strings.NewReader(scnr.Text())
			if loaded {
				err = c.AddCorpus(i)
				if err != nil {
					log.Fatalln(err)
				}
			} else {
				cn, err := babble.ReadCorpus(i)
				if err != nil {
					log.Fatalln(err)
				}
				c = *cn
			}
		}
		// Save chain.
		f, err = os.Create(chain)
		defer f.Close()
		if err != nil {
			log.Fatalln(err)
		}
		err = c.WriteChain(f)
		if err != nil {
			log.Fatalln(err)
		}
	}

	log.Printf("%v", c.Babble())
}
