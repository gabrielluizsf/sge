//go:build !cgo

#include "textflag.h"

// C-ABI entry points the runtime calls instead of runtime/cgo's C code.
// Each saves the System V callee-saved registers, then calls the Go
// (ABI0, nosplit) half with its C arguments as Go stack arguments.

#define SAVE MOVQ BX, 16(SP); MOVQ R12, 24(SP); MOVQ R13, 32(SP); MOVQ R14, 40(SP); MOVQ R15, 48(SP)
#define RESTORE MOVQ 16(SP), BX; MOVQ 24(SP), R12; MOVQ 32(SP), R13; MOVQ 40(SP), R14; MOVQ 48(SP), R15

// void x_cgo_init(G *g, void (*setg)(void*))
TEXT cgoInitC<>(SB), NOSPLIT, $56-0
	SAVE
	MOVQ DI, 0(SP)
	MOVQ SI, 8(SP)
	CALL ·cgoInit(SB)
	RESTORE
	RET

// void x_cgo_thread_start(ThreadStart *ts)
TEXT threadStartC<>(SB), NOSPLIT, $56-0
	SAVE
	MOVQ DI, 0(SP)
	CALL ·threadStart(SB)
	RESTORE
	RET

// void *threadentry(void *ts) — the pthread start routine.
TEXT threadEntryC<>(SB), NOSPLIT, $56-0
	SAVE
	MOVQ DI, 0(SP)
	CALL ·threadEntry(SB)
	MOVQ 8(SP), AX
	CALL AX // mstart; returns when the M exits
	RESTORE
	XORL AX, AX
	RET

// x_cgo_setenv(char *arg[2]) -> setenv(arg[0], arg[1], 1)
TEXT setenvC<>(SB), NOSPLIT|NOFRAME, $0-0
	MOVQ 8(DI), SI
	MOVQ 0(DI), DI
	MOVL $1, DX
	JMP native_setenv(SB)

// x_cgo_unsetenv(char *arg[1]) -> unsetenv(arg[0])
TEXT unsetenvC<>(SB), NOSPLIT|NOFRAME, $0-0
	MOVQ 0(DI), DI
	JMP native_unsetenv(SB)

// Nothing waits on runtime init: C never calls back into Go here.
TEXT noopC<>(SB), NOSPLIT|NOFRAME, $0-0
	RET

DATA noopFuncval<>+0(SB)/8, $noopC<>(SB)
GLOBL noopFuncval<>(SB), RODATA, $8

// The runtime's hook variables. Ours carry content, so the linker keeps
// them over the runtime's zero (BSS) definitions.
DATA _cgo_init+0(SB)/8, $cgoInitC<>(SB)
GLOBL _cgo_init(SB), NOPTR, $8
DATA _cgo_thread_start+0(SB)/8, $threadStartC<>(SB)
GLOBL _cgo_thread_start(SB), NOPTR, $8
DATA _cgo_notify_runtime_init_done+0(SB)/8, $noopC<>(SB)
GLOBL _cgo_notify_runtime_init_done(SB), NOPTR, $8
DATA _cgo_pthread_key_created+0(SB)/8, $·keyCreated(SB)
GLOBL _cgo_pthread_key_created(SB), NOPTR, $8
DATA runtime·_cgo_setenv+0(SB)/8, $setenvC<>(SB)
GLOBL runtime·_cgo_setenv(SB), NOPTR, $8
DATA runtime·_cgo_unsetenv+0(SB)/8, $unsetenvC<>(SB)
GLOBL runtime·_cgo_unsetenv(SB), NOPTR, $8
DATA runtime·set_crosscall2+0(SB)/8, $noopFuncval<>(SB)
GLOBL runtime·set_crosscall2(SB), NOPTR, $8
DATA runtime·iscgo+0(SB)/1, $1
GLOBL runtime·iscgo(SB), NOPTR, $1

DATA ·threadEntryPC+0(SB)/8, $threadEntryC<>(SB)
GLOBL ·threadEntryPC(SB), NOPTR, $8
