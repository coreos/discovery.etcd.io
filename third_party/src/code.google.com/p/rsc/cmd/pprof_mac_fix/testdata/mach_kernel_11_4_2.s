_version:
	.ascii "Darwin Kernel Version 11.4.2: Thu Aug 23 16:25:48 PDT 2012; root:xnu-1699.32.7~1/RELEASE_X86_64\0"

.globl _main
_main:
	ret

.globl _psignal_internal
_psignal_internal:
	ret

.globl _task_vtimer_clear
_task_vtimer_clear:
	ret

.globl _task_vtimer_set
_task_vtimer_set:
	ret

.globl _current_thread
_current_thread:
// 0xffffff80002be1c0 <current_thread+0>:	push   %rbp
	.byte 0x55;
// 0xffffff80002be1c1 <current_thread+1>:	mov    %rsp,%rbp
	.byte 0x48; .byte 0x89; .byte 0xe5;
// 0xffffff80002be1c4 <current_thread+4>:	mov    %gs:0x8,%rax
	.byte 0x65; .byte 0x48; .byte 0x8b; .byte 0x04; .byte 0x25; .byte 0x08; .byte 0x00; .byte 0x00; .byte 0x00;
// 0xffffff80002be1cd <current_thread+13>:	pop    %rbp
	.byte 0x5d;
// 0xffffff80002be1ce <current_thread+14>:	retq   
	.byte 0xc3;
// 0xffffff80002be1cf <current_thread+15>:	nop    
	.byte 0x90;

.globl _bsd_ast
_bsd_ast:
// 0xffffff8000553410 <bsd_ast+0>:	push   %rbp
	.byte 0x55;
// 0xffffff8000553411 <bsd_ast+1>:	mov    %rsp,%rbp
	.byte 0x48; .byte 0x89; .byte 0xe5;
// 0xffffff8000553414 <bsd_ast+4>:	push   %r15
	.byte 0x41; .byte 0x57;
// 0xffffff8000553416 <bsd_ast+6>:	push   %r14
	.byte 0x41; .byte 0x56;
// 0xffffff8000553418 <bsd_ast+8>:	push   %rbx
	.byte 0x53;
// 0xffffff8000553419 <bsd_ast+9>:	sub    $0x18,%rsp
	.byte 0x48; .byte 0x83; .byte 0xec; .byte 0x18;
// 0xffffff800055341d <bsd_ast+13>:	mov    %rdi,%rbx
	.byte 0x48; .byte 0x89; .byte 0xfb;
// 0xffffff8000553420 <bsd_ast+16>:	callq  0xffffff80005e1db0 <current_proc>
	.byte 0xe8; .byte 0x8b; .byte 0xe9; .byte 0x08; .byte 0x00;
// 0xffffff8000553425 <bsd_ast+21>:	mov    %rax,%r14
	.byte 0x49; .byte 0x89; .byte 0xc6;
// 0xffffff8000553428 <bsd_ast+24>:	mov    %rbx,%rdi
	.byte 0x48; .byte 0x89; .byte 0xdf;
// 0xffffff800055342b <bsd_ast+27>:	callq  0xffffff8000245340 <get_bsdthread_info>
	.byte 0xe8; .byte 0x10; .byte 0x1f; .byte 0xcf; .byte 0xff;
// 0xffffff8000553430 <bsd_ast+32>:	mov    %rax,%rbx
	.byte 0x48; .byte 0x89; .byte 0xc3;
// 0xffffff8000553433 <bsd_ast+35>:	test   %r14,%r14
	.byte 0x4d; .byte 0x85; .byte 0xf6;
// 0xffffff8000553436 <bsd_ast+38>:	je     0xffffff8000553779 <bsd_ast+873>
	.byte 0x0f; .byte 0x84; .byte 0x3d; .byte 0x03; .byte 0x00; .byte 0x00;
// 0xffffff800055343c <bsd_ast+44>:	mov    0x158(%r14),%eax
	.byte 0x41; .byte 0x8b; .byte 0x86; .byte 0x58; .byte 0x01; .byte 0x00; .byte 0x00;
// 0xffffff8000553443 <bsd_ast+51>:	test   $0x80,%ah
	.byte 0xf6; .byte 0xc4; .byte 0x80;
// 0xffffff8000553446 <bsd_ast+54>:	je     0xffffff8000553474 <bsd_ast+100>
	.byte 0x74; .byte 0x2c;
// 0xffffff8000553448 <bsd_ast+56>:	test   $0x20,%al
	.byte 0xa8; .byte 0x20;
// 0xffffff800055344a <bsd_ast+58>:	je     0xffffff8000553474 <bsd_ast+100>
	.byte 0x74; .byte 0x28;
// 0xffffff800055344c <bsd_ast+60>:	lea    0x158(%r14),%r15
	.byte 0x4d; .byte 0x8d; .byte 0xbe; .byte 0x58; .byte 0x01; .byte 0x00; .byte 0x00;
// 0xffffff8000553453 <bsd_ast+67>:	callq  0xffffff80002c0b70 <get_useraddr>
	.byte 0xe8; .byte 0x18; .byte 0xd7; .byte 0xd6; .byte 0xff;
// 0xffffff8000553458 <bsd_ast+72>:	mov    %eax,%esi
	.byte 0x89; .byte 0xc6;
// 0xffffff800055345a <bsd_ast+74>:	mov    $0x1,%edx
	.byte 0xba; .byte 0x01; .byte 0x00; .byte 0x00; .byte 0x00;
// 0xffffff800055345f <bsd_ast+79>:	mov    %r14,%rdi
	.byte 0x4c; .byte 0x89; .byte 0xf7;
