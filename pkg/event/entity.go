package event

type IDer interface {
	ID() string
}

type EntityNamer interface {
	EntityName() string
}

type Entity struct {
	id   string
	name string
}

func (e Entity) ID() string         { return e.id }
func (e Entity) EntityName() string { return e.name }
func (e Entity) Equals(o IDer) bool { return e.id == o.ID() }

func (e *Entity) setID(id string)     { e.id = id }
func (e *Entity) setName(name string) { e.name = name }
