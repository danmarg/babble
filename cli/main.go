package main

import (
	"flag"
	"github.com/cenkalti/log"
	babble "github.com/danmarg/babble/lib"
	"os"
)

var (
	corpus    string
	chain     string
	verbosity int
)

func init() {
	flag.StringVar(&corpus, "add_corpus", "", "Add words from this corpus.")
	flag.StringVar(&chain, "chain_file", "", "File to load/save chain from.")
	flag.IntVar(&verbosity, "v", 2, "Verbosity.")
}

func main() {
	flag.Parse()
	// Check args are specified.
	if chain == "" {
		log.Fatalln("missing required --chain")
	}
	log.SetLevel([]log.Level{log.DEBUG, log.INFO, log.NOTICE, log.WARNING, log.ERROR, log.CRITICAL}[verbosity])
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
		if loaded {
			err = c.AddCorpus(f)
			if err != nil {
				log.Fatalln(err)
			}
		} else {
			cn, err := babble.ReadCorpus(f)
			if err != nil {
				log.Fatalln(err)
			}
			c = *cn
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

	log.Noticef("XXX: %v", c.Babble())
}
