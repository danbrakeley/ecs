package test

import (
	"testing"
	"time"

	"github.com/danbrakeley/ecs"
)

//
// baseSys: no tick, no receiver

type baseSys struct {
	UpdatedEntity *ecs.Entity
	DroppedEntity *ecs.Entity
}

func (b *baseSys) GetName() string {
	return "baseSys"
}
func (b *baseSys) UpdateEntity(e *ecs.Entity) {
	b.UpdatedEntity = e
}
func (b *baseSys) DropEntity(e *ecs.Entity) {
	b.DroppedEntity = e
}

func createMgrAndBaseSys() (*ecs.Manager, *baseSys) {
	mgr := ecs.NewManager()
	b := &baseSys{}
	mgr.RegisterSystem(b)
	return mgr, b
}

//
// tickSys: no receiver

type tickSys struct {
	baseSys
	WasTicked bool
}

func (t *tickSys) GetName() string {
	return "tickSys"
}
func (t *tickSys) Tick(dt time.Duration) {
	t.WasTicked = true
}

func createMgrAndTickSys() (*ecs.Manager, *tickSys) {
	mgr := ecs.NewManager()
	t := &tickSys{}
	mgr.RegisterSystem(t)
	return mgr, t
}

//
// receiverSys: no tick

type receiverSys struct {
	baseSys
	Queue []*ecs.Entity
}

func (r *receiverSys) GetName() string {
	return "receiverSys"
}
func (r *receiverSys) ClearQueue() {
	r.Queue = []*ecs.Entity{}
}
func (r *receiverSys) HandleBroadcast(e *ecs.Entity) {
	r.Queue = append(r.Queue, e)
}

func createMgrAndReceiveSys() (*ecs.Manager, *receiverSys) {
	mgr := ecs.NewManager()
	r := &receiverSys{}
	mgr.RegisterSystem(r)
	return mgr, r
}

//
// Base Tests

func Test_ECS_SystemsUpdatedOnAddComponent(t *testing.T) {
	mgr, s := createMgrAndBaseSys()
	c := &EmptyCom{}
	e := mgr.NewEntity()
	e.AddComponent(c)

	if s.UpdatedEntity != e {
		t.Errorf("System did not update entity when component was added")
	}
}

func Test_ECS_SystemsAlertedToDrop(t *testing.T) {
	mgr, s := createMgrAndBaseSys()
	c := &EmptyCom{}
	e := mgr.NewEntity()
	e.AddComponent(c)

	e.Drop()
	if s.DroppedEntity != e {
		t.Errorf("System was not told to drop entity when entity called Drop")
	}
}

//
// Tick Tests

func Test_ECS_RegisteredSystemGetsTick(t *testing.T) {
	mgr, tickSys := createMgrAndTickSys()
	tickSys.WasTicked = false
	mgr.Tick(time.Second)
	if tickSys.WasTicked != true {
		t.Errorf("Registered system did not get a Tick")
	}
}

//
// Receiver Tests

func Test_ECS_RegisteredSystemGetsBroadcast(t *testing.T) {
	mgr, recSys := createMgrAndReceiveSys()
	e := mgr.NewEntity()
	sentID := e.GetID()
	mgr.Broadcast(e)
	if len(recSys.Queue) != 1 {
		t.Errorf("Registered system did not receive a broadcast")
		return
	}
	receivedID := recSys.Queue[0].GetID()
	if sentID != receivedID {
		t.Errorf("Broadcasted id %v, but heard id %v", sentID, receivedID)
		return
	}
}
