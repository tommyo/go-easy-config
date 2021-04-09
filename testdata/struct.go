package testdata

type SmallExample struct {
	Foo  string  `json:"foo"`
	Bar  float64 `json:"bar"`
	Baz  uint    `json:"baz"`
	None bool    `json:"none" usage:"A boolean flag"`
}

type FullExample struct {
	Simple    string       `json:"simple"`
	Described bool         `json:"described" usage:"A short description"`
	Skipped   string       `json:"-"`
	Deep      SmallExample `json:"deep"`
}
