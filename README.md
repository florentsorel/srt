# SRT – SubRip Parser for Go

A Go library for reading, writing, and manipulating SubRip (SRT) subtitle files.

---

## Features

- Parse SRT files from file path.
- Parse SRT files from `io.Reader`.
- Work with a simple data model: `model.Subtitles` and `model.Cue`.
- Shift subtitles in time, remove cues, or re-serialize back to SRT.
- UTF-8 only: supports clean parsing and writing without hidden conversions.
---

## Installation

Requires **Go 1.18+**.

```bash
go get github.com/florentsorel/srt
```
