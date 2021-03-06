;
; Example of filling two 10x1 arrays (A & B)
; Filling C 40x1 array with the output from A + B
; Sum all C elements of the array in to R20
;    (backward jumps flavor)
;

; Number of words to process
LLI     R20, 40                    ; R20 = 40 (iterations)

; Initializacion memory
LLI     R1, 0                      ; R1 = 0 (index)
LLI     R10, 64                    ; A's memory index
LLI     R11, 320                   ; B's memory index
LLI     R12, 576                   ; C's memory index

MEMORY_LOOP:					   ; for (RI=0; R1 < R20; R1++)
ADDI    R1, R1, 1                  ; R1 += 1

SW      R10, R1, 0                 ; A[R10] = R1
SLI     R11, 10                    ; B[R11] = 10
LW      R13, R11, 0                ; R13 = MEM[B[I]]
ADD     R14, R1, R13               ; R14 = A[I] + B[I]
SW      R12, R14, 0                ; C[R12] = R14 = R1 + B[R11]

; Increment array indexes
ADDI    R10, R10, 4                ; R10 += 1 (A's index)
ADDI    R11, R11, 4                ; R11 += 1 (B's index)
ADDI    R12, R12, 4                ; R11 += 1 (C's index)

BEQ     R1, R20, END_MEMORY_LOOP   ; breaks when R1 = R20
J 		MEMORY_LOOP
	
END_MEMORY_LOOP:

; ---------- Expected Values -----------
;          A      B       C
;         0x40   0x140   0x240
;  0x00    1      10      11
;  0x04    2      10      12
;  0x08    3      10      13
;  0x0C    4      10      14
;  ....    .      .       ..
;  0xC4    50     10      60
; 
; --------------------------------------

; Sum all C's elements
LLI     R1, 0                       ; R1 = 0 (index)
LLI     R12, 576                    ; C's memory index
LLI     R15, 0                      ; Total of C's elements

PROCESS_LOOP:					    ; for (RI=0; R1 < R20; R1++)
ADDI    R1, R1, 1                   ; R1 += 1

LW      R16, R12, 0                 ; R16 = C[I]
ADD     R15, R15, R16               ; R15 += C[I]

; Increment C's index
ADDI    R12, R12, 4                 ; R12 += 1 (C's index)

BEQ     R1, R20, END_PROCESS_LOOP  ; breaks when R1 = R20
J 		PROCESS_LOOP

END_PROCESS_LOOP:

; ---------- Expected Values -----------
;
;   R15 = 11 + 12 + ... + 59 + 60
;   R15 = 0x6EF = 1775
; 
; --------------------------------------