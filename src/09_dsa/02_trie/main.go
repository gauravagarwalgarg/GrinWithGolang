/*
What this teaches:
    Trie data structure: insert, search, startsWith, and autocomplete. Use cases
    include routing tables, autocomplete systems, and prefix-based searches.

Beginner analogy:
    "A Trie is like a phone's predictive text: each letter narrows down possible
     words. After typing 'hel', the trie knows 'hello', 'help', 'helmet' are all
     reachable it's a tree of prefixes."

C++ comparison:
    "Same structure as C++ using array[26] or unordered_map<char, Node*>. Go uses
     map[rune]*TrieNode for Unicode support. No pointer arithmetic, but the
     algorithmic complexity is identical."

Interview relevance:
    Tries appear in autocomplete, spell-check, IP routing, and word-break problems.
    Interviewers ask for insert/search in O(m) where m is word length, and
    autocomplete with DFS traversal.
*/

package main

import "fmt"

// --- Trie Node ---

type TrieNode struct {
	children map[rune]*TrieNode
	isEnd    bool
	word     string // Store complete word at terminal nodes
}

func newTrieNode() *TrieNode {
	return &TrieNode{children: make(map[rune]*TrieNode)}
}

// --- Trie ---

type Trie struct {
	root *TrieNode
	size int
}

func NewTrie() *Trie {
	return &Trie{root: newTrieNode()}
}

// Insert adds a word to the trie. O(m) where m = len(word).
func (t *Trie) Insert(word string) {
	node := t.root
	for _, ch := range word {
		if _, ok := node.children[ch]; !ok {
			node.children[ch] = newTrieNode()
		}
		node = node.children[ch]
	}
	if !node.isEnd {
		t.size++
	}
	node.isEnd = true
	node.word = word
}

// Search returns true if the exact word exists. O(m).
func (t *Trie) Search(word string) bool {
	node := t.findNode(word)
	return node != nil && node.isEnd
}

// StartsWith returns true if any word has the given prefix. O(m).
func (t *Trie) StartsWith(prefix string) bool {
	return t.findNode(prefix) != nil
}

// Autocomplete returns all words with the given prefix, up to limit.
func (t *Trie) Autocomplete(prefix string, limit int) []string {
	node := t.findNode(prefix)
	if node == nil {
		return nil
	}

	var results []string
	t.collect(node, &results, limit)
	return results
}

// Delete removes a word from the trie. Returns true if it existed.
func (t *Trie) Delete(word string) bool {
	return t.deleteHelper(t.root, []rune(word), 0)
}

func (t *Trie) Size() int {
	return t.size
}

// --- Internal helpers ---

func (t *Trie) findNode(prefix string) *TrieNode {
	node := t.root
	for _, ch := range prefix {
		child, ok := node.children[ch]
		if !ok {
			return nil
		}
		node = child
	}
	return node
}

func (t *Trie) collect(node *TrieNode, results *[]string, limit int) {
	if len(*results) >= limit {
		return
	}
	if node.isEnd {
		*results = append(*results, node.word)
	}
	for _, child := range node.children {
		t.collect(child, results, limit)
	}
}

func (t *Trie) deleteHelper(node *TrieNode, runes []rune, depth int) bool {
	if depth == len(runes) {
		if !node.isEnd {
			return false
		}
		node.isEnd = false
		t.size--
		return len(node.children) == 0
	}

	ch := runes[depth]
	child, ok := node.children[ch]
	if !ok {
		return false
	}

	shouldDelete := t.deleteHelper(child, runes, depth+1)
	if shouldDelete {
		delete(node.children, ch)
		return !node.isEnd && len(node.children) == 0
	}
	return false
}

func main() {
	fmt.Println("=== Trie (Prefix Tree) ===")

	trie := NewTrie()

	// Insert words
	words := []string{"hello", "help", "helmet", "hero", "her", "heap", "go", "golang", "good"}
	fmt.Println("\n--- Inserting words ---")
	for _, w := range words {
		trie.Insert(w)
	}
	fmt.Printf("  Inserted %d words: %v\n", trie.Size(), words)

	// Search
	fmt.Println("\n--- Search ---")
	for _, w := range []string{"hello", "hell", "hero", "golang", "gopher"} {
		fmt.Printf("  Search(%q) = %v\n", w, trie.Search(w))
	}

	// StartsWith
	fmt.Println("\n--- StartsWith ---")
	for _, p := range []string{"hel", "go", "xyz", "he"} {
		fmt.Printf("  StartsWith(%q) = %v\n", p, trie.StartsWith(p))
	}

	// Autocomplete
	fmt.Println("\n--- Autocomplete ---")
	fmt.Printf("  Autocomplete('hel', 5) = %v\n", trie.Autocomplete("hel", 5))
	fmt.Printf("  Autocomplete('go', 3)  = %v\n", trie.Autocomplete("go", 3))
	fmt.Printf("  Autocomplete('he', 10) = %v\n", trie.Autocomplete("he", 10))

	// Delete
	fmt.Println("\n--- Delete ---")
	fmt.Printf("  Before delete: Search('help') = %v\n", trie.Search("help"))
	trie.Delete("help")
	fmt.Printf("  After delete:  Search('help') = %v\n", trie.Search("help"))
	fmt.Printf("  StartsWith('hel') still = %v (other words remain)\n", trie.StartsWith("hel"))
	fmt.Printf("  Size: %d\n", trie.Size())

	// Key takeaways
	fmt.Println("\n--- Key Takeaways ---")
	fmt.Println("1. Insert/Search/StartsWith are all O(m) length of the word")
	fmt.Println("2. Space: O(alphabet_size × total_characters) in worst case")
	fmt.Println("3. Autocomplete: DFS from prefix node collects all terminal descendants")
	fmt.Println("4. map[rune] supports Unicode; array[26] is faster for ASCII-only")
	fmt.Println("5. Use cases: autocomplete, spell check, IP routing, word games")
}
