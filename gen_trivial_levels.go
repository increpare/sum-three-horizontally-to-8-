package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime/pprof"
	"time"
)

var breite int
var höhe int
var zahlmax int
var zielsumme int
var zielstückzahl int

var maxTiefe int = 20

func leeresBrett() ([][]uint8, []uint8) {

	b := make([][]uint8, breite)
	zeilen := make([]uint8, breite*höhe)
	for i := 0; i < breite; i++ {
		b[i] = zeilen[i*höhe : (i+1)*höhe]
	}

	return b, zeilen
}

func kopieBrett(zeilen []uint8, zeilen2 []uint8) {
	copy(zeilen2, zeilen)
}

func basteleSilhouette(brett [][]uint8) (stückzahl int) {

	for i := range brett {

		var gipfelhöhe = rand.Intn(höhe + 1)
		for j := 0; j < gipfelhöhe; j++ {
			brett[i][j] = 0
		}
		for j := gipfelhöhe; j < höhe; j++ {
			brett[i][j] = 1 + uint8(rand.Intn(zahlmax))
			stückzahl++
		}
		// for j := 0; j < höhe; j++ {
		// 	brett[i][j] = 1 + uint8(rand.Intn(zahlmax))
		// 	stückzahl++
		// }
	}

	randZeile := uint8(rand.Intn(breite))
	for j := 0; j < höhe; j++ {
		if brett[randZeile][j] != 0 {
			break
		}

		brett[randZeile][j] = 1 + uint8(rand.Intn(zahlmax))
		stückzahl++
	}

	if brett[0][höhe-1] == 0 && brett[1][höhe-1] == 0 {
		brett[1][höhe-1] = 1 + uint8(rand.Intn(zahlmax))
		stückzahl++
	}

	if brett[breite-2][höhe-1] == 0 && brett[breite-1][höhe-1] == 0 {
		brett[breite-2][höhe-1] = 1 + uint8(rand.Intn(zahlmax))
		stückzahl++
	}
	return
}

const charmap = ".123456"
const erstcharmap = "pqwerty"

func druckBrett(brett [][]uint8) {
	for i := 0; i < breite+2; i++ {
		fmt.Printf("#")
	}
	fmt.Printf("\n")

	for j := 0; j < höhe; j++ {
		fmt.Printf("#")
		for i := 0; i < breite; i++ {
			var val = brett[i][j]
			if i == 0 && j == 0 {
				fmt.Printf(string(erstcharmap[val]))
			} else {
				fmt.Printf(string(charmap[val]))
			}
		}
		fmt.Printf("#\n")
	}

	for i := 0; i < breite+2; i++ {
		fmt.Printf("#")
	}
	fmt.Printf("\n\n")
}

func stabil(brett [][]uint8, maske [][]uint8) bool {
	wasgefunden := false
	for i := 0; i < breite-2; i++ {
		for j := 0; j < höhe; j++ {
			var e1 = brett[i][j]
			var e2 = brett[i+1][j]
			var e3 = brett[i+2][j]
			if e1 == 0 || e2 == 0 || e3 == 0 {
				continue
			}
			if int(e1+e2+e3) != zielsumme {
				continue
			}

			maske[i][j] = 1
			maske[i+1][j] = 1
			maske[i+2][j] = 1

			wasgefunden = true
		}
	}

	// for i := 0; i < breite; i++ {
	// 	for j := 0; j < höhe-2; j++ {
	// 		var e1 = brett[i][j]
	// 		var e2 = brett[i][j+1]
	// 		var e3 = brett[i][j+2]
	// 		if e1 == 0 || e2 == 0 || e3 == 0 {
	// 			continue
	// 		}
	// 		if int(e1+e2+e3) != zielsumme {
	// 			continue
	// 		}

	// 		maske[i][j] = 1
	// 		maske[i][j+1] = 1
	// 		maske[i][j+2] = 1

	// 		wasgefunden = true
	// 	}
	// }

	return wasgefunden == false
}

func brettLeerung(brett [][]uint8, maske [][]uint8) int {
	removed := 0
	for i := 0; i < breite; i++ {
		for j := 0; j < höhe; j++ {
			if maske[i][j] == 1 {
				maske[i][j] = 0
				brett[i][j] = 0
				removed++
				//fallen lassen
				for k := j; k > 0; k-- {
					brett[i][k] = brett[i][k-1]
				}
				brett[i][0] = 0
			}
		}
	}
	return removed
}

