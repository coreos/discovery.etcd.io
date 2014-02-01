all: \
	mach_kernel_10_8_0 \
	mach_kernel_10_8_0_fix \
	mach_kernel_10_8_0_i386 \
	mach_kernel_10_8_0_i386_fix \
	mach_kernel_11_4_2 \
	mach_kernel_11_4_2_fix \
	mach_kernel_12_4_0 \
	mach_kernel_12_4_0_fix \

%_i386: %_i386.s
	gcc -m32 -o $*_i386 $*_i386.s

%_i386_fix: %_i386_fix.s
	gcc -m32 -o $*_i386_fix $*_i386_fix.s

%: %.s
	gcc -o $* $*.s

