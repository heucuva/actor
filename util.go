package actor

import "time"

// PostSpawnInitialize calls an actor's PostSpawnInitialize() function, if it has one
func PostSpawnInitialize(a Actor) error {
	if t, ok := a.(PostSpawnInitializeIntf); ok {
		return t.PostSpawnInitialize()
	}

	return nil
}

// ExecuteConstruction calls an actor's ExecuteConstruction() function, if it has one
func ExecuteConstruction(a Actor) error {
	if t, ok := a.(ExecuteConstructionIntf); ok {
		return t.ExecuteConstruction()
	}

	return nil
}

// OnConstruction calls an actor's OnConstruction() function, if it has one
func OnConstruction(a Actor) error {
	if t, ok := a.(OnConstructionIntf); ok {
		return t.OnConstruction()
	}

	return nil
}

// PostActorConstruction calls an actor's PostActorConstruction() function, if it has one
func PostActorConstruction(a Actor) error {
	if t, ok := a.(PostActorConstructionIntf); ok {
		return t.PostActorConstruction()
	}

	return nil
}

// PreInitializeComponents calls an actor's PreInitializeComponents() function, if it has one
func PreInitializeComponents(a Actor) error {
	if t, ok := a.(PreInitializeComponentsIntf); ok {
		return t.PreInitializeComponents()
	}

	return nil
}

// InitializeComponents calls an actor's InitializeComponents() function, if it has one
func InitializeComponents(a Actor) error {
	if t, ok := a.(InitializeComponentsIntf); ok {
		return t.InitializeComponents()
	}

	return nil
}

// PostInitializeComponents calls an actor's PostInitializeComponents() function, if it has one
func PostInitializeComponents(a Actor) error {
	if t, ok := a.(PostInitializeComponentsIntf); ok {
		return t.PostInitializeComponents()
	}

	return nil
}

// OnActorSpawned calls an actor's OnActorSpawned() function, if it has one
func OnActorSpawned(a Actor) error {
	if t, ok := a.(OnActorSpawnedIntf); ok {
		return t.OnActorSpawned()
	}

	return nil
}

// BeginPlay calls an actor's BeginPlay() function, if it has one
func BeginPlay(a Actor) error {
	if t, ok := a.(BeginPlayIntf); ok {
		return t.BeginPlay()
	}

	return nil
}

// Tick calls an actor's Tick() function, if it has one
func Tick(a Actor, deltaTime time.Duration) error {
	if t, ok := a.(TickIntf); ok {
		return t.Tick(deltaTime)
	}

	return nil
}

// WantTick calls an actor's WantTick() function, if it has one
func WantTick(a Actor) (bool, error) {
	if t, ok := a.(WantTickIntf); ok {
		return t.WantTick()
	}

	return true, nil
}

// EndPlay calls an actor's EndPlay() function, if it has one
func EndPlay(a Actor, endPlayReason error) error {
	if t, ok := a.(EndPlayIntf); ok {
		return t.EndPlay(endPlayReason)
	}

	return nil
}

// BeginDestroy calls an actor's BeginDestroy() function, if it has one
func BeginDestroy(a Actor) error {
	if t, ok := a.(BeginDestroyIntf); ok {
		return t.BeginDestroy()
	}

	return nil
}

// FinishDestroy calls an actor's FinishDestroy() function, if it has one
func FinishDestroy(a Actor) error {
	if t, ok := a.(FinishDestroyIntf); ok {
		return t.FinishDestroy()
	}

	return nil
}