// 0xffffff8000553462 <bsd_ast+82>:	callq  0xffffff8000561220 <addupc_task>
	.byte 0xe8; .byte 0xb9; .byte 0xdd; .byte 0x00; .byte 0x00;
// 0xffffff8000553467 <bsd_ast+87>:	mov    $0xffff7fff,%edi
	.byte 0xbf; .byte 0xff; .byte 0x7f; .byte 0xff; .byte 0xff;
// 0xffffff800055346c <bsd_ast+92>:	mov    %r15,%rsi
	.byte 0x4c; .byte 0x89; .byte 0xfe;
// 0xffffff800055346f <bsd_ast+95>:	callq  0xffffff80005e2af0 <OSBitAndAtomic>
	.byte 0xe8; .byte 0x7c; .byte 0xf6; .byte 0x08; .byte 0x00;
// 0xffffff8000553474 <bsd_ast+100>:	cmpq   $0x0,0x1c0(%r14)
	.byte 0x49; .byte 0x83; .byte 0xbe; .byte 0xc0; .byte 0x01; .byte 0x00; .byte 0x00; .byte 0x00;
// 0xffffff800055347c <bsd_ast+108>:	jne    0xffffff8000553488 <bsd_ast+120>
	.byte 0x75; .byte 0x0a;
// 0xffffff800055347e <bsd_ast+110>:	cmpl   $0x0,0x1c8(%r14)
	.byte 0x41; .byte 0x83; .byte 0xbe; .byte 0xc8; .byte 0x01; .byte 0x00; .byte 0x00; .byte 0x00;
// 0xffffff8000553486 <bsd_ast+118>:	je     0xffffff80005534f6 <bsd_ast+230>
	.byte 0x74; .byte 0x6e;
// 0xffffff8000553488 <bsd_ast+120>:	mov    0x18(%r14),%rdi
	.byte 0x49; .byte 0x8b; .byte 0x7e; .byte 0x18;
// 0xffffff800055348c <bsd_ast+124>:	mov    $0x1,%esi
	.byte 0xbe; .byte 0x01; .byte 0x00; .byte 0x00; .byte 0x00;
// 0xffffff8000553491 <bsd_ast+129>:	lea    -0x1c(%rbp),%rdx
	.byte 0x48; .byte 0x8d; .byte 0x55; .byte 0xe4;
// 0xffffff8000553495 <bsd_ast+133>:	callq  0xffffff8000236330 <task_vtimer_update>
	.byte 0xe8; .byte 0x96; .byte 0x2e; .byte 0xce; .byte 0xff;
// 0xffffff800055349a <bsd_ast+138>:	mov    -0x1c(%rbp),%edx
	.byte 0x8b; .byte 0x55; .byte 0xe4;
// 0xffffff800055349d <bsd_ast+141>:	lea    0x1b0(%r14),%rsi
	.byte 0x49; .byte 0x8d; .byte 0xb6; .byte 0xb0; .byte 0x01; .byte 0x00; .byte 0x00;
// 0xffffff80005534a4 <bsd_ast+148>:	mov    %r14,%rdi
	.byte 0x4c; .byte 0x89; .byte 0xf7;
// 0xffffff80005534a7 <bsd_ast+151>:	callq  0xffffff800055c640 <itimerdecr>
	.byte 0xe8; .byte 0x94; .byte 0x91; .byte 0x00; .byte 0x00;
// 0xffffff80005534ac <bsd_ast+156>:	test   %eax,%eax
	.byte 0x85; .byte 0xc0;
// 0xffffff80005534ae <bsd_ast+158>:	jne    0xffffff80005534f6 <bsd_ast+230>
	.byte 0x75; .byte 0x46;
// 0xffffff80005534b0 <bsd_ast+160>:	cmpq   $0x0,0x1c0(%r14)
	.byte 0x49; .byte 0x83; .byte 0xbe; .byte 0xc0; .byte 0x01; .byte 0x00; .byte 0x00; .byte 0x00;
// 0xffffff80005534b8 <bsd_ast+168>:	jne    0xffffff80005534c4 <bsd_ast+180>
	.byte 0x75; .byte 0x0a;
// 0xffffff80005534ba <bsd_ast+170>:	cmpl   $0x0,0x1c8(%r14)
	.byte 0x41; .byte 0x83; .byte 0xbe; .byte 0xc8; .byte 0x01; .byte 0x00; .byte 0x00; .byte 0x00;
// 0xffffff80005534c2 <bsd_ast+178>:	je     0xffffff80005534d4 <bsd_ast+196>
	.byte 0x74; .byte 0x10;
// 0xffffff80005534c4 <bsd_ast+180>:	mov    0x18(%r14),%rdi
	.byte 0x49; .byte 0x8b; .byte 0x7e; .byte 0x18;
// 0xffffff80005534c8 <bsd_ast+184>:	mov    $0x1,%esi
	.byte 0xbe; .byte 0x01; .byte 0x00; .byte 0x00; .byte 0x00;
// 0xffffff80005534cd <bsd_ast+189>:	callq  0xffffff8000236090 <task_vtimer_set>
	call _task_vtimer_set
// 0xffffff80005534d2 <bsd_ast+194>:	jmp    0xffffff80005534e2 <bsd_ast+210>
	.byte 0xeb; .byte 0x0e;
// 0xffffff80005534d4 <bsd_ast+196>:	mov    0x18(%r14),%rdi
	.byte 0x49; .byte 0x8b; .byte 0x7e; .byte 0x18;
// 0xffffff80005534d8 <bsd_ast+200>:	mov    $0x1,%esi
	.byte 0xbe; .byte 0x01; .byte 0x00; .byte 0x00; .byte 0x00;
