//go:build gots_overlay

// Overlay for the stack intrinsics: JavaScript has no program counters or
// stack pointers to give.

package sys

func GetCallerPC() uintptr { return 0 }

func GetCallerSP() uintptr { return 0 }

func GetClosurePtr() uintptr { return 0 }
