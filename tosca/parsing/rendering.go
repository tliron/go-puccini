package parsing

//
// Renderable
//

type Renderable interface {
	Render()
}

// From [Renderable] interface
func Render(entityPtr EntityPtr) bool {
	if renderable, ok := entityPtr.(Renderable); ok {
		renderable.Render()
		return true
	} else {
		return false
	}
}
