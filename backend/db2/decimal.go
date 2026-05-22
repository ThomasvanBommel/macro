package db2

import (
	"encoding/json"
	"math"
)

// Decimal represents a decimal number with accuracy to two decimal places.
type Decimal float64

const ScaleFactor = 100

func (d Decimal) MarshalJSON() ([]byte, error) {
	return json.Marshal(float64(d))
}

func (d *Decimal) UnmarshalJSON(data []byte) error {
	var f float64
	if err := json.Unmarshal(data, &f); err != nil {
		return err
	}
	*d = Decimal(f)
	return nil
}

func (d Decimal) ToInt() int {
	return int(math.Round(float64(d) * ScaleFactor))
}

func ToDecimal(i int) Decimal {
	return Decimal(float64(i) / ScaleFactor)
}
