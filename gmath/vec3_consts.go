package gmath

// --- Constructors ---

func Zero3() Vector3 {
	return Vec3(0, 0, 0)
}

func One3() Vector3 {
	return Vec3(1, 1, 1)
}

func Up3() Vector3 {
	return Vec3(0, 1, 0)
}

func Down3() Vector3 {
	return Vec3(0, -1, 0)
}

func Right3() Vector3 {
	return Vec3(1, 0, 0)
}

func Left3() Vector3 {
	return Vec3(-1, 0, 0)
}

func Forward3() Vector3 {
	return Vec3(0, 0, 1)
}

func Backward3() Vector3 {
	return Vec3(0, 0, -1)
}

