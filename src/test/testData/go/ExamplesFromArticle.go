package main

import "strings"

// Go port of the article-style golden cases used in Kotlin/Java/Dart.
// Where Go lacks the original construct (notably try/catch and ternary),
// the closest control-flow equivalent is used while preserving the overall shape.

//go:generate complexity 7
func exampleOne(a, b, c, d, e, f bool) {
	if a && b && c || d || e && f { // +1 for `if`, +1, +1, +1
	}

	if a && // +1 for `if`, +1
		!(b && c) { // +1
	}
}

type ArticleEntry struct {
	version  int
	previous *ArticleEntry
}

type ArticleTxn struct {
	active bool
	status int
}

type ArticleIndex struct{}

func (ArticleIndex) wwDependency(version, status, wait int) int { return 0 }

const articleTimedOut = -1
const articleAborted = 1

//go:generate complexity 24
func exampleTwo(entry *ArticleEntry, txn ArticleTxn, frst *ArticleEntry, ti ArticleIndex) *ArticleEntry {
	for { // +1
		if frst != nil { // +2 (nesting = 1)
			if frst.version > entry.version { // +3 (nesting = 2)
				return nil
			}
			if txn.active { // +3 (nesting = 2)
				e := frst
				for e != nil { // +4 (nesting = 3)
					version := e.version
					depends := ti.wwDependency(version, txn.status, 0)
					if depends == articleTimedOut { // +5 (nesting = 4)
						return nil
					}
					if depends != 0 && depends != articleAborted { // +5 (nesting = 4), +1
						return nil
					}
					e = e.previous
				}
			}
			entry.previous = frst
			frst = entry
			break
		}
	}
	return frst
} // total complexity == 24

//go:generate complexity 9
func exampleThree(a, b, c bool) {
	if a { // +1
		for i := 1; i < 10; i++ { // +2 (nesting=1)
			for b { // +3 (nesting=2)
			}
		}
		if c { // +2 (nesting=1)
		} else { // +1
		}
	}
} // Closest Go equivalent to the article's try/catch example.

//go:generate complexity 2
func exampleFour(a bool) {
	go func() { // +0 (but nesting level is now 1)
		if a { // +2 (nesting=1)
		}
	}()
} // Cognitive Complexity 2

//go:generate complexity 7
func exampleFive(max int) int {
	total := 0
OUT:
	for i := 1; i < max; i++ { // +1
		for j := 2; j < i; j++ { // +2
			if i%j == 0 { // +3
				continue OUT // +1
			}
		}
		total += i
	}
	return total
} // Cognitive Complexity 7

//go:generate complexity 1
func exampleSix(number int) string {
	switch number { // +1
	case 1:
		return "one"
	case 2:
		return "a couple"
	case 3:
		return "a few"
	default:
		return "lots"
	}
} // Cognitive Complexity 1

type ArticleMethodSymbol struct{}

type ArticleOverrideSymbol struct {
	isMethod bool
	isStatic bool
	method   *ArticleMethodSymbol
}

type ArticleSymbolTable struct {
	items []ArticleOverrideSymbol
}

func (t ArticleSymbolTable) lookup(name string) []ArticleOverrideSymbol {
	_ = name
	return t.items
}

type ArticleClassType struct {
	unknown bool
	symbol  ArticleSymbolTable
}

var unknownMethodSymbol = &ArticleMethodSymbol{}

func canOverride(method *ArticleMethodSymbol) bool {
	return method != nil
}

func checkOverridingParameters(method *ArticleMethodSymbol, classType ArticleClassType) *bool {
	_ = method
	_ = classType
	return nil
}

//go:generate complexity 19
func exampleSeven(classType ArticleClassType, name string) *ArticleMethodSymbol {
	if classType.unknown { // +1
		return unknownMethodSymbol
	}
	unknownFound := false
	symbols := classType.symbol.lookup(name)
	for _, overrideSymbol := range symbols { // +1
		if overrideSymbol.isMethod && !overrideSymbol.isStatic { // +2 (nesting = 1), +1
			methodJavaSymbol := overrideSymbol.method
			if canOverride(methodJavaSymbol) { // +3 (nesting = 2)
				overriding := checkOverridingParameters(methodJavaSymbol, classType)
				if overriding == nil { // +4 (nesting = 3)
					if !unknownFound { // +5 (nesting = 4)
						unknownFound = true
					}
				} else if *overriding { // +1
					return methodJavaSymbol
				}
			}
		}
	}
	if unknownFound { // +1
		return unknownMethodSymbol
	}
	return nil
} // total complexity == 19

const articleSpecialChars = ".+()|^$"

func isSlash(ch byte) bool {
	return ch == '/' || ch == '\\'
}

//go:generate complexity 20
func exampleEight(antPattern, directorySeparator string) string {
	escapedDirectorySeparator := "\\" + directorySeparator
	var sb strings.Builder
	sb.Grow(len(antPattern))
	sb.WriteByte('^')

	i := 0
	if strings.HasPrefix(antPattern, "/") || strings.HasPrefix(antPattern, "\\") { // +1, +1
		i = 1
	}
	for i < len(antPattern) { // +1
		ch := antPattern[i]
		if strings.IndexByte(articleSpecialChars, ch) != -1 { // +2 (nesting = 1)
			sb.WriteByte('\\')
			sb.WriteByte(ch)
		} else if ch == '*' { // +1
			if i+1 < len(antPattern) && antPattern[i+1] == '*' { // +3 (nesting = 2), +1
				if i+2 < len(antPattern) && isSlash(antPattern[i+2]) { // +4 (nesting = 3), +1
					sb.WriteString("(?:.*")
					sb.WriteString(escapedDirectorySeparator)
					sb.WriteString("|)")
					i += 2
				} else { // +1
					sb.WriteString(".*")
					i += 1
				}
			} else { // +1
				sb.WriteString("[^")
				sb.WriteString(escapedDirectorySeparator)
				sb.WriteString("]*?")
			}
		} else if ch == '?' { // +1
			sb.WriteString("[^")
			sb.WriteString(escapedDirectorySeparator)
			sb.WriteString("]")
		} else if isSlash(ch) { // +1
			sb.WriteString(escapedDirectorySeparator)
		} else { // +1
			sb.WriteByte(ch)
		}
		i++
	}

	sb.WriteByte('$')
	return sb.String()
} // total complexity = 20
