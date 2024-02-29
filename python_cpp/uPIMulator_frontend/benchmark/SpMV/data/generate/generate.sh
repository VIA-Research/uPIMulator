
R="1 4 16 64 256"

make

for r in $R; do
    ./replicate ../bcsstk30.mtx $r ../bcsstk30.mtx.$r.mtx
done

