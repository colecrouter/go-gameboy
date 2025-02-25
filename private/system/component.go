package system

type MComponent interface {
	Component
	MClock()
}

type TComponent interface {
	Component
	TClock()
}

type Component interface {
	Run()
}
