package babble

import (
	"bufio"
	"encoding/gob"
	"io"
	"log"
	"math/rand"
	"regexp"
)

const prefixLen int = 2

// Chains are made up of Tokens.
type Token interface{}

// StartString tokens are indications that the string starts here.
type StartString struct{}

// EndString tokens are indications that the string ends here.
type EndString struct{}

// Word is a single word.
type Word string

// Prefix is a n-gram of prefixLen length.
type Prefix [prefixLen]Token

// A Chain represents the mapping of prefix to Suffix or EndString Tokens.
type Chain struct {
	Ignore *regexp.Regexp
	Links  map[Prefix]map[Token]int
}

func init() {
	// Gob types.
	gob.Register(Chain{})
	gob.Register(EndString{})
	gob.Register(Prefix{})
	gob.Register(StartString{})
	gob.Register(Word(""))
}

// ReadChain reads a chain from the reader.
func ReadChain(reader io.Reader) (Chain, error) {
	var c Chain
	de := gob.NewDecoder(reader)
	err := de.Decode(&c)
	return c, err
}

// WriteChain saves the chain to the writer.
func (c *Chain) WriteChain(writer io.Writer) error {
	en := gob.NewEncoder(writer)
	return en.Encode(c)
}

// ReadCorpus reads text from reader and creates a new Chain from it.
func ReadCorpus(reader io.Reader) (*Chain, error) {
	c := new(Chain)
	err := c.AddCorpus(reader)
	return c, err
}

func startPrefix() Prefix {
	var p Prefix
	for i, _ := range p {
		p[i] = StartString{}
	}
	return p
}

// AddCorpus adds words from the reader to the Chain.
func (c *Chain) AddCorpus(reader io.Reader) error {
	if c.Links == nil {
		c.Links = make(map[Prefix]map[Token]int)
	}
	// Read tokens from reader.
	scnr := bufio.NewScanner(reader)
	scnr.Split(bufio.ScanWords)
	prev := startPrefix()
	for scnr.Scan() {
		cur := scnr.Text()
		if c.Ignore != nil && c.Ignore.MatchString(cur) {
			continue
		}
		var sfx Word = Word(cur)
		v, ok := c.Links[prev]
		if !ok {
			v = make(map[Token]int)
		}
		v[sfx]++
		c.Links[prev] = v
		copy(prev[:], prev[1:])
		prev[prefixLen-1] = Word(cur)
	}
	v, ok := c.Links[prev]
	if !ok {
		v = make(map[Token]int)
	}
	v[EndString{}]++
	c.Links[prev] = v
	return scnr.Err()
}

// Generate some babble.
func (c *Chain) Babble() string {
	pfx := startPrefix()
	var toks []Token
	for {
		sfxs, ok := c.Links[pfx]
		if !ok {
			log.Printf("Bailing on prefix %v", pfx)
			break
		}
		cnt := 0
		for _, x := range sfxs {
			cnt += x
		}
		r := rand.Intn(cnt)
		var sfx Token
		i := 0
		for c, x := range sfxs {
			i += x
			if i >= r {
				// This is our suffix.
				sfx = c
			}
		}
		toks = append(toks, sfx)
		// Construct a new prefix.
		if (sfx == EndString{}) {
			break
		}
		copy(pfx[:], pfx[1:])
		pfx[prefixLen-1] = sfx
	}
	r := ""
	for _, t := range toks {
		switch t := t.(type) {
		case Word:
			r += " " + string(t)
		}
	}
	return r
}
