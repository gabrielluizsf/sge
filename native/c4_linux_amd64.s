#include "textflag.h"

// func hookCall(f *hookFrame)
// BX and R12 are callee-saved in the System V ABI, so they survive the call.
TEXT ·hookCall(SB), NOSPLIT, $0-8
	MOVQ f+0(FP), BX
	MOVQ SP, R12
	ANDQ $~15, SP
	MOVQ 8(BX), DI
	MOVQ 16(BX), SI
	MOVQ 24(BX), DX
	MOVQ 32(BX), CX
	MOVQ 0(BX), R10
	XORL AX, AX
	CALL R10
	MOVQ R12, SP
	MOVQ AX, 40(BX)
	RET
