package routerStruct

func FindMainDirection(route Route) int {
	for i, direction := range route.Directions {
		if direction.Current == 1 {
			return i
		}
	}
	return 0
}

func CheckAllDead(route Route) bool{
	for _, direction := range route.Directions{
		if direction.Alive == 1{
			return false
		}
	}
	return true
}

func GetNextDirection(route Route) int{
	return 0
}

func ResetDirections(route *Route) {
	for _, direction := range route.Directions {
		direction.Alive = 1
	}
}

