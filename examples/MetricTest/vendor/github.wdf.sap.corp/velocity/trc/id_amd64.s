#define get_tls(r)  MOVQ TLS, r
#define g(r)        0(r)(TLS*1)
#define NOSPLIT     4

TEXT ·id(SB), NOSPLIT, $0
	MOVQ offset+0(FP), BX
	get_tls(CX)
	MOVQ g(CX), CX
	MOVQ 0(BX)(CX*1), AX
	MOVQ AX, ret+8(FP)
	RET

