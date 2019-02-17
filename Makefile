DISPLAY="DISPLAY=:1"

THREADS="--threads:on"

CC=nim

mwm: $CC "c" $THREADS main.nim

test: $DISPLAY $CC "c" $THREADS -r main.nim