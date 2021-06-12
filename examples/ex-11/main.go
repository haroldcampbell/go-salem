package main

import (
	"fmt"
	"go-salem"

	"github.com/haroldcampbell/go_utils/utils"
)

type item struct {
	SKU string
}

type basic struct {
	Children []item
}

// Sequences Example
// An example showing the difference between sequencing based on item-index vs sequence-index
// EnsureSequence: The sequence items are chosen based on the item index
// EnsureSequenceAcross: The sequence items are chosen based on the sequence index

func main() {
	seq := []string{"a", "b", "c", "d"}

	seqIndexTap := salem.Tap().
		EnsureSequenceAcross("Children.SKU", seq[0], seq[1], seq[2], seq[3]).
		WithExactItems(2)

	itemIndexTap := salem.Tap().
		EnsureSequence("Children.SKU", seq[0], seq[1], seq[2], seq[3]).
		WithExactItems(2)

	execute(seqIndexTap, "Sequence-index-based mocks")
	execute(itemIndexTap, "Items-index-based mocks")
}

func execute(tap *salem.Factory, msg string) {
	factory := salem.Mock(basic{})

	factory.Ensure("Children", tap). // Uses thethe item index
						WithExactItems(2)

	results := factory.Execute()

	str := utils.PrettyMongoString(results)
	fmt.Printf("%s:\n%v\n\n", msg, str)
}
