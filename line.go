package geojson

import (
	"encoding/json"
	"errors"
	"fmt"
)

type LineString []Coordinate

func (l LineString) Type() GeometryType {
	return LineStringGeometryType
}

func (p LineString) MarshalJSON() ([]byte, error) {
	return defaultMarshal(p.Type(), []Coordinate(p))
}

func (l *LineString) UnmarshalJSON(data []byte) error {
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}

	if m["type"] != string(LineStringGeometryType) {
		return fmt.Errorf("unexpected geometry type %q", m["type"])
	}

	coords, ok := m["coordinates"]
	if !ok {
		return errors.New("missing coordinates")
	}

	line, err := getLineOrMultipoint(coords)
	*l = line
	return err
}

type MultiLineString [][]Coordinate

func (m MultiLineString) Type() GeometryType {
	return MultiLineStringGeometryType
}

func (m MultiLineString) MarshalJSON() ([]byte, error) {
	return defaultMarshal(m.Type(), [][]Coordinate(m))
}

func (m *MultiLineString) UnmarshalJSON(data []byte) error {
	var d map[string]any
	if err := json.Unmarshal(data, &d); err != nil {
		return err
	}

	if d["type"] != string(MultiLineStringGeometryType) {
		return fmt.Errorf("unexpected geometry type %q", d["type"])
	}

	coords, ok := d["coordinates"]
	if !ok {
		return errors.New("missing coordinates")
	}

	lines, err := getPolygonOrMultiline(coords)
	*m = lines
	return err
}