// 0xffffff80005534dd <bsd_ast+205>:	callq  0xffffff8000236060 <task_vtimer_clear>
	call _task_vtimer_clear
// 0xffffff80005534e2 <bsd_ast+210>:	xor    %esi,%esi
	.byte 0x31; .byte 0xf6;
// 0xffffff80005534e4 <bsd_ast+212>:	xor    %ecx,%ecx
	.byte 0x31; .byte 0xc9;
// 0xffffff80005534e6 <bsd_ast+214>:	mov    $0x1a,%r8d
	.byte 0x41; .byte 0xb8; .byte 0x1a; .byte 0x00; .byte 0x00; .byte 0x00;
// 0xffffff80005534ec <bsd_ast+220>:	mov    %r14,%rdi
	.byte 0x4c; .byte 0x89; .byte 0xf7;
// 0xffffff80005534ef <bsd_ast+223>:	xor    %edx,%edx
	.byte 0x31; .byte 0xd2;
// 0xffffff80005534f1 <bsd_ast+225>:	callq  0xffffff80005523c0 <setsigvec+560>
	call _psignal_internal
// 0xffffff80005534f6 <bsd_ast+230>:	cmpq   $0x0,0x1e0(%r14)
	.byte 0x49; .byte 0x83; .byte 0xbe; .byte 0xe0; .byte 0x01; .byte 0x00; .byte 0x00; .byte 0x00;
// 0xffffff80005534fe <bsd_ast+238>:	jne    0xffffff800055350a <bsd_ast+250>
	.byte 0x75; .byte 0x0a;
// 0xffffff8000553500 <bsd_ast+240>:	cmpl   $0x0,0x1e8(%r14)
	.byte 0x41; .byte 0x83; .byte 0xbe; .byte 0xe8; .byte 0x01; .byte 0x00; .byte 0x00; .byte 0x00;
// 0xffffff8000553508 <bsd_ast+248>:	je     0xffffff8000553578 <bsd_ast+360>
	.byte 0x74; .byte 0x6e;
// 0xffffff800055350a <bsd_ast+250>:	mov    0x18(%r14),%rdi
	.byte 0x49; .byte 0x8b; .byte 0x7e; .byte 0x18;
// 0xffffff800055350e <bsd_ast+254>:	mov    $0x2,%esi
	.byte 0xbe; .byte 0x02; .byte 0x00; .byte 0x00; .byte 0x00;
// 0xffffff8000553513 <bsd_ast+259>:	lea    -0x20(%rbp),%rdx
	.byte 0x48; .byte 0x8d; .byte 0x55; .byte 0xe0;
// 0xffffff8000553517 <bsd_ast+263>:	callq  0xffffff8000236330 <task_vtimer_update>
	.byte 0xe8; .byte 0x14; .byte 0x2e; .byte 0xce; .byte 0xff;
// 0xffffff800055351c <bsd_ast+268>:	mov    -0x20(%rbp),%edx
	.byte 0x8b; .byte 0x55; .byte 0xe0;
// 0xffffff800055351f <bsd_ast+271>:	lea    0x1d0(%r14),%rsi
	.byte 0x49; .byte 0x8d; .byte 0xb6; .byte 0xd0; .byte 0x01; .byte 0x00; .byte 0x00;
// 0xffffff8000553526 <bsd_ast+278>:	mov    %r14,%rdi
	.byte 0x4c; .byte 0x89; .byte 0xf7;
// 0xffffff8000553529 <bsd_ast+281>:	callq  0xffffff800055c640 <itimerdecr>
	.byte 0xe8; .byte 0x12; .byte 0x91; .byte 0x00; .byte 0x00;
// 0xffffff800055352e <bsd_ast+286>:	test   %eax,%eax
	.byte 0x85; .byte 0xc0;
// 0xffffff8000553530 <bsd_ast+288>:	jne    0xffffff8000553578 <bsd_ast+360>
	.byte 0x75; .byte 0x46;
// 0xffffff8000553532 <bsd_ast+290>:	cmpq   $0x0,0x1e0(%r14)
	.byte 0x49; .byte 0x83; .byte 0xbe; .byte 0xe0; .byte 0x01; .byte 0x00; .byte 0x00; .byte 0x00;
// 0xffffff800055353a <bsd_ast+298>:	jne    0xffffff8000553546 <bsd_ast+310>
	.byte 0x75; .byte 0x0a;
// 0xffffff800055353c <bsd_ast+300>:	cmpl   $0x0,0x1e8(%r14)
	.byte 0x41; .byte 0x83; .byte 0xbe; .byte 0xe8; .byte 0x01; .byte 0x00; .byte 0x00; .byte 0x00;
// 0xffffff8000553544 <bsd_ast+308>:	je     0xffffff8000553556 <bsd_ast+326>
	.byte 0x74; .byte 0x10;
// 0xffffff8000553546 <bsd_ast+310>:	mov    0x18(%r14),%rdi
	.byte 0x49; .byte 0x8b; .byte 0x7e; .byte 0x18;
// 0xffffff800055354a <bsd_ast+314>:	mov    $0x2,%esi
	.byte 0xbe; .byte 0x02; .byte 0x00; .byte 0x00; .byte 0x00;
// 0xffffff800055354f <bsd_ast+319>:	callq  0xffffff8000236090 <task_vtimer_set>
	call _task_vtimer_set
// 0xffffff8000553554 <bsd_ast+324>:	jmp    0xffffff8000553564 <bsd_ast+340>
	.byte 0xeb; .byte 0x0e;