func machTauschen(brett [][]uint8, brettZeilen []uint8, stückzahl int, tiefe uint) (erfolgreich int, gesamt int, interessante int, lösungstiefe uint, lösungsEingabe []int) {
	lösungstiefe = 1000

	maske := masken[tiefe]

	brett2 := bretter[tiefe]
	brett2Zeilen := bretterZeilen[tiefe]

	for i := 0; i < breite-1; i++ {
		for j := 0; j < höhe; j++ {

			w1 := brett[i][j]
			w2 := brett[i+1][j]

			//wenn leer oder identitsch, überspring!
			if w1 == w2 {
				continue
			}

			gesamt++

			kopieBrett(brettZeilen, brett2Zeilen)

			stückzahl2 := stückzahl

			//tauschen
			brett2[i][j], brett2[i+1][j] = brett2[i+1][j], brett2[i][j]

			if brett2[i][j] == 0 {
				for k := j; k > 0; k-- {
					brett2[i][k] = brett2[i][k-1]
				}
				brett2[i][0] = 0
			} else if brett2[i+1][j] == 0 {
				for k := j; k > 0; k-- {
					brett2[i+1][k] = brett2[i+1][k-1]
				}
				brett2[i+1][0] = 0
			}

			if j+1 < höhe {
				if brett2[i][j+1] == 0 {
					k := j + 1
					for k < höhe && brett2[i][k] == 0 {
						brett2[i][k] = brett2[i][k-1]
						brett2[i][k-1] = 0
						k++
					}
				}

				if brett2[i+1][j+1] == 0 {
					k := j + 1
					for k < höhe && brett2[i+1][k] == 0 {
						brett2[i+1][k] = brett2[i+1][k-1]
						brett2[i+1][k-1] = 0
						k++
					}
				}
			}

			istinteressant := false

			// vorstückzahl2 := stückzahl2

			for !stabil(brett2, maske) {
				removed := brettLeerung(brett2, maske)
				stückzahl2 -= removed
				istinteressant = true
			}

			// if istinteressant {
			// 	fmt.Printf("vor (%d):\n", vorstückzahl2)
			// 	druckBrett(brett)
			// 	fmt.Printf("(%d,%d)\n", i, j)

			// 	fmt.Printf("nach (%d):\n", stückzahl2)
			// 	druckBrett(brett2)
			// }

			if istinteressant {
				interessante++
			} else {
				continue
			}

			if stückzahl2 == 0 {

				erfolgreich++
				if lösungstiefe > tiefe {
					//gelöst
					lösungstiefe = tiefe
					lösungsEingabe = make([]int, tiefe+1)

					lösungsEingabe[tiefe] = i + j*breite
				}
			} else {
				_erfolgreich, _gesamt, _interessante, _lösungstiefe, _lösungsEingabe := machTauschen(brett2, brett2Zeilen, stückzahl2, tiefe+1)
				erfolgreich += _erfolgreich
				gesamt += _gesamt
				interessante += _interessante
				if _erfolgreich > 0 && _lösungstiefe < lösungstiefe {
					lösungstiefe = _lösungstiefe
					lösungsEingabe = _lösungsEingabe
					lösungsEingabe[tiefe] = i + j*breite

				}
			}

		}
	}
	return
}

func löseBrett(brett [][]uint8, brettZeilen []uint8, stückzahl int, minverhältnis float32) float32 {
	erfolgreich, gesamt, interessante, lösungstiefe, lösungsEingabe := machTauschen(brett, brettZeilen, stückzahl, 0)

	if erfolgreich != 1 {
		return minverhältnis
	}
	if lösungstiefe > 0 {
		return minverhältnis
	}
	if interessante > 1 {
		return minverhältnis
	}
	verhältnis := float32(stückzahl)
	if minverhältnis >= verhältnis {
		return minverhältnis
	}

	druckBrett(brett)
	fmt.Printf("(\nLösung: ")
	for _, eingabe := range lösungsEingabe {
		if eingabe == -1 {
			break
		}
		var x = eingabe % breite
		var y = eingabe / breite

		fmt.Printf("(%d,%d) ", x, y)
	}
	fmt.Printf("\n")

	fmt.Printf("Verhältnis3: %f\nStückzahl: %d\nErfolg: %d\nGesamt: %d\nInteressante: %d\nLösungstiefe: %d\n)\n\n\n", verhältnis, stückzahl, erfolgreich, gesamt, interessante, lösungstiefe)

	return verhältnis

}

func transpose(slice [][]uint8) [][]uint8 {
	xl := len(slice[0])
	yl := len(slice)
	result := make([][]uint8, xl)
	for i := range result {
		result[i] = make([]uint8, yl)
	}
	for i := 0; i < xl; i++ {
		for j := 0; j < yl; j++ {
			result[i][j] = slice[j][i]
		}
	}
	return result
}

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

var masken [][][]uint8
var maskenZeilen [][]uint8

var bretter [][][]uint8
var bretterZeilen [][]uint8

func main() {

	rand.Seed(time.Now().UnixNano())

	breitePtr := flag.Int("Breite", 3, "")
	höhePtr := flag.Int("Höhe", 3, "")
	zahlmaxPtr := flag.Int("Zahlmax", 6, "")
	zielsummePtr := flag.Int("Summe", 8, "")
	zielstückzahlPtr := flag.Int("Stückzahl", -1, "")

	flag.Parse()

	profiling := false
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
		profiling = true
	}

	breite = *breitePtr
	höhe = *höhePtr
	zahlmax = *zahlmaxPtr
	zielsumme = *zielsummePtr
	zielstückzahl = *zielstückzahlPtr

	masken = make([][][]uint8, maxTiefe+1)
	maskenZeilen = make([][]uint8, maxTiefe+1)

	bretter = make([][][]uint8, maxTiefe+1)
	bretterZeilen = make([][]uint8, maxTiefe+1)

	for i := 0; i <= maxTiefe; i++ {
		masken[i], maskenZeilen[i] = leeresBrett()
		bretter[i], bretterZeilen[i] = leeresBrett()
	}

	var minverhältnis float32 = 0.00001

	// {
	// 	brett := transpose([][]uint8{
	// 		{4, 0, 2},
	// 		{2, 0, 5},
	// 		{1, 0, 4},
	// 	})
	// 	stückzahl := 6
	// 	löseBrett(brett, stückzahl, minverhältnis)
	// }

	brett, brettZeilen := leeresBrett()
	maske, _ := leeresBrett()

	for zz := 0; (!profiling) || (zz < 10000); zz++ {

		stückzahl := basteleSilhouette(brett)
		if zielstückzahl != -1 && stückzahl != zielstückzahl {
			continue
		}

		stabilät := stabil(brett, maske)

		if !stabilät {
			continue
		}

		minverhältnis = löseBrett(brett, brettZeilen, stückzahl, minverhältnis)

	}
}
