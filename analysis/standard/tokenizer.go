package standard

import (
	. "github.com/balzaczyy/golucene/core/analysis"
	. "github.com/balzaczyy/golucene/core/analysis/tokenattributes"
	"github.com/balzaczyy/golucene/core/util"
	"io"
)

// standard/StandardTokenizer.java

const (
	ACRONYM_DEP = 8 // deprecated 3.1
)

/* String token types that correspond to token type int constants */
var TOKEN_TYPES = []string{
	"<ALPHANUM>",
	"<APOSTROPHE>",
	"<ACRONYM>",
	"<COMPANY>",
	"<EMAIL>",
	"<HOST>",
	"<NUM>",
	"<CJ>",
	"<ACRONYM_DEP>",
	"<SOUTHEAST_ASIAN>",
	"<IDEOGRAPHIC>",
	"<HIRAGANA>",
	"<KATAKANA>",
	"<HANGUL>",
}

/*
A grammar-based tokenizer constructed with JFlex.

As of Lucene version 3.1, this class implements the Word Break rules
from the Unicode Text Segmentation algorithm, as specified in Unicode
standard Annex #29.

Many applications have specific tokenizer needs. If this tokenizer
does not suit your application, please consider copying this source
code directory to your project and maintaining your own grammar-based
tokenizer.

Version

You must specify the required Version compatibility when creating
StandardTokenizer:

	- As of 3.4, Hiragana and Han characters are no longer wrongly
	split from their combining characters. If you use a previous
	version number, you get the exact broken behavior for backwards
	compatibility.
	- As of 3.1, StandardTokenizer implements Unicode text segmentation.
	If you use a previous version number, you get the exact behavior of
	ClassicTokenizer for backwards compatibility.
*/
type StandardTokenizer struct {
	*Tokenizer
	input io.ReadCloser

	// A private instance of the JFlex-constructed scanner
	scanner StandardTokenizerInterface

	skippedPositions int
	maxTokenLength   int

	// this tokenizer generates three attributes:
	// term offset, positionIncrement and type

	termAtt    CharTermAttribute
	offsetAtt  OffsetAttribute
	posIncrAtt PositionIncrementAttribute
	typeAtt    TypeAttribute
}

/*
Creates a new instance of the StandardTokenizer. Attaches the input
to the newly created JFlex scanner.
*/
func newStandardTokenizer(matchVersion util.Version, input io.ReadCloser) *StandardTokenizer {
	ans := &StandardTokenizer{
		Tokenizer: NewTokenizer(input),
		input:     input,
	}
	ans.init(matchVersion)
	return ans
}

func (t *StandardTokenizer) init(matchVersion util.Version) {
	// GoLucene support >=4.5 only
	t.scanner = newStandardTokenizerImpl(nil)
}

func (t *StandardTokenizer) IncrementToken() (bool, error) {
	t.Attributes().Clear()
	t.skippedPositions = 0

	for {
		tokenType, err := t.scanner.nextToken()
		if tokenType == YYEOF || err != nil {
			return false, err
		}

		if t.scanner.yylength() <= t.maxTokenLength {
			t.posIncrAtt.SetPositionIncrement(t.skippedPositions + 1)
			t.scanner.text(t.termAtt)
			start := t.scanner.yychar()
			t.offsetAtt.SetOffset(t.CorrectOffset(start), t.CorrectOffset(start+t.termAtt.Length()))
			// This 'if' should be removed in the next release. For now,
			// it converts invalid acronyms to HOST. When removed, only the
			// 'else' part should remain.
			if tokenType == ACRONYM_DEP {
				panic("not implemented yet")
			} else {
				t.typeAtt.SetType(TOKEN_TYPES[tokenType])
			}
			return true, nil
		} else {
			// When we skip a too-long term, we still increment the positionincrement
			t.skippedPositions++
		}
	}
}

func (t *StandardTokenizer) End() error {
	panic("not implemented yet")
}

func (t *StandardTokenizer) Reset() error {
	t.scanner.yyreset(t.input)
	t.skippedPositions = 0
	return nil
}

// standard/StandardTokenizerInterface.java

/* This character denotes the end of file */
const YYEOF = -1

/* Internal interface for supporting versioned grammars. */
type StandardTokenizerInterface interface {
	// Copies the matched text into the CharTermAttribute
	text(CharTermAttribute)
	// Returns the current position.
	yychar() int
	// Resets the scanner to read from a new input stream.
	// Does not close the old reader.
	//
	// All internal variables are reset, the old input stream cannot be
	// reused (internal buffer) is discarded and lost). Lexical state
	// is set to ZZ_INITIAL.
	yyreset(io.ReadCloser)
	// Returns the length of the matched text region.
	yylength() int
	// Resumes scanning until the next regular expression is matched,
	// the end of input is encountered or an I/O-Error occurs.
	nextToken() (int, error)
}