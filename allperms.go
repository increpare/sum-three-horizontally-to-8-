//z.B. go run allperms.go  3 6 8

package main

func partition(zellen int, maxziffer int, summe int) [][]int {
	if zellen == 1 {
		if maxziffer >= summe {
			return [][]int{{summe}}
		}

		return [][]int{}
	}

	result := make([][]int, 0)

	for z := 1; z <= maxziffer; z++ {
		neusumme := summe - z
		if neusumme < (zellen - 1) {
			continue
		}

		unterpartitionen := partition(zellen-1, maxziffer, neusumme)
		for _, p := range unterpartitionen {
			neueZeile := append([]int{z}, p...)
			result = append(result, neueZeile)
		}
	}
	return result
}

func rückwärtsString(s string) string {
	zeichen := []rune(s)
	for i, j := 0, len(zeichen)-1; i < j; i, j = i+1, j-1 {
		zeichen[i], zeichen[j] = zeichen[j], zeichen[i]
	}
	return string(zeichen)
}

// func main() {
// 	A := os.Args[1:]
// 	if len(A) == 0 {
// 		fmt.Printf("allperms ZELLEN MAXZIFFER SUMME\n")
// 		return
// 	}

// 	zellen, _ := strconv.Atoi(A[0])
// 	maxziffer, _ := strconv.Atoi(A[1])
// 	summe, _ := strconv.Atoi(A[2])

// 	partitionen := partition(zellen, maxziffer, summe)

// 	for _, p := range partitionen {
// 		//		horizontal [player no spawnfall] [ k1 | k3 | k4 ] -> [player no spawnfall] [ k1 dokill | k3 dokill | k4 dokill]
// 		s := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(p)), ""), "[]")
// 		srv := rückwärtsString(s)
// 		if s == srv {
// 			fmt.Printf("right ")
// 		} else if s < srv {
// 			fmt.Printf("horizontal ")
// 		} else {
// 			continue
// 		}
// 		k1str := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(p)), " | k"), "[]")
// 		k2str := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(p)), " dokill | k"), "[]")
// 		fmt.Printf("[player no spawnfall] [ k%s ] -> [player no spawnfall] [ k%s dokill ]\n", k1str, k2str)
// 		// fmt.Printf("%v\n", p)
// 	}
// }
