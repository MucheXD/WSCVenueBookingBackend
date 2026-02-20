package systemPermission

func Map(units ...SystemPermission) uint64 {
	var permMap uint64 = 0
	for _, unit := range units {
		permMap |= uint64(unit)
	}
	return permMap
}

func Check(permMap uint64, required SystemPermission) bool {
	return (permMap & uint64(required)) != 0
}

func Satisfy(permMap uint64, required ...SystemPermission) bool {
	for _, req := range required {
		if (permMap & uint64(req)) != 0 {
			return true
		}
	}
	return false
}
