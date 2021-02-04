#include "go_asm.h"
#include "textflag.h"
#include "funcdata.h"

TEXT ·hackCallWithStructArg(SB),$24-8
    NO_LOCAL_POINTERS
	MOVQ fn+8(FP), AX
	MOVQ args+16(FP), BX
	MOVQ BX,0(SP)
	MOVQ 0(AX),AX
	CALL AX
	RET


TEXT ·hackCallWithStructInterface(SB),$32-16
    NO_LOCAL_POINTERS
	MOVQ fn+8(FP), AX
	MOVQ itype+16(FP), BX
	MOVQ idata+24(FP), CX
	MOVQ CX,8(SP)
	MOVQ BX,0(SP)
	MOVQ 0(AX),AX
	CALL AX
	RET
