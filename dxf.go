package dxf

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/ma2ma/godxf/color"
	"github.com/ma2ma/godxf/drawing"
	"github.com/ma2ma/godxf/table"
)

// Default values.
var (
	DefaultColor    = color.White
	DefaultLineType = table.LT_CONTINUOUS
)

// NewDrawing creates a drawing.
func NewDrawing() *drawing.Drawing {
	return drawing.New()
}

// Open opens the named DXF file.
func OpenReader(r io.Reader) (*drawing.Drawing, error) {
	var err error
	//defer r.Close()
	scanner := bufio.NewScanner(r)
	d := NewDrawing()
	var code, value string
	parsers := []func(*drawing.Drawing, int, [][2]string) error{
		ParseHeader,
		ParseClasses,
		ParseTables,
		ParseBlocks,
		ParseEntities,
		ParseObjects,
	}
	data := make([][2]string, 0)
	setparser := false
	var parser func(*drawing.Drawing, int, [][2]string) error
	line := 0
	startline := 0
	for scanner.Scan() {
		line++
		if line%2 == 1 {
			code = strings.TrimSpace(scanner.Text())
			if err != nil {
				return d, err
			}
		} else {
			value = scanner.Text()
			if setparser {
				if code != "2" {
					return d, fmt.Errorf("line %d: invalid group code: %s", line, code)
				}
				ind := drawing.SectionTypeValue(strings.ToUpper(value))
				if ind < 0 {
					return d, fmt.Errorf("line %d: unknown section name: %s", line, value)
				}
				parser = parsers[ind]
				startline = line + 1
				setparser = false
			} else {
				if code == "0" {
					switch strings.ToUpper(value) {
					case "EOF":
						return d, nil
					case "SECTION":
						setparser = true
					case "ENDSEC":
						err := parser(d, startline, data)
						if err != nil {
							return d, err
						}
						data = make([][2]string, 0)
						startline = line + 1
					default:
						data = append(data, [2]string{code, scanner.Text()})
					}
				} else {
					data = append(data, [2]string{code, scanner.Text()})
				}
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return d, err
	}
	if len(data) > 0 {
		err := parser(d, startline, data)
		if err != nil {
			return d, err
		}
	}
	return d, nil
}

// Open opens the named DXF file.
func Open(filename string) (*drawing.Drawing, error) {
	var err error
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	d := NewDrawing()
	var code, value string
	parsers := []func(*drawing.Drawing, int, [][2]string) error{
		ParseHeader,
		ParseClasses,
		ParseTables,
		ParseBlocks,
		ParseEntities,
		ParseObjects,
	}
	data := make([][2]string, 0)
	setparser := false
	var parser func(*drawing.Drawing, int, [][2]string) error
	line := 0
	startline := 0
	for scanner.Scan() {
		line++
		if line%2 == 1 {
			code = strings.TrimSpace(scanner.Text())
			if err != nil {
				return d, err
			}
		} else {
			value = scanner.Text()
			if setparser {
				if code != "2" {
					return d, fmt.Errorf("line %d: invalid group code: %s", line, code)
				}
				ind := drawing.SectionTypeValue(strings.ToUpper(value))
				if ind < 0 {
					return d, fmt.Errorf("line %d: unknown section name: %s", line, value)
				}
				parser = parsers[ind]
				startline = line + 1
				setparser = false
			} else {
				if code == "0" {
					switch strings.ToUpper(value) {
					case "EOF":
						return d, nil
					case "SECTION":
						setparser = true
					case "ENDSEC":
						err := parser(d, startline, data)
						if err != nil {
							return d, err
						}
						data = make([][2]string, 0)
						startline = line + 1
					default:
						data = append(data, [2]string{code, scanner.Text()})
					}
				} else {
					data = append(data, [2]string{code, scanner.Text()})
				}
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return d, err
	}
	if len(data) > 0 {
		err := parser(d, startline, data)
		if err != nil {
			return d, err
		}
	}
	return d, nil
}

// ColorIndex converts RGB value to corresponding color number.
func ColorIndex(cl []int) color.ColorNumber {
	minind := 0
	minval := 1000000
	for i, c := range color.ColorRGB {
		tmpval := 0
		for j := 0; j < 3; j++ {
			tmpval += (cl[j] - int(c[j])) * (cl[j] - int(c[j]))
		}
		if tmpval < minval {
			minind = i
			minval = tmpval
			if minval == 0 {
				break
			}
		}
	}
	return color.ColorNumber(minind)
}

// IndexColor converts color number to RGB value.
func IndexColor(index uint8) []uint8 {
	return color.ColorRGB[index]
}