// 0xffffff8000553556 <bsd_ast+326>:	mov    0x18(%r14),%rdi
	.byte 0x49; .byte 0x8b; .byte 0x7e; .byte 0x18;
// 0xffffff800055355a <bsd_ast+330>:	mov    $0x2,%esi
	.byte 0xbe; .byte 0x02; .byte 0x00; .byte 0x00; .byte 0x00;
// 0xffffff800055355f <bsd_ast+335>:	callq  0xffffff8000236060 <task_vtimer_clear>
	call _task_vtimer_clear
// 0xffffff8000553564 <bsd_ast+340>:	xor    %esi,%esi
	.byte 0x31; .byte 0xf6;
// 0xffffff8000553566 <bsd_ast+342>:	xor    %ecx,%ecx
	.byte 0x31; .byte 0xc9;
// 0xffffff8000553568 <bsd_ast+344>:	mov    $0x1b,%r8d
	.byte 0x41; .byte 0xb8; .byte 0x1b; .byte 0x00; .byte 0x00; .byte 0x00;
// 0xffffff800055356e <bsd_ast+350>:	mov    %r14,%rdi
	.byte 0x4c; .byte 0x89; .byte 0xf7;
// 0xffffff8000553571 <bsd_ast+353>:	xor    %edx,%edx
	.byte 0x31; .byte 0xd2;
// 0xffffff8000553573 <bsd_ast+355>:	callq  0xffffff80005523c0 <setsigvec+560>
	call _psignal_internal
// 0xffffff8000553578 <bsd_ast+360>:	cmpq   $0x0,0x1f0(%r14)
	.byte 0x49; .byte 0x83; .byte 0xbe; .byte 0xf0; .byte 0x01; .byte 0x00; .byte 0x00; .byte 0x00;
// 0xffffff8000553580 <bsd_ast+368>:	jne    0xffffff8000553590 <bsd_ast+384>
	.byte 0x75; .byte 0x0e;
// 0xffffff8000553582 <bsd_ast+370>:	cmpl   $0x0,0x1f8(%r14)
	.byte 0x41; .byte 0x83; .byte 0xbe; .byte 0xf8; .byte 0x01; .byte 0x00; .byte 0x00; .byte 0x00;
// 0xffffff800055358a <bsd_ast+378>:	je     0xffffff800055363e <bsd_ast+558>
	.byte 0x0f; .byte 0x84; .byte 0xae; .byte 0x00; .byte 0x00; .byte 0x00;
// 0xffffff8000553590 <bsd_ast+384>:	mov    0x18(%r14),%rdi
	.byte 0x49; .byte 0x8b; .byte 0x7e; .byte 0x18;
// 0xffffff8000553594 <bsd_ast+388>:	mov    $0x4,%esi
	.byte 0xbe; .byte 0x04; .byte 0x00; .byte 0x00; .byte 0x00;
// 0xffffff8000553599 <bsd_ast+393>:	lea    -0x28(%rbp),%rdx
	.byte 0x48; .byte 0x8d; .byte 0x55; .byte 0xd8;
// 0xffffff800055359d <bsd_ast+397>:	callq  0xffffff8000236330 <task_vtimer_update>
	.byte 0xe8; .byte 0x8e; .byte 0x2d; .byte 0xce; .byte 0xff;
// 0xffffff80005535a2 <bsd_ast+402>:	mov    %r14,%rdi
	.byte 0x4c; .byte 0x89; .byte 0xf7;
// 0xffffff80005535a5 <bsd_ast+405>:	callq  0xffffff80005460f0 <proc_spinlock>
	.byte 0xe8; .byte 0x46; .byte 0x2b; .byte 0xff; .byte 0xff;
// 0xffffff80005535aa <bsd_ast+410>:	cmpq   $0x0,0x1f0(%r14)
	.byte 0x49; .byte 0x83; .byte 0xbe; .byte 0xf0; .byte 0x01; .byte 0x00; .byte 0x00; .byte 0x00;
// 0xffffff80005535b2 <bsd_ast+418>:	jle    0xffffff80005535c0 <bsd_ast+432>
	.byte 0x7e; .byte 0x0c;
// 0xffffff80005535b4 <bsd_ast+420>:	lea    0x1f8(%r14),%rax
	.byte 0x49; .byte 0x8d; .byte 0x86; .byte 0xf8; .byte 0x01; .byte 0x00; .byte 0x00;
// 0xffffff80005535bb <bsd_ast+427>:	mov    -0x28(%rbp),%ecx
	.byte 0x8b; .byte 0x4d; .byte 0xd8;
// 0xffffff80005535be <bsd_ast+430>:	jmp    0xffffff80005535d3 <bsd_ast+451>
	.byte 0xeb; .byte 0x13;
// 0xffffff80005535c0 <bsd_ast+432>:	mov    -0x28(%rbp),%ecx
	.byte 0x8b; .byte 0x4d; .byte 0xd8;
// 0xffffff80005535c3 <bsd_ast+435>:	cmp    %ecx,0x1f8(%r14)
	.byte 0x41; .byte 0x39; .byte 0x8e; .byte 0xf8; .byte 0x01; .byte 0x00; .byte 0x00;
// 0xffffff80005535ca <bsd_ast+442>:	jle    0xffffff80005535fe <bsd_ast+494>
	.byte 0x7e; .byte 0x32;
// 0xffffff80005535cc <bsd_ast+444>:	lea    0x1f8(%r14),%rax
	.byte 0x49; .byte 0x8d; .byte 0x86; .byte 0xf8; .byte 0x01; .byte 0x00; .byte 0x00;
