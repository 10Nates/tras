/*
This is based on the speech generation of old TRAS, available at https://github.com/10Nates/tras-old/blob/main/speechgen.js.
Note, I made tras-old before GPT2, and I was completely unaware of LLMS, so that was basically my uneducatedversion of it.

The basic idea is to algorithmically generate sentences in the same theme as Noam Chomsky's "Colorless green ideas sleep furiously."
That is, a context-free grammar generator using transformational grammer. In the future, my plan is to give it meaning, or at least
better structure, and perhaps even make my own neural network that isn't quite as heavy as the GPTs.

After a lot of back and forth with ChatGPT including a lot of acronyms for subjects that didn't exist, some code surgery,
and a few original ideas, I made this "mockup" in JS (it's fully functional but JS is a toy language)

```js
const WordPOS = require('wordpos');
const wordpos = new WordPOS();

// Function to get a random element from an array

	function getRandomElement(array) {
	    return array[Math.floor(Math.random() * array.length)];
	}

// Function to get a random noun

	async function randomNoun() {
	    const nouns = await wordpos.randNoun();
	    return getRandomElement(nouns);
	}

// Function to get a random verb

	async function randomVerb() {
	    const verbs = await wordpos.randVerb();
	    return getRandomElement(verbs);
	}

// Function to get a random adverb

	async function randomAdverb() {
	    const adverbs = await wordpos.randAdverb();
	    return getRandomElement(adverbs);
	}

// Function to get a random adjective

	async function randomAdjective() {
	    const adjectives = await wordpos.randAdjective();
	    return getRandomElement(adjectives);
	}

	const grammar = {
	    S: ['NP VP'],
	    NP: ['Det N', 'Det AP N', 'Det N PP'],
	    VP: ['V', 'V NP', 'V NP PP', 'V Adv', 'V Adv NP'],
	    PP: ['P NP'],
	    AP: ['A', 'A AP'],
	    Det: [...],
	    N: [ randomNoun ],
	    V: [ randomVerb ],
	    P: [...],
	    Adv: [ randomAdverb ],
	    A: [ randomAdjective ],
	};

// Function to expand a non-terminal symbol

	async function expand(symbol) {
	    if (grammar[symbol]) {
	        const expansions = grammar[symbol];
	        let chosenExpansion = expansions[Math.floor(Math.random() * expansions.length)];
	        if (typeof chosenExpansion === "function") {
	            chosenExpansion = await chosenExpansion();
	        }
	        const symbols = chosenExpansion.split(' ');
	        let asyncJoin = "";
	        for (let i = 0; i < symbols.length; i++) {
	            asyncJoin += (i > 0 ? " " : "") + await expand(symbols[i])
	        }
	        return asyncJoin;
	    }
	    return symbol;
	}

// Generate a sentence

	async function generateSentence() {
	    return await expand('S');
	}

// Generate and log a sentence
generateSentence().then(console.log)
```

I then moved to Bing Chat in order to get more detail on grammar types that were not yet included and I was not yet aware off, then I
had it convert ITS Python into Go. This is is what I got from
that:

```go
// A context-free grammar generator using phase structure rules
// The grammar can generate sentences with the following parts:
// - Subject (S)
// - Verb (V)
// - Object (O)
// - Complement (C)
// - Adverbial (A)
// - Determiner (D)
// - Adjective (Adj)
// - Preposition (P)
// - Noun (N)

// Import the packages for natural language processing and random number generation
package main

import (

	"fmt"
	"math/rand"
	"time"

)

// Define the grammar rules as a map of slices

	var grammar_rules = map[string][]string{
		"S":   {"NP VP", "NP VP A"},
		"NP":  {"D N", "D Adj N", "D N PP"},
		"VP":  {"V", "V NP", "V NP PP", "V NP C"},
		"PP":  {"P NP"},
		"C":   {"Adj", "PP"},
		"A":   {"Adv", "PP"},
		"D":   {"the", "a", "an"},
		"N":   {"cat", "dog", "book", "table", "chair", "man", "woman", "boy", "girl"},
		"V":   {"is", "has", "reads", "writes", "likes", "hates", "sees", "gives"},
		"Adj": {"big", "small", "red", "blue", "happy", "sad", "smart", "stupid"},
		"P":   {"on", "under", "in", "with", "to", "from", "for"},
		"Adv": {"quickly", "slowly", "loudly", "quietly", "well", "badly"},
	}

// Define a function to generate a random sentence from the grammar

	func generate_sentence(symbol string) string {
		// If the symbol is a terminal, return it as it is
		if _, ok := GrammarRules[symbol]; !ok {
			return symbol
		}
		// If the symbol is a non-terminal, choose a random production rule and apply it recursively
		rand.Seed(time.Now().UnixNano())
		rule := GrammarRules[symbol][rand.Intn(len(GrammarRules[symbol]))]
		sentence := ""
		for _, part := range split(rule, " ") {
			sentence += generate_sentence(part) + " "
		}
		return sentence
	}

// Define a helper function to split a string by a separator

	func split(str string, sep string) []string {
		result := []string{}
		start := 0
		for i := 0; i < len(str); i++ {
			if str[i:i+1] == sep {
				result = append(result, str[start:i])
				start = i + 1
			}
		}
		result = append(result, str[start:])
		return result
	}

// Print 10 random sentences generated by the grammar

	func main() {
		for i := 0; i < 10; i++ {
			fmt.Println(generate_sentence("S"))
		}
	}

```

This document is result of this great wonderful ride of Generative Pre-trained Transformers. In order to make a worse
version of what they do. Technology is truly magical. I used Bing Chat's Go code as a base and imposed my significant alteration
of ChatGPT code on top of it (ChatGPT also originally had one of those small lists. They had very similar bases. Probably
could have done this with just one of the models. Too late now).
*/
package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// Define the grammar rules as a map of slices
var GrammarRules = map[string][]string{
	// grammar tree
	"S":  {"NP VP", "NP VP A"},
	"NP": {"D N", "D Adj N", "D N PP"},
	"VP": {"V", "V NP", "V NP PP", "V NP C"},
	"PP": {"P NP"},
	"C":  {"Adj", "PP"},
	"A":  {"Adv", "PP"},
	// word lists
	"D": {"the", "a", "an"},
	"P": {"on", "under", "in", "with", "to", "from", "for"},
}

// Define a function to generate a random sentence from the grammar
func generateSentenceFromSymbol(symbol string) string {
	// If the symbol is a terminal, return it as it is
	if _, ok := GrammarRules[symbol]; !ok {
		return symbol
	}
	// If the symbol is a non-terminal, choose a random production rule and apply it recursively
	rand.Seed(time.Now().UnixNano())
	rule := GrammarRules[symbol][rand.Intn(len(GrammarRules[symbol]))]
	sentence := ""
	for _, part := range strings.Split(rule, " ") {
		switch part {
		// TODO: add custom wordpos implementation to pull random word from dictionary.
		case "N":
		case "V":
		case "Adj":
		case "Adv":
		default:
			sentence += generateSentenceFromSymbol(part) + " "
		}
	}
	return strings.TrimSpace(sentence)
}

// Print 10 random sentences generated by the grammar
func TestSpeechGeneration() {
	for i := 0; i < 10; i++ {
		fmt.Println(generateSentenceFromSymbol("S"))
	}
}
