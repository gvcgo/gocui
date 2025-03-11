// Copyright 2014 The gocui Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gocui

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteRunes(t *testing.T) {
	tests := []struct {
		existingLines  []string
		stringsToWrite []string
		expectedLines  []string
	}{
		{
			[]string{},
			[]string{""},
			[]string{""},
		},
		{
			[]string{},
			[]string{"1\n"},
			[]string{"1\x00"},
		},
		{
			[]string{},
			[]string{"1\n", "2\n"},
			[]string{"1\x00", "2\x00"},
		},
		{
			[]string{"a"},
			[]string{"1\n"},
			[]string{"1\x00"},
		},
		{
			[]string{"a\x00"},
			[]string{"1\n"},
			[]string{"1\x00"},
		},
		{
			[]string{"ab"},
			[]string{"1\n"},
			[]string{"1b"},
		},
		{
			[]string{"abc"},
			[]string{"1\n"},
			[]string{"1bc"},
		},
		{
			[]string{},
			[]string{"1\r"},
			[]string{"1\x00"},
		},
		{
			[]string{"a"},
			[]string{"1\r"},
			[]string{"1\x00"},
		},
		{
			[]string{"a\x00"},
			[]string{"1\r"},
			[]string{"1\x00"},
		},
		{
			[]string{"ab"},
			[]string{"1\r"},
			[]string{"1b"},
		},
		{
			[]string{"abc"},
			[]string{"1\r"},
			[]string{"1bc"},
		},
	}

	for _, test := range tests {
		v := NewView("name", 0, 0, 10, 10, OutputNormal)
		for _, l := range test.existingLines {
			v.lines = append(v.lines, stringToCells(l))
		}
		for _, s := range test.stringsToWrite {
			v.writeRunes([]rune(s))
		}
		var resultingLines []string
		for _, l := range v.lines {
			resultingLines = append(resultingLines, cellsToString(l))
		}
		assert.Equal(t, test.expectedLines, resultingLines)
	}
}

func TestUpdatedCursorAndOrigin(t *testing.T) {
	tests := []struct {
		prevOrigin     int
		size           int
		cursor         int
		expectedCursor int
		expectedOrigin int
	}{
		{0, 10, 0, 0, 0},
		{0, 10, 10, 10, 0},
		{0, 10, 11, 10, 1},
		{0, 10, 20, 10, 10},
		{20, 10, 19, 0, 19},
		{20, 10, 25, 5, 20},
	}

	for _, test := range tests {
		cursor, origin := updatedCursorAndOrigin(test.prevOrigin, test.size, test.cursor)
		assert.EqualValues(t, test.expectedCursor, cursor, "Cursor is wrong")
		assert.EqualValues(t, test.expectedOrigin, origin, "Origin in wrong")
	}
}

func TestContainsColoredText(t *testing.T) {
	hexColor := func(text string, hexStr string) []cell {
		cells := make([]cell, len(text))
		hex := GetColor(hexStr)
		for i, chr := range text {
			cells[i] = cell{fgColor: hex, chr: chr}
		}
		return cells
	}
	red := "#ff0000"
	green := "#00ff00"
	redStr := func(text string) []cell { return hexColor(text, red) }
	greenStr := func(text string) []cell { return hexColor(text, green) }

	concat := func(lines ...[]cell) []cell {
		var cells []cell
		for _, line := range lines {
			cells = append(cells, line...)
		}
		return cells
	}

	tests := []struct {
		lines      [][]cell
		fgColorStr string
		text       string
		expected   bool
	}{
		{
			lines:      [][]cell{concat(redStr("a"))},
			fgColorStr: red,
			text:       "a",
			expected:   true,
		},
		{
			lines:      [][]cell{concat(redStr("a"))},
			fgColorStr: red,
			text:       "b",
			expected:   false,
		},
		{
			lines:      [][]cell{concat(redStr("a"))},
			fgColorStr: green,
			text:       "b",
			expected:   false,
		},
		{
			lines:      [][]cell{concat(redStr("hel"), greenStr("lo"), redStr(" World!"))},
			fgColorStr: red,
			text:       "hello",
			expected:   false,
		},
		{
			lines:      [][]cell{concat(redStr("hel"), greenStr("lo"), redStr(" World!"))},
			fgColorStr: green,
			text:       "lo",
			expected:   true,
		},
		{
			lines: [][]cell{
				redStr("hel"),
				redStr("lo"),
			},
			fgColorStr: red,
			text:       "hello",
			expected:   false,
		},
	}

	for i, test := range tests {
		v := &View{lines: test.lines}
		assert.Equal(t, test.expected, v.ContainsColoredText(test.fgColorStr, test.text), "Test %d failed", i)
	}
}