// 0xffffff80005535d3 <bsd_ast+451>:	movq   $0x0,-0x30(%rbp)
	.byte 0x48; .byte 0xc7; .byte 0x45; .byte 0xd0; .byte 0x00; .byte 0x00; .byte 0x00; .byte 0x00;
// 0xffffff80005535db <bsd_ast+459>:	mov    (%rax),%edx
	.byte 0x8b; .byte 0x10;
// 0xffffff80005535dd <bsd_ast+461>:	sub    %ecx,%edx
	.byte 0x29; .byte 0xca;
// 0xffffff80005535df <bsd_ast+463>:	mov    %edx,(%rax)
	.byte 0x89; .byte 0x10;
// 0xffffff80005535e1 <bsd_ast+465>:	test   %edx,%edx
	.byte 0x85; .byte 0xd2;
// 0xffffff80005535e3 <bsd_ast+467>:	jns    0xffffff80005535f4 <bsd_ast+484>
	.byte 0x79; .byte 0x0f;
// 0xffffff80005535e5 <bsd_ast+469>:	decq   0x1f0(%r14)
	.byte 0x49; .byte 0xff; .byte 0x8e; .byte 0xf0; .byte 0x01; .byte 0x00; .byte 0x00;
// 0xffffff80005535ec <bsd_ast+476>:	add    $0xf4240,%edx
	.byte 0x81; .byte 0xc2; .byte 0x40; .byte 0x42; .byte 0x0f; .byte 0x00;
// 0xffffff80005535f2 <bsd_ast+482>:	mov    %edx,(%rax)
	.byte 0x89; .byte 0x10;
// 0xffffff80005535f4 <bsd_ast+484>:	mov    %r14,%rdi
	.byte 0x4c; .byte 0x89; .byte 0xf7;
// 0xffffff80005535f7 <bsd_ast+487>:	callq  0xffffff80005460d0 <proc_spinunlock>
	.byte 0xe8; .byte 0xd4; .byte 0x2a; .byte 0xff; .byte 0xff;
// 0xffffff80005535fc <bsd_ast+492>:	jmp    0xffffff800055363e <bsd_ast+558>
	.byte 0xeb; .byte 0x40;
// 0xffffff80005535fe <bsd_ast+494>:	movl   $0x0,0x1f8(%r14)
	.byte 0x41; .byte 0xc7; .byte 0x86; .byte 0xf8; .byte 0x01; .byte 0x00; .byte 0x00; .byte 0x00; .byte 0x00; .byte 0x00; .byte 0x00;
// 0xffffff8000553609 <bsd_ast+505>:	movq   $0x0,0x1f0(%r14)
	.byte 0x49; .byte 0xc7; .byte 0x86; .byte 0xf0; .byte 0x01; .byte 0x00; .byte 0x00; .byte 0x00; .byte 0x00; .byte 0x00; .byte 0x00;
// 0xffffff8000553614 <bsd_ast+516>:	mov    %r14,%rdi
	.byte 0x4c; .byte 0x89; .byte 0xf7;
// 0xffffff8000553617 <bsd_ast+519>:	callq  0xffffff80005460d0 <proc_spinunlock>
	.byte 0xe8; .byte 0xb4; .byte 0x2a; .byte 0xff; .byte 0xff;
// 0xffffff800055361c <bsd_ast+524>:	mov    0x18(%r14),%rdi
	.byte 0x49; .byte 0x8b; .byte 0x7e; .byte 0x18;
// 0xffffff8000553620 <bsd_ast+528>:	mov    $0x4,%esi
	.byte 0xbe; .byte 0x04; .byte 0x00; .byte 0x00; .byte 0x00;
// 0xffffff8000553625 <bsd_ast+533>:	callq  0xffffff8000236060 <task_vtimer_clear>
	call _task_vtimer_clear
// 0xffffff800055362a <bsd_ast+538>:	xor    %esi,%esi
	.byte 0x31; .byte 0xf6;
// 0xffffff800055362c <bsd_ast+540>:	xor    %ecx,%ecx
	.byte 0x31; .byte 0xc9;
// 0xffffff800055362e <bsd_ast+542>:	mov    $0x18,%r8d
	.byte 0x41; .byte 0xb8; .byte 0x18; .byte 0x00; .byte 0x00; .byte 0x00;
// 0xffffff8000553634 <bsd_ast+548>:	mov    %r14,%rdi
	.byte 0x4c; .byte 0x89; .byte 0xf7;
// 0xffffff8000553637 <bsd_ast+551>:	xor    %edx,%edx
	.byte 0x31; .byte 0xd2;
// 0xffffff8000553639 <bsd_ast+553>:	callq  0xffffff80005523c0 <setsigvec+560>
	call _psignal_internal
// 0xffffff800055363e <bsd_ast+558>:	movzbl 0x25d(%rbx),%r8d
	.byte 0x44; .byte 0x0f; .byte 0xb6; .byte 0x83; .byte 0x5d; .byte 0x02; .byte 0x00; .byte 0x00;
// 0xffffff8000553646 <bsd_ast+566>:	test   %r8d,%r8d
	.byte 0x45; .byte 0x85; .byte 0xc0;
// 0xffffff8000553649 <bsd_ast+569>:	je     0xffffff8000553660 <bsd_ast+592>
	.byte 0x74; .byte 0x15;
// 0xffffff800055364b <bsd_ast+571>:	movb   $0x0,0x25d(%rbx)
	.byte 0xc6; .byte 0x83; .byte 0x5d; .byte 0x02; .byte 0x00; .byte 0x00; .byte 0x00;
// 0xffffff8000553652 <bsd_ast+578>:	xor    %esi,%esi
	.byte 0x31; .byte 0xf6;
