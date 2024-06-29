package geojson

import (
	"encoding/json"
	"errors"
	"fmt"
)

// GeometryType is a supported geometry type.
type GeometryType string

// Types of geometry.
const (
	PointGeometryType           GeometryType = "Point"
	MultiPointGeometryType      GeometryType = "MultiPoint"
	LineStringGeometryType      GeometryType = "LineString"
	MultiLineStringGeometryType GeometryType = "MultiLineString"
	PolygonGeometryType         GeometryType = "Polygon"
	MultiPolygonGeometryType    GeometryType = "MultiPolygon"
	GeometryCollectionType      GeometryType = "GeometryCollection"
)

type Geometry interface {
	Type() GeometryType
}

type Coordinate struct {
	Latitude  float64
	Longitude float64
}

func (c *Coordinate) UnmarshalJSON(data []byte) error {
	var list [2]float64
	if err := json.Unmarshal(data, &list); err != nil {
		return err
	}

	*c = Coordinate{Longitude: list[0], Latitude: list[1]}
	return nil
}

func (c Coordinate) MarshalJSON() ([]byte, error) {
	return json.Marshal([]float64{c.Longitude, c.Latitude})
}

type GeometryCollection struct {
	GT         GeometryType `json:"type"`
	Geometries []Geometry   `json:"geometries"`
}

func (g GeometryCollection) Type() GeometryType {
	return GeometryCollectionType
}

func (g *GeometryCollection) UnmarshalJSON(data []byte) error {
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}

	if m["type"] != string(GeometryCollectionType) {
		return fmt.Errorf("unexpected geometry type %s", m["type"])
	}

	geomAnys, exists := m["geometries"]
	if !exists {
		return fmt.Errorf("geometries must not be nil")
	}

	geoms, ok := normalizeToSlice(geomAnys)
	if !ok {
		return fmt.Errorf("geometries must be a slice")
	}

	newGeometries := make([]Geometry, 0, len(geoms))
	for _, geom := range geoms {
		remarshaled, err := json.Marshal(geom)
		if err != nil {
			return err
		}

		geomap, ok := geom.(map[string]any)
		if !ok {
			return fmt.Errorf("invalid geometry '%s'", string(remarshaled))
		}

		geotype, typed := geomap["type"]
		if !typed {
			return errors.New("geometry type must be specified")
		}

		switch GeometryType(fmt.Sprintf("%v", geotype)) {
		case PointGeometryType:
			var p Point
			if err = json.Unmarshal(remarshaled, &p); err != nil {
				return err
			}
			newGeometries = append(newGeometries, p)
		case MultiPointGeometryType:
			var p MultiPoint
			if err = json.Unmarshal(remarshaled, &p); err != nil {
				return err
			}
			newGeometries = append(newGeometries, p)
		case LineStringGeometryType:
			var l LineString
			if err = json.Unmarshal(remarshaled, &l); err != nil {
				return err
			}
			newGeometries = append(newGeometries, l)
		case MultiLineStringGeometryType:
			var l MultiLineString
			if err = json.Unmarshal(remarshaled, &l); err != nil {
				return err
			}
			newGeometries = append(newGeometries, l)
		case PolygonGeometryType:
			var p Polygon
			if err = json.Unmarshal(remarshaled, &p); err != nil {
				return err
			}
			newGeometries = append(newGeometries, p)
		case MultiPolygonGeometryType:
			var p MultiPolygon
			if err = json.Unmarshal(remarshaled, &p); err != nil {
				return err
			}
			newGeometries = append(newGeometries, p)
		default:
			return fmt.Errorf("unrecognized geometry type %v", geotype)
		}
	}

	*g = GeometryCollection{
		GT:         GeometryCollectionType,
		Geometries: newGeometries,
	}

	return nil
}
