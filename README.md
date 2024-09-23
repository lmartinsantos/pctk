# Point & Click Toolkit

Welcome to the pctk repository: the toolkit to make Point&Click adventure 
games. 

This is still work in progress. Please come back later if you want to find 
something useful.

## TODO

Related to objects:

- Move obects under Room (missing onInit / onExit)
- Complete `processControlInputs`
  - Complete hovering over actors (filtering ego)
  - Complex actions (use X on Y, Give X to Y)
- Simplify mouse collision/hove with a `GetTarget` function (should actor be an object?)
- Ego Commands (for objet's scripts)
- Update object commands (to set classes, change state etc,.)
- Rendering room objects (including actors) using z-index (updating `Position` to add Z coord)
