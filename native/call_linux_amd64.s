#include "textflag.h"

// void callFrameC(Frame *f) — C ABI, reached through runtime.asmcgocall.
// Loads 6 integer and 8 SSE registers, copies ints[6:16] to the stack,
// calls f.fn and stores RAX and XMM0 back into the frame.
TEXT callFrameC<>(SB), NOSPLIT|NOFRAME, $0-0
	PUSHQ BP
	MOVQ SP, BP
	PUSHQ BX
	MOVQ DI, BX
	SUBQ $80, SP
	ANDQ $~15, SP
	MOVQ 56(BX), AX
	MOVQ AX, 0(SP)
	MOVQ 64(BX), AX
	MOVQ AX, 8(SP)
	MOVQ 72(BX), AX
	MOVQ AX, 16(SP)
	MOVQ 80(BX), AX
	MOVQ AX, 24(SP)
	MOVQ 88(BX), AX
	MOVQ AX, 32(SP)
	MOVQ 96(BX), AX
	MOVQ AX, 40(SP)
	MOVQ 104(BX), AX
	MOVQ AX, 48(SP)
	MOVQ 112(BX), AX
	MOVQ AX, 56(SP)
	MOVQ 120(BX), AX
	MOVQ AX, 64(SP)
	MOVQ 128(BX), AX
	MOVQ AX, 72(SP)
	MOVSD 136(BX), X0
	MOVSD 144(BX), X1
	MOVSD 152(BX), X2
	MOVSD 160(BX), X3
	MOVSD 168(BX), X4
	MOVSD 176(BX), X5
	MOVSD 184(BX), X6
	MOVSD 192(BX), X7
	MOVQ 8(BX), DI
	MOVQ 16(BX), SI
	MOVQ 24(BX), DX
	MOVQ 32(BX), CX
	MOVQ 40(BX), R8
	MOVQ 48(BX), R9
	MOVQ 0(BX), R10
	MOVL $8, AX // upper bound of vector registers, for variadic callees
	CALL R10
	MOVQ AX, 200(BX)
	MOVSD X0, 208(BX)
	LEAQ -8(BP), SP
	POPQ BX
	POPQ BP
	RET

DATA ·callFramePC+0(SB)/8, $callFrameC<>(SB)
GLOBL ·callFramePC(SB), NOPTR, $8
