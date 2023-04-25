# Meduce

Library for processing huge amounts of data on your PC by using MapReduce paradigm.

### Usage
The paradigm is pretty simple. You only need two (or three) functions to process
all of your data.

#### 1. `func Mapper(key KeyIn, value ValueIn, emit Emitter[KeyOut, ValueOut])`
This function maps data that you supplied to key-value pairs. 
For each piece of data you can emit as many key-value pairs 
as you want by calling `emit` function. 
The function should not have any side effects.

#### 2. `func Reducer(key KeyOut, values []ValueOut) ValueOut`
This function reduces all values that were mapped to the same key. 
It can also be called as combiner function, so it should be idempotent
and have no side effects.

#### 3. `func Finalizer(key KeyOut, value ValueOut) ValueOut` (optional)
This function is called after all data was processed and reduced.

