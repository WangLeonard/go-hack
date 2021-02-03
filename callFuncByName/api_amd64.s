#include "go_asm.h"
#include "textflag.h"
#include "funcdata.h"

TEXT ·hackCallWithStructArg(SB),$16-0
    NO_LOCAL_POINTERS
	MOVQ fn+8(FP), AX
	MOVQ fn+16(FP), BX
	MOVQ BX,0(SP)
	MOVQ 0(AX),AX
	CALL AX
	RET


TEXT ·hackCallWithStructInterface(SB),$16-0
    NO_LOCAL_POINTERS
	MOVQ fn+8(FP), AX
	MOVQ fn+16(FP), BX
	MOVQ BX,0(SP)
	MOVQ 0(AX),AX
	CALL AX
	RET
