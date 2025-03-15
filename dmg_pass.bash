#!/bin/sh

i=0

while IFS='' read -r line || [[ -n "$line" ]]; do
    
    ((i++))

    if [ $i -lt $3 ]; then
        echo "$i: $line skipped"
        continue
    fi

    echo "$i: $line"
        
    res="$(printf $line | hdiutil attach -quiet -stdinpass $2 )" 

#    echo "RESULT: $?"
    if [ $? -eq 0 ]; then
    
        echo "================================================================"
        echo ""
        echo "            The disk is successfully mounted!!! "
        echo "            Your password is: $line"
        echo ""
        echo "      Do not forget to donate some ETH to the developer:"
        echo "  Ethereum Address: 0x281694Fabfdd9735e01bB59942B18c469b6e3df6"             
        echo "================================================================"
        exit    
    fi
    
done < "$1"