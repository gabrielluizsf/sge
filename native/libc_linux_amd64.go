package native

// The libc functions the thread hooks and the loader need. glibc 2.34+
// carries dlopen and pthread in libc.so.6 itself.

//go:cgo_import_dynamic _ _ "libc.so.6"
//go:cgo_import_dynamic native_malloc malloc "libc.so.6"
//go:cgo_import_dynamic native_free free "libc.so.6"
//go:cgo_import_dynamic native_abort abort "libc.so.6"
//go:cgo_import_dynamic native_setenv setenv "libc.so.6"
//go:cgo_import_dynamic native_unsetenv unsetenv "libc.so.6"
//go:cgo_import_dynamic native_sigfillset sigfillset "libc.so.6"
//go:cgo_import_dynamic native_pthread_sigmask pthread_sigmask "libc.so.6"
//go:cgo_import_dynamic native_pthread_self pthread_self "libc.so.6"
//go:cgo_import_dynamic native_pthread_create pthread_create "libc.so.6"
//go:cgo_import_dynamic native_pthread_attr_init pthread_attr_init "libc.so.6"
//go:cgo_import_dynamic native_pthread_attr_destroy pthread_attr_destroy "libc.so.6"
//go:cgo_import_dynamic native_pthread_attr_getstack pthread_attr_getstack "libc.so.6"
//go:cgo_import_dynamic native_pthread_attr_getstacksize pthread_attr_getstacksize "libc.so.6"
//go:cgo_import_dynamic native_pthread_attr_setdetachstate pthread_attr_setdetachstate "libc.so.6"
//go:cgo_import_dynamic native_pthread_getattr_np pthread_getattr_np "libc.so.6"
//go:cgo_import_dynamic native_dlopen dlopen "libc.so.6"
//go:cgo_import_dynamic native_dlsym dlsym "libc.so.6"
//go:cgo_import_dynamic native_dlclose dlclose "libc.so.6"
//go:cgo_import_dynamic native_dlerror dlerror "libc.so.6"

// libc holds the addresses of the PLT stubs, filled by libc_linux_amd64.s.
type libcTable struct {
	malloc, free, abort, setenv, unsetenv                     uintptr
	sigfillset, sigmask, self, create                         uintptr
	attrInit, attrDestroy, attrGetstack, attrGetstacksize     uintptr
	attrSetdetach, getattrNP, dlopen, dlsym, dlclose, dlerror uintptr
}

var libc libcTable
