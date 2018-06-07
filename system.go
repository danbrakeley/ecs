package ecs

import "time"

// System is the interface that a registered/managed system must expose.
type System interface {
	GetName() string        // return a unique name for this system (ie "FooSystem")
	UpdateEntity(e *Entity) // this entity's component makeup has changed
	DropEntity(e *Entity)   // this entity is about to be deleted (so remove it from your systems)
}

// Ticker is the interface for systems that can recieve Ticks
type Ticker interface {
	Tick(dt time.Duration) // called every frame; dt is time since last frame
}

// Receiver is the interface for systems that want to hear broadcasts
type Receiver interface {
	HandleBroadcast(e *Entity)
}