// 0xffffff8000553654 <bsd_ast+580>:	xor    %ecx,%ecx
	.byte 0x31; .byte 0xc9;
// 0xffffff8000553656 <bsd_ast+582>:	mov    %r14,%rdi
	.byte 0x4c; .byte 0x89; .byte 0xf7;
// 0xffffff8000553659 <bsd_ast+585>:	xor    %edx,%edx
	.byte 0x31; .byte 0xd2;
// 0xffffff800055365b <bsd_ast+587>:	callq  0xffffff80005523c0 <setsigvec+560>
	.byte 0xe8; .byte 0x60; .byte 0xed; .byte 0xff; .byte 0xff;
// 0xffffff8000553660 <bsd_ast+592>:	cmpb   $0x0,0x25c(%rbx)
	.byte 0x80; .byte 0xbb; .byte 0x5c; .byte 0x02; .byte 0x00; .byte 0x00; .byte 0x00;
// 0xffffff8000553667 <bsd_ast+599>:	je     0xffffff8000553691 <bsd_ast+641>
	.byte 0x74; .byte 0x28;
// 0xffffff8000553669 <bsd_ast+601>:	movb   $0x0,0x25c(%rbx)
	.byte 0xc6; .byte 0x83; .byte 0x5c; .byte 0x02; .byte 0x00; .byte 0x00; .byte 0x00;
// 0xffffff8000553670 <bsd_ast+608>:	mov    %r14,%rdi
	.byte 0x4c; .byte 0x89; .byte 0xf7;
// 0xffffff8000553673 <bsd_ast+611>:	callq  0xffffff8000545d10 <proc_lock>
	.byte 0xe8; .byte 0x98; .byte 0x26; .byte 0xff; .byte 0xff;
// 0xffffff8000553678 <bsd_ast+616>:	movb   $0x1,0x270(%r14)
	.byte 0x41; .byte 0xc6; .byte 0x86; .byte 0x70; .byte 0x02; .byte 0x00; .byte 0x00; .byte 0x01;
// 0xffffff8000553680 <bsd_ast+624>:	mov    %r14,%rdi
	.byte 0x4c; .byte 0x89; .byte 0xf7;
// 0xffffff8000553683 <bsd_ast+627>:	callq  0xffffff8000545ce0 <proc_unlock>
	.byte 0xe8; .byte 0x58; .byte 0x26; .byte 0xff; .byte 0xff;
// 0xffffff8000553688 <bsd_ast+632>:	mov    0x18(%r14),%rdi
	.byte 0x49; .byte 0x8b; .byte 0x7e; .byte 0x18;
// 0xffffff800055368c <bsd_ast+636>:	callq  0xffffff8000237270 <task_suspend>
	.byte 0xe8; .byte 0xdf; .byte 0x3b; .byte 0xce; .byte 0xff;
// 0xffffff8000553691 <bsd_ast+641>:	mov    0x260(%rbx),%rdi
	.byte 0x48; .byte 0x8b; .byte 0xbb; .byte 0x60; .byte 0x02; .byte 0x00; .byte 0x00;
// 0xffffff8000553698 <bsd_ast+648>:	test   %rdi,%rdi
	.byte 0x48; .byte 0x85; .byte 0xff;
// 0xffffff800055369b <bsd_ast+651>:	je     0xffffff80005536f2 <bsd_ast+738>
	.byte 0x74; .byte 0x55;
// 0xffffff800055369d <bsd_ast+653>:	callq  0xffffff800054a610 <proc_find>
	.byte 0xe8; .byte 0x6e; .byte 0x6f; .byte 0xff; .byte 0xff;
// 0xffffff80005536a2 <bsd_ast+658>:	movq   $0x0,0x260(%rbx)
	.byte 0x48; .byte 0xc7; .byte 0x83; .byte 0x60; .byte 0x02; .byte 0x00; .byte 0x00; .byte 0x00; .byte 0x00; .byte 0x00; .byte 0x00;
// 0xffffff80005536ad <bsd_ast+669>:	test   %rax,%rax
	.byte 0x48; .byte 0x85; .byte 0xc0;
// 0xffffff80005536b0 <bsd_ast+672>:	je     0xffffff80005536f2 <bsd_ast+738>
	.byte 0x74; .byte 0x40;
// 0xffffff80005536b2 <bsd_ast+674>:	mov    %rax,%r15
	.byte 0x49; .byte 0x89; .byte 0xc7;
// 0xffffff80005536b5 <bsd_ast+677>:	mov    %r15,%rdi
	.byte 0x4c; .byte 0x89; .byte 0xff;
// 0xffffff80005536b8 <bsd_ast+680>:	callq  0xffffff8000545d10 <proc_lock>
	.byte 0xe8; .byte 0x53; .byte 0x26; .byte 0xff; .byte 0xff;
// 0xffffff80005536bd <bsd_ast+685>:	cmpb   $0x0,0x270(%r15)
	.byte 0x41; .byte 0x80; .byte 0xbf; .byte 0x70; .byte 0x02; .byte 0x00; .byte 0x00; .byte 0x00;
// 0xffffff80005536c5 <bsd_ast+693>:	je     0xffffff80005536e2 <bsd_ast+722>
	.byte 0x74; .byte 0x1b;
// 0xffffff80005536c7 <bsd_ast+695>:	movb   $0x0,0x270(%r15)
	.byte 0x41; .byte 0xc6; .byte 0x87; .byte 0x70; .byte 0x02; .byte 0x00; .byte 0x00; .byte 0x00;
// 0xffffff80005536cf <bsd_ast+703>:	mov    %r15,%rdi
	.byte 0x4c; .byte 0x89; .byte 0xff;
