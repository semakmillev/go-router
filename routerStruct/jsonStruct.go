package routerStruct

type JsonRouter struct {
	Routes []Route `json:"routes"`
}

type Direction struct {
	URL     string `json:"url"`
	Current int
	Alive int
}

type Route struct {
	Name       string      `json:"name"`
	Port       string      `json:"port"`
	Directions []Direction `json:"directions"`
}
