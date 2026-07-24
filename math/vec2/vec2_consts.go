package vec2

func Up() Vec2 {
	return New(0.0, 1.0)
}

func Down() Vec2 {
	return New(0.0, -1.0)
}

func Right() Vec2 {
	return New(1.0, 0.0)
}

func Left() Vec2 {
	return New(-1.0, 0.0)
}

func Zero() Vec2 {
	return New(0.0, 0.0)
}

func One() Vec2 {
	return New(1.0, 1.0)
}

func NegOne() Vec2 {
	return New(-1.0, -1.0)
}
