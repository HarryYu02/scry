package indexer

import (
	"strings"

	"github.com/kljensen/snowball"
)

func stem(word string) (string, error) {
	stemmed, err := snowball.Stem(word, "english", true)
	if err != nil {
		return "", err
	}
	upperStemmed := strings.ToUpper(stemmed)
	return upperStemmed, nil
}

func isStopWord(word string) bool {
	word = strings.ToUpper(strings.TrimSpace(word))
	switch word {
	case "A", "ABOUT", "ABOVE", "AFTER", "AGAIN", "AGAINST", "ALL", "AM", "AN",
		"AND", "ANY", "ARE", "AS", "AT", "BE", "BECAUSE", "BEEN", "BEFORE",
		"BEING", "BELOW", "BETWEEN", "BOTH", "BUT", "BY", "CAN", "DID", "DO",
		"DOES", "DOING", "DON", "DOWN", "DURING", "EACH", "FEW", "FOR", "FROM",
		"FURTHER", "HAD", "HAS", "HAVE", "HAVING", "HE", "HER", "HERE", "HERS",
		"HERSELF", "HIM", "HIMSELF", "HIS", "HOW", "I", "IF", "IN", "INTO", "IS",
		"IT", "ITS", "ITSELF", "JUST", "ME", "MORE", "MOST", "MY", "MYSELF",
		"NO", "NOR", "NOT", "NOW", "OF", "OFF", "ON", "ONCE", "ONLY", "OR",
		"OTHER", "OUR", "OURS", "OURSELVES", "OUT", "OVER", "OWN", "S", "SAME",
		"SHE", "SHOULD", "SO", "SOME", "SUCH", "T", "THAN", "THAT", "THE", "THEIR",
		"THEIRS", "THEM", "THEMSELVES", "THEN", "THERE", "THESE", "THEY",
		"THIS", "THOSE", "THROUGH", "TO", "TOO", "UNDER", "UNTIL", "UP",
		"VERY", "WAS", "WE", "WERE", "WHAT", "WHEN", "WHERE", "WHICH", "WHILE",
		"WHO", "WHOM", "WHY", "WILL", "WITH", "YOU", "YOUR", "YOURS", "YOURSELF",
		"YOURSELVES":
		return true
	}
	return false
}
