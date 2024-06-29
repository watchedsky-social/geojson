package geojson

import (
	"encoding/json"
	"errors"
	"fmt"
)

type Polygon [][]Coordinate

func (p Polygon) Type() GeometryType {
	return PolygonGeometryType
}

func (p Polygon) MarshalJSON() ([]byte, error) {
	return defaultMarshal(p.Type(), [][]Coordinate(p))
}

func (p *Polygon) UnmarshalJSON(data []byte) error {
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}

	if m["type"] != string(PolygonGeometryType) {
		return fmt.Errorf("unexpected geometry type %q", m["type"])
	}

	coords, ok := m["coordinates"]
	if !ok {
		return errors.New("missing coordinates")
	}

	polygon, err := getPolygonOrMultiline(coords)
	*p = polygon
	return err
}

type MultiPolygon [][][]Coordinate

func (m MultiPolygon) Type() GeometryType {
	return MultiPolygonGeometryType
}

func (p MultiPolygon) MarshalJSON() ([]byte, error) {
	return defaultMarshal(p.Type(), [][][]Coordinate(p))
}

func (m *MultiPolygon) UnmarshalJSON(data []byte) error {
	var d map[string]any
	if err := json.Unmarshal(data, &d); err != nil {
		return err
	}

	if d["type"] != string(MultiPolygonGeometryType) {
		return fmt.Errorf("unexpected geometry type %q", d["type"])
	}

	coords, ok := d["coordinates"]
	if !ok {
		return errors.New("missing coordinates")
	}

	polygons, err := getMultiPolygon(coords)
	*m = polygons
	return err
}
