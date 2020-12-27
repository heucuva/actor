# goactor

A Pure Go actor/tick/manager library, vaguely similar to Unreal Engine 4's system

# Usage

## Creating an actor

Either spawn an actor (using the `actor.SpawnActor()` call) or build your own via standard Go mechanisms.

### `actor.SpawnActor()`

If you use the `actor.SpawnActor()` call, then you get a series of function callbacks for free. These callbacks are optionally-defined and happen in a specific order:

1. `PostSpawnInitialize`
2. `ExecuteConstruction`
3. `OnConstruction`
4. `PostActorConstruction`
5. `PreInitializeComponents`
6. `InitializeComponents`
7. `PostInitializeComponents`
8. `OnActorSpawned`

If an actor was created with the `actor.DeferredSpawnActor()` option, then the callback sequence pauses after `PostSpawnInitialize` and will not continue until after a call to `actor.FinishSpawningActor()` is made (and making sure to pass in the same SpawnActorOptions)

## Getting Ticks

There are a few ways to get ticks on an non-zero time interval for the actors. The easiest way is to call the `AddActor()` function on the default global manager, found at `actor.GetManager()` and pass along the `actor.TickInterval` option with your desired non-zero time interval.  This manager is configured to run at application startup and runs in on the background context.

If you want a set of actors to tick on every _frame_, then you must pass in the `actor.TickEveryFrame` option.  This sets the actor to not automatically tick on a given interval, but instead, it will have its `Tick` callback called after every call to the manager's `TickFrame()` function. **NOTE**: you must call `TickFrame()` yourself for this to work as expected.

When you call the `AddActor()` function, the actor will receive this optional callback before any `Tick` callbacks will fire:

9. `BeginPlay`

If you add an actor to fire on a specific interval from within the scope of an existing tick event, it will not get a `Tick` callback until the next cycle of the interval, which may be significantly more or less than the expected interval duration. Be sure to consider the `deltaTime` value that is passed along with the `Tick` callback.

## Making Your Own Manager Instances

Sure, why not?  Have as many as you'd like.  The default-constructed global one is probably fine for most tasks, though.

## Shutting Down Actors

When you are done with a singular actor, simply ask the Manager it's registered to remove it via a call to `RemoveActor()`.  This will trigger the following optional callback(s):

10. `EndPlay`

The `EndPlay` callback will include the reason for the callback, provided as an `error` value.

If you are wanting to shut down a manager and trigger `EndPlay` on all its actors, simply ask the manager to do so by calling its `Stop()` function.  Once you do this, however, the manager will no longer be valid for use and cannot be reset.
