
#include "textflag.h"

// One JMP per imported symbol; the linker routes it through the PLT.
TEXT malloc_jmp<>(SB), NOSPLIT|NOFRAME, $0-0
	JMP native_malloc(SB)

TEXT free_jmp<>(SB), NOSPLIT|NOFRAME, $0-0
	JMP native_free(SB)

TEXT abort_jmp<>(SB), NOSPLIT|NOFRAME, $0-0
	JMP native_abort(SB)

TEXT setenv_jmp<>(SB), NOSPLIT|NOFRAME, $0-0
	JMP native_setenv(SB)

TEXT unsetenv_jmp<>(SB), NOSPLIT|NOFRAME, $0-0
	JMP native_unsetenv(SB)

TEXT sigfillset_jmp<>(SB), NOSPLIT|NOFRAME, $0-0
	JMP native_sigfillset(SB)

TEXT pthread_sigmask_jmp<>(SB), NOSPLIT|NOFRAME, $0-0
	JMP native_pthread_sigmask(SB)

TEXT pthread_self_jmp<>(SB), NOSPLIT|NOFRAME, $0-0
	JMP native_pthread_self(SB)

TEXT pthread_create_jmp<>(SB), NOSPLIT|NOFRAME, $0-0
	JMP native_pthread_create(SB)

TEXT pthread_attr_init_jmp<>(SB), NOSPLIT|NOFRAME, $0-0
	JMP native_pthread_attr_init(SB)

TEXT pthread_attr_destroy_jmp<>(SB), NOSPLIT|NOFRAME, $0-0
	JMP native_pthread_attr_destroy(SB)

TEXT pthread_attr_getstack_jmp<>(SB), NOSPLIT|NOFRAME, $0-0
	JMP native_pthread_attr_getstack(SB)

TEXT pthread_attr_getstacksize_jmp<>(SB), NOSPLIT|NOFRAME, $0-0
	JMP native_pthread_attr_getstacksize(SB)

TEXT pthread_attr_setdetachstate_jmp<>(SB), NOSPLIT|NOFRAME, $0-0
	JMP native_pthread_attr_setdetachstate(SB)

TEXT pthread_getattr_np_jmp<>(SB), NOSPLIT|NOFRAME, $0-0
	JMP native_pthread_getattr_np(SB)

TEXT dlopen_jmp<>(SB), NOSPLIT|NOFRAME, $0-0
	JMP native_dlopen(SB)

TEXT dlsym_jmp<>(SB), NOSPLIT|NOFRAME, $0-0
	JMP native_dlsym(SB)

TEXT dlclose_jmp<>(SB), NOSPLIT|NOFRAME, $0-0
	JMP native_dlclose(SB)

TEXT dlerror_jmp<>(SB), NOSPLIT|NOFRAME, $0-0
	JMP native_dlerror(SB)

DATA ·libc+0(SB)/8, $malloc_jmp<>(SB)
DATA ·libc+8(SB)/8, $free_jmp<>(SB)
DATA ·libc+16(SB)/8, $abort_jmp<>(SB)
DATA ·libc+24(SB)/8, $setenv_jmp<>(SB)
DATA ·libc+32(SB)/8, $unsetenv_jmp<>(SB)
DATA ·libc+40(SB)/8, $sigfillset_jmp<>(SB)
DATA ·libc+48(SB)/8, $pthread_sigmask_jmp<>(SB)
DATA ·libc+56(SB)/8, $pthread_self_jmp<>(SB)
DATA ·libc+64(SB)/8, $pthread_create_jmp<>(SB)
DATA ·libc+72(SB)/8, $pthread_attr_init_jmp<>(SB)
DATA ·libc+80(SB)/8, $pthread_attr_destroy_jmp<>(SB)
DATA ·libc+88(SB)/8, $pthread_attr_getstack_jmp<>(SB)
DATA ·libc+96(SB)/8, $pthread_attr_getstacksize_jmp<>(SB)
DATA ·libc+104(SB)/8, $pthread_attr_setdetachstate_jmp<>(SB)
DATA ·libc+112(SB)/8, $pthread_getattr_np_jmp<>(SB)
DATA ·libc+120(SB)/8, $dlopen_jmp<>(SB)
DATA ·libc+128(SB)/8, $dlsym_jmp<>(SB)
DATA ·libc+136(SB)/8, $dlclose_jmp<>(SB)
DATA ·libc+144(SB)/8, $dlerror_jmp<>(SB)
GLOBL ·libc(SB), NOPTR, $152