// 0xffffff80005536d2 <bsd_ast+706>:	callq  0xffffff8000545ce0 <proc_unlock>
	.byte 0xe8; .byte 0x09; .byte 0x26; .byte 0xff; .byte 0xff;
// 0xffffff80005536d7 <bsd_ast+711>:	mov    0x18(%r15),%rdi
	.byte 0x49; .byte 0x8b; .byte 0x7f; .byte 0x18;
// 0xffffff80005536db <bsd_ast+715>:	callq  0xffffff8000237040 <task_resume>
	.byte 0xe8; .byte 0x60; .byte 0x39; .byte 0xce; .byte 0xff;
// 0xffffff80005536e0 <bsd_ast+720>:	jmp    0xffffff80005536ea <bsd_ast+730>
	.byte 0xeb; .byte 0x08;
// 0xffffff80005536e2 <bsd_ast+722>:	mov    %r15,%rdi
	.byte 0x4c; .byte 0x89; .byte 0xff;
// 0xffffff80005536e5 <bsd_ast+725>:	callq  0xffffff8000545ce0 <proc_unlock>
	.byte 0xe8; .byte 0xf6; .byte 0x25; .byte 0xff; .byte 0xff;
// 0xffffff80005536ea <bsd_ast+730>:	mov    %r15,%rdi
	.byte 0x4c; .byte 0x89; .byte 0xff;
// 0xffffff80005536ed <bsd_ast+733>:	callq  0xffffff800054a3e0 <proc_rele>
	.byte 0xe8; .byte 0xee; .byte 0x6c; .byte 0xff; .byte 0xff;
// 0xffffff80005536f2 <bsd_ast+738>:	callq  0xffffff80002be1c0 <current_thread>
	.byte 0xe8; .byte 0xc9; .byte 0xaa; .byte 0xd6; .byte 0xff;
// 0xffffff80005536f7 <bsd_ast+743>:	mov    %rax,%rdi
	.byte 0x48; .byte 0x89; .byte 0xc7;
// 0xffffff80005536fa <bsd_ast+746>:	callq  0xffffff800023a080 <thread_should_halt>
	.byte 0xe8; .byte 0x81; .byte 0x69; .byte 0xce; .byte 0xff;
// 0xffffff80005536ff <bsd_ast+751>:	test   %eax,%eax
	.byte 0x85; .byte 0xc0;
// 0xffffff8000553701 <bsd_ast+753>:	jne    0xffffff8000553763 <bsd_ast+851>
	.byte 0x75; .byte 0x60;
// 0xffffff8000553703 <bsd_ast+755>:	testb  $0x4,0x15d(%r14)
	.byte 0x41; .byte 0xf6; .byte 0x86; .byte 0x5d; .byte 0x01; .byte 0x00; .byte 0x00; .byte 0x04;
// 0xffffff800055370b <bsd_ast+763>:	mov    0x148(%rbx),%eax
	.byte 0x8b; .byte 0x83; .byte 0x48; .byte 0x01; .byte 0x00; .byte 0x00;
// 0xffffff8000553711 <bsd_ast+769>:	mov    0x150(%rbx),%ecx
	.byte 0x8b; .byte 0x8b; .byte 0x50; .byte 0x01; .byte 0x00; .byte 0x00;
// 0xffffff8000553717 <bsd_ast+775>:	je     0xffffff800055371d <bsd_ast+781>
	.byte 0x74; .byte 0x04;
// 0xffffff8000553719 <bsd_ast+777>:	xor    %edx,%edx
	.byte 0x31; .byte 0xd2;
// 0xffffff800055371b <bsd_ast+779>:	jmp    0xffffff8000553724 <bsd_ast+788>
	.byte 0xeb; .byte 0x07;
// 0xffffff800055371d <bsd_ast+781>:	mov    0x2bc(%r14),%edx
	.byte 0x41; .byte 0x8b; .byte 0x96; .byte 0xbc; .byte 0x02; .byte 0x00; .byte 0x00;
// 0xffffff8000553724 <bsd_ast+788>:	or     %edx,%ecx
	.byte 0x09; .byte 0xd1;
// 0xffffff8000553726 <bsd_ast+790>:	or     $0x10100,%ecx
	.byte 0x81; .byte 0xc9; .byte 0x00; .byte 0x01; .byte 0x01; .byte 0x00;
// 0xffffff800055372c <bsd_ast+796>:	xor    $0xfffefeff,%ecx
	.byte 0x81; .byte 0xf1; .byte 0xff; .byte 0xfe; .byte 0xfe; .byte 0xff;
// 0xffffff8000553732 <bsd_ast+802>:	test   %ecx,%eax
	.byte 0x85; .byte 0xc8;
// 0xffffff8000553734 <bsd_ast+804>:	je     0xffffff8000553763 <bsd_ast+851>
	.byte 0x74; .byte 0x2d;
// 0xffffff8000553736 <bsd_ast+806>:	mov    %r14,%rdi
	.byte 0x4c; .byte 0x89; .byte 0xf7;
// 0xffffff8000553739 <bsd_ast+809>:	callq  0xffffff8000552d60 <issignal>
	.byte 0xe8; .byte 0x22; .byte 0xf6; .byte 0xff; .byte 0xff;
// 0xffffff800055373e <bsd_ast+814>:	test   %eax,%eax
	.byte 0x85; .byte 0xc0;
// 0xffffff8000553740 <bsd_ast+816>:	je     0xffffff8000553763 <bsd_ast+851>
	.byte 0x74; .byte 0x21;
