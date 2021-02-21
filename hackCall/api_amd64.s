#include "go_asm.h"
#include "textflag.h"
#include "funcdata.h"

// func MethodCall(fn interface{}, obj interface{}, args interface{})
TEXT ·MethodCall(SB),$24-48
    NO_LOCAL_POINTERS
    MOVQ fn_data+8(FP), AX
    MOVQ 0(AX),AX
    MOVQ obj_data+24(FP),BX
    MOVQ arg_itab+32(FP),CX
    MOVQ arg_data+40(FP),DX
    MOVQ BX,obj_data-24(SP)
    MOVQ CX,arg_itab-16(SP)
    MOVQ DX,arg_data-8(SP)
    CALL AX
    RET


// func MethodCallByPtr(fn uintptr, obj interface{}, args interface{})
TEXT ·MethodCallByPtr(SB),$24-40
    NO_LOCAL_POINTERS
    MOVQ fn_data(FP), AX
    MOVQ obj_data+16(FP),BX
    MOVQ arg_itab+24(FP),CX
    MOVQ arg_data+32(FP),DX
    MOVQ BX,obj_data-24(SP)
    MOVQ CX,arg_itab-16(SP)
    MOVQ DX,arg_data-8(SP)
    CALL AX
    RET
