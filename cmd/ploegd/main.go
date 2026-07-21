// ploegd is Ploeg's single daemon: webhook ingest, provider SPI host,
// lease manager, outcome ingestion. Pre-alpha placeholder.
package main

import "fmt"

var version = "0.0.0-dev"

func main() {
	fmt.Printf("ploegd %s — pre-alpha, see docs/design.md\n", version)
}
