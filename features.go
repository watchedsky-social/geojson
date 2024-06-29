package geojson

import (
	"encoding/json"
	"fmt"
)

type JSONObject map[string]any

func (j JSONObject) StringValue(key string) string {
	s, _ := j[key].(string)
	return s
}

func (j JSONObject) IntValue(key string) int64 {
	i, _ := j[key].(int64)
	return i
}

func (j JSONObject) FloatValue(key string) float64 {
	f, _ := j[key].(float64)
	return f
}

type Feature struct {
	ID         string     `json:"id"`
	Geometry   Geometry   `json:"geometry"`
	Properties JSONObject `json:"properties"`
}

type FeatureCollection struct {
	Features []Feature `json:"features"`
}

type intermediateFeature struct {
	ID         string         `json:"id"`
	Geometry   map[string]any `json:"geometry"`
	Properties map[string]any `json:"properties"`
}

func (f *Feature) UnmarshalJSON(data []byte) error {
	var intF intermediateFeature
	if err := json.Unmarshal(data, &intF); err != nil {
		return err
	}

	(*f).ID = intF.ID
	(*f).Properties = intF.Properties

	if intF.Geometry != nil {
		if rawType, ok := intF.Geometry["type"]; ok {
			if geoType, ok := rawType.(string); ok {
				geoBytes, err := json.Marshal(intF.Geometry)
				if err != nil {
					return err
				}

				switch GeometryType(geoType) {
				case PointGeometryType:
					g := Point{}
					if err = json.Unmarshal(geoBytes, &g); err != nil {
						return err
					}

					(*f).Geometry = g
				case MultiPointGeometryType:
					g := MultiPoint{}
					if err = json.Unmarshal(geoBytes, &g); err != nil {
						return err
					}

					(*f).Geometry = g
				case LineStringGeometryType:
					g := LineString{}
					if err = json.Unmarshal(geoBytes, &g); err != nil {
						return err
					}

					(*f).Geometry = g
				case MultiLineStringGeometryType:
					g := MultiLineString{}
					if err = json.Unmarshal(geoBytes, &g); err != nil {
						return err
					}

					(*f).Geometry = g
				case PolygonGeometryType:
					g := Polygon{}
					if err = json.Unmarshal(geoBytes, &g); err != nil {
						return err
					}

					(*f).Geometry = g
				case MultiPolygonGeometryType:
					g := MultiPolygon{}
					if err = json.Unmarshal(geoBytes, &g); err != nil {
						return err
					}

					(*f).Geometry = g
				case GeometryCollectionType:
					g := GeometryCollection{}
					if err = json.Unmarshal(geoBytes, &g); err != nil {
						return err
					}

					(*f).Geometry = g
				default:
					return fmt.Errorf("unrecognized geometry type %q", geoType)
				}

			}
		}
	}

	return nil
}