// 0xffffff8000553742 <bsd_ast+818>:	mov    %eax,%edi
	.byte 0x89; .byte 0xc7;
// 0xffffff8000553744 <bsd_ast+820>:	nopw   0x0(%rax,%rax,1)
	.byte 0x66; .byte 0x0f; .byte 0x1f; .byte 0x44; .byte 0x00; .byte 0x00;
// 0xffffff800055374a <bsd_ast+826>:	nopw   0x0(%rax,%rax,1)
	.byte 0x66; .byte 0x0f; .byte 0x1f; .byte 0x44; .byte 0x00; .byte 0x00;
// 0xffffff8000553750 <bsd_ast+832>:	callq  0xffffff8000551aa0 <postsig>
	.byte 0xe8; .byte 0x4b; .byte 0xe3; .byte 0xff; .byte 0xff;
// 0xffffff8000553755 <bsd_ast+837>:	mov    %r14,%rdi
	.byte 0x4c; .byte 0x89; .byte 0xf7;
// 0xffffff8000553758 <bsd_ast+840>:	callq  0xffffff8000552d60 <issignal>
	.byte 0xe8; .byte 0x03; .byte 0xf6; .byte 0xff; .byte 0xff;
// 0xffffff800055375d <bsd_ast+845>:	test   %eax,%eax
	.byte 0x85; .byte 0xc0;
// 0xffffff800055375f <bsd_ast+847>:	mov    %eax,%edi
	.byte 0x89; .byte 0xc7;
// 0xffffff8000553761 <bsd_ast+849>:	jne    0xffffff8000553750 <bsd_ast+832>
	.byte 0x75; .byte 0xed;
// 0xffffff8000553763 <bsd_ast+851>:	mov    0x356167(%rip),%al        # 0xffffff80008a98d0
	.byte 0x8a; .byte 0x05; .byte 0x67; .byte 0x61; .byte 0x35; .byte 0x00;
// 0xffffff8000553769 <bsd_ast+857>:	test   %al,%al
	.byte 0x84; .byte 0xc0;
// 0xffffff800055376b <bsd_ast+859>:	jne    0xffffff8000553779 <bsd_ast+873>
	.byte 0x75; .byte 0x0c;
// 0xffffff800055376d <bsd_ast+861>:	movb   $0x1,0x35615c(%rip)        # 0xffffff80008a98d0
	.byte 0xc6; .byte 0x05; .byte 0x5c; .byte 0x61; .byte 0x35; .byte 0x00; .byte 0x01;
// 0xffffff8000553774 <bsd_ast+868>:	callq  0xffffff8000527680 <bsdinit_task>
	.byte 0xe8; .byte 0x07; .byte 0x3f; .byte 0xfd; .byte 0xff;
// 0xffffff8000553779 <bsd_ast+873>:	add    $0x18,%rsp
	.byte 0x48; .byte 0x83; .byte 0xc4; .byte 0x18;
// 0xffffff800055377d <bsd_ast+877>:	pop    %rbx
	.byte 0x5b;
// 0xffffff800055377e <bsd_ast+878>:	pop    %r14
	.byte 0x41; .byte 0x5e;
// 0xffffff8000553780 <bsd_ast+880>:	pop    %r15
	.byte 0x41; .byte 0x5f;
// 0xffffff8000553782 <bsd_ast+882>:	pop    %rbp
	.byte 0x5d;
// 0xffffff8000553783 <bsd_ast+883>:	retq   
	.byte 0xc3;
// 0xffffff8000553784 <bsd_ast+884>:	nopw   0x0(%rax,%rax,1)
	.byte 0x66; .byte 0x0f; .byte 0x1f; .byte 0x44; .byte 0x00; .byte 0x00;
// 0xffffff800055378a <bsd_ast+890>:	nopw   0x0(%rax,%rax,1)
	.byte 0x66; .byte 0x0f; .byte 0x1f; .byte 0x44; .byte 0x00; .byte 0x00;
// 0xffffff8000553790 <bsd_ast+896>:	push   %rbp
	.byte 0x55;
// 0xffffff8000553791 <bsd_ast+897>:	mov    %rsp,%rbp
	.byte 0x48; .byte 0x89; .byte 0xe5;
// 0xffffff8000553794 <bsd_ast+900>:	mov    (%rsi),%r8d
	.byte 0x44; .byte 0x8b; .byte 0x06;
// 0xffffff8000553797 <bsd_ast+903>:	xor    %esi,%esi
	.byte 0x31; .byte 0xf6;
// 0xffffff8000553799 <bsd_ast+905>:	xor    %edx,%edx
	.byte 0x31; .byte 0xd2;
// 0xffffff800055379b <bsd_ast+907>:	xor    %ecx,%ecx
	.byte 0x31; .byte 0xc9;
// 0xffffff800055379d <bsd_ast+909>:	callq  0xffffff80005523c0 <setsigvec+560>
	.byte 0xe8; .byte 0x1e; .byte 0xec; .byte 0xff; .byte 0xff;
// 0xffffff80005537a2 <bsd_ast+914>:	xor    %eax,%eax
	.byte 0x31; .byte 0xc0;
// 0xffffff80005537a4 <bsd_ast+916>:	pop    %rbp
	.byte 0x5d;
// 0xffffff80005537a5 <bsd_ast+917>:	retq   
	.byte 0xc3;
// 0xffffff80005537a6 <bsd_ast+918>:	nopw   %cs:0x0(%rax,%rax,1)
	.byte 0x66; .byte 0x2e; .byte 0x0f; .byte 0x1f; .byte 0x84; .byte 0x00; .byte 0x00; .byte 0x00; .byte 0x00; .byte 0x00;
