// Package animal gives animals made without a skeleton one, and moves any
// animal's skeleton by where its feet must be: bones laid out along a
// model, skin weights from the nearest bone, and clips, standing, walking,
// grazing, biting and falling, written as keyframes onto the bones of the
// model, for package bake to pose.
//
// A rigged model faces +Z, stands on Y=0 and is centred on X=0 and Z=0, as
// package bake expects; sides are along X.
package animal
