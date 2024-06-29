package geojson

import (
	"encoding/json"
	"errors"
	"fmt"
)

type Point Coordinate

func (p Point) Type() GeometryType {
	return PointGeometryType
}

type intermediateGeometry struct {
	Type        string `json:"type" bson:"type"`
	Coordinates any    `json:"coordinates" bson:"coordinates"`
}

func defaultMarshal(g GeometryType, c any) ([]byte, error) {
	return json.Marshal(intermediateGeometry{
		Type:        string(g),
		Coordinates: c,
	})
}

func (p Point) MarshalJSON() ([]byte, error) {
	return defaultMarshal(p.Type(), Coordinate(p))
}

func (p *Point) UnmarshalJSON(data []byte) error {
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}

	if m["type"] != string(PointGeometryType) {
		return fmt.Errorf("unexpected geometry type %q", m["type"])
	}

	l, ok := m["coordinates"].([]float64)
	if !ok {
		return errors.New("incorrect coordinate format")
	}

	*p = Point{Longitude: l[0], Latitude: l[1]}
	return nil
}

type MultiPoint []Coordinate

func (m MultiPoint) Type() GeometryType {
	return MultiPointGeometryType
}

func (p MultiPoint) MarshalJSON() ([]byte, error) {
	return defaultMarshal(p.Type(), []Coordinate(p))
}

func (p *MultiPoint) UnmarshalJSON(data []byte) error {
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}

	if m["type"] != string(MultiPointGeometryType) {
		return fmt.Errorf("unexpected geometry type %q", m["type"])
	}

	coords, ok := m["coordinates"]
	if !ok {
		return errors.New("missing coordinates")
	}

	points, err := getLineOrMultipoint(coords)
	*p = points
	return err
}
