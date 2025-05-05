package cyrlat

import (
	"strings"
	"sync"
	"unicode/utf8"
)

// var isMongolian = regexp.MustCompile(`^[А-Яа-я0-9ЁёӨөҮү\-.,!?:"' ]+$`).MatchString

// var isMongolianAbbr = regexp.MustCompile(`[А-ЯЁӨҮ]{2,}`).MatchString

var charMap = map[rune][]string{
	'а': {"a"},
	'б': {"b"},
	'в': {"v", "w"},
	'г': {"g"},
	'д': {"d"},
	'е': {"ye"},
	'ё': {"yo"},
	'ж': {"j"},
	'з': {"z"},
	'и': {"i"},
	'й': {"i"},
	'к': {"k"},
	'л': {"l"},
	'м': {"m"},
	'н': {"n"},
	'о': {"o"},
	'ө': {"o", "u"},
	'п': {"p"},
	'р': {"r"},
	'с': {"s"},
	'т': {"t"},
	'у': {"u"},
	'ү': {"u"},
	'ф': {"f", "p"},
	'х': {"h", "kh"},
	'ц': {"c", "ts"},
	'ч': {"c", "ch"},
	'ш': {"sh"},
	'щ': {"shch"},
	'ь': {"i"},
	'ъ': {""},
	'ы': {"ii"},
	'э': {"e"},
	'ю': {"y", "yu"},
	'я': {"y", "ya"},
}

func GetLatins(word string) []string {
	var wg sync.WaitGroup
	ch := make(chan *string)

	wg.Add(1)
	go convert(&wg, ch, "", strings.ToLower(word), map[rune]string{})
	go func() {
		wg.Wait()
		ch <- nil
	}()

	page := []string{}
	for {
		itm := <-ch
		if itm == nil {
			break
		}
		page = append(page, *itm)
	}
	return page
}

func convert(wg *sync.WaitGroup, ch chan *string, cur, rest string, used map[rune]string) {
	defer wg.Done()
	r, size := utf8.DecodeRuneInString(rest)
	if size == 0 {
		tmp := cur
		ch <- &tmp
		return
	}
	if r == utf8.RuneError {
		// If not valid utf-8, skip
		wg.Add(1)
		go convert(wg, ch, cur, trimLeftChars(rest, 1), used)
		return
	}

	copyUsed := map[rune]string{}
	for k2, v2 := range used {
		copyUsed[k2] = v2
	}

	if char, ok := copyUsed[r]; ok {
		wg.Add(1)
		go convert(wg, ch, cur+char, trimLeftChars(rest, 1), copyUsed)
		return
	}

	chars, ok := charMap[r]
	if !ok {
		// Add rune blindly if not in char map.
		wg.Add(1)
		go convert(wg, ch, cur+string(r), trimLeftChars(rest, 1), copyUsed)
		return
	}

	for _, char := range chars {
		copyUsed2 := map[rune]string{}
		for k2, v2 := range used {
			copyUsed2[k2] = v2
		}
		copyUsed2[r] = char
		wg.Add(1)
		go convert(wg, ch, cur+string(char), trimLeftChars(rest, 1), copyUsed2)
	}
}

func trimLeftChars(s string, n int) string {
	m := 0
	for i := range s {
		if m >= n {
			return s[i:]
		}
		m++
	}
	return s[:0]
}

// func main() {
// 	page := GetLatins("Сайнуу өвөө")
// 	for _, word := range page {
// 		log.Println(word)
// 	}
// 	// log.Println(GetLatins("Худалдаа"))
// 	// log.Println(GetLatins("Өглөө"))
// 	// log.Println(GetLatins("Чуцуцаша"))
// }
