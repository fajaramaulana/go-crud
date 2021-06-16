package models

// User schema of the user table
type Bbm struct {
	ID           int64  `json:"id"`
	Jumlah_liter string `json:"jumlah_liter"`
	Premium      int64  `json:"premium"`
	Pertalite    int64  `json:"pertalite"`
}