func stringToCells(s string) []cell {
	var cells []cell
	for _, c := range s {
		cells = append(cells, cell{chr: c})
	}
	return cells
}

func cellsToString(cells []cell) string {
	var s string
	for _, c := range cells {
		s += string(c.chr)
	}
	return s
}

func TestLineWrap(t *testing.T) {
	testCases := []struct {
		name     string
		line     string
		columns  int
		expected []string
	}{
		{
			name:    "Wrap on space",
			line:    "Hello World",
			columns: 5,
			expected: []string{
				"Hello",
				"World",
			},
		},
		{
			name:    "Wrap on hyphen",
			line:    "Hello-World",
			columns: 6,
			expected: []string{
				"Hello-",
				"World",
			},
		},
		{
			name:    "Wrap on hyphen 2",
			line:    "Blah Hello-World",
			columns: 12,
			expected: []string{
				"Blah Hello-",
				"World",
			},
		},
		{
			name:    "Wrap on hyphen 3",
			line:    "Blah Hello-World",
			columns: 11,
			expected: []string{
				"Blah Hello-",
				"World",
			},
		},
		{
			name:    "Wrap on hyphen 4",
			line:    "Blah Hello-World",
			columns: 10,
			expected: []string{
				"Blah Hello",
				"-World",
			},
		},
		{
			name:    "Wrap on space 2",
			line:    "Blah Hello World",
			columns: 10,
			expected: []string{
				"Blah Hello",
				"World",
			},
		},
		{
			name:    "Wrap on space with more words",
			line:    "Longer word here",
			columns: 10,
			expected: []string{
				"Longer",
				"word here",
			},
		},
		{
			name:    "Split word that's too long",
			line:    "ThisWordIsWayTooLong",
			columns: 10,
			expected: []string{
				"ThisWordIs",
				"WayTooLong",
			},
		},
		{
			name:    "Split word that's too long over multiple lines",
			line:    "ThisWordIsWayTooLong",
			columns: 5,
			expected: []string{
				"ThisW",
				"ordIs",
				"WayTo",
				"oLong",
			},
		},
		{
			name:    "Lots of hyphens",
			line:    "one-two-three-four-five",
			columns: 8,
			expected: []string{
				"one-two-",
				"three-",
				"four-",
				"five",
			},
		},
		{
			name:    "Several lines using all the available width",
			line:    "aaa bb cc ddd-ee ff",
			columns: 5,
			expected: []string{
				"aaa",
				"bb cc",
				"ddd-",
				"ee ff",
			},
		},
		{
			name:    "Multi-cell runes",
			line:    "🐤🐤🐤 🐝🐝 🙉 🦊🦊🦊-🐬🐬 🦢🦢",
			columns: 9,
			expected: []string{
				"🐤🐤🐤",
				"🐝🐝 🙉",
				"🦊🦊🦊-",
				"🐬🐬 🦢🦢",
			},
		},
		{
			name:    "Space in last column",
			line:    "hello world",
			columns: 6,
			expected: []string{
				"hello",
				"world",
			},
		},
		{
			name:    "Hyphen in last column",
			line:    "hello-world",
			columns: 6,
			expected: []string{
				"hello-",
				"world",
			},
		},
		{
			name:    "English text",
			line:    "+The sea reach of the Thames stretched before us like the bedinnind of an interminable waterway. In the offind the sea and the sky were welded todether without a joint, and in the luminous space the tanned sails of the bardes drifting blah blah",
			columns: 81,
			expected: []string{
				"+The sea reach of the Thames stretched before us like the bedinnind of an",
				"interminable waterway. In the offind the sea and the sky were welded todether",
				"without a joint, and in the luminous space the tanned sails of the bardes",
				"drifting blah blah",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			lineCells := stringToCells(tc.line)

			result := lineWrap(lineCells, tc.columns)

			resultStrings := make([]string, len(result))
			for i, line := range result {
				resultStrings[i] = cellsToString(line)
			}

			assert.EqualValues(t, tc.expected, resultStrings)
		})
	}
}
