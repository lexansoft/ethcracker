# Ethereum password cracker

If you forgot your password for Ethereum private key or presale file, but think you still remember 
the list of possible substrings from wich you constracted the password, then try this progam to 
routinly go through all the possible combinations and find the working password.

PLEASE DO NOT TRUST ANYONE TO COMPILE THE PROGRAM FOR YOU. ALWAYS USE THE SOURCE CODE DOWNLOADED FROM 
THE GITHUB. THIS WAY YOU CAN BE SURE, THE PROGRAM DOES NOT HAVE ANY MALICIOUS CODE !!!

# Gitter Channel

https://gitter.im/lexansoft/ethcracker

# Usage 

    ethcracker -pk ~/test/pk.txt -t ~/test/templates.txt

    -pk path to the private key file
    -t  path to the template file
    -l  path to the file with all the possible variants (every line has one variant) If -l is specified, -t is ignored
    -presale  for cracking presale JSON file
    -threads Number of threads
    -v Verbosity ( 0, 1, 2 )
    -start_from Skip first N combinations ( you can specify N as percentage. F.e. : 30% )
    -keep_order Keep the order of the lines ( no permutations )
    -re Report every N-th combination
    -dump path Just dump all the variants into text file
    -min_len Minimum password length
    -max_len Maximum password length
    

# Template file format

Every line contains the possible variants of the substring. For example file:

    a1 a2 a3
    b

will generate all those combinations

    a1
    a2
    a3
    b
    a1b
    ba1
    a2b
    ba2
    a3b
    ba3

Note: you can use \s to specify white space. ( "a\sb" means "a b" )


# Template line flags 

You can also specify some keys for every line. 

    ~a always use some value from this string
    ~c Try both: capitalized and not-capitalized versions of all words. 
    
For example the template file

    a 
    ~ac test

will generate all those combinations

    test
    Test
    atest
    aTest
    testa
    Testa


# Installing

Install Go Language


    git clone https://github.com/lexansoft/ethcracker
    cd ethcracker
    
    go get github.com/ethereum/go-ethereum
    go get github.com/pborman/uuid
    go get golang.org/x/crypto/pbkdf2
    go get golang.org/x/crypto/ripemd160
    go get golang.org/x/crypto/scrypt
    go get github.com/rjeczalik/notify
    
    
# To update all previously installed packages use this:    
    
    go get -u all
    
# Run the cracker    
    
    go run src/ethcracker.go -pk PATH_TO_FILE -t PATH_TO_TEMPLATE_FILE -threads 4 
    
# Installing on Windows 

On windwos you need to install the Chocolatey:  https://chocolatey.org 

Then install git, golang and mingw 

    choco install git
    choco install golang
    choco install mingw
    
After that make all the steps from the Installing section.    
    
# Cracking your Mac DMG file password
If you stored your keys in the encrypted mac DMG image and forgot the password, do this:

1. dump all the possible variants of your password into a file

        go run src/ethcracker.go ... -dump ~/v.txt 
        
2. Use dmg_pass.bash script to try all the variants form v.txt

        ./dmg_pass.bash v.txt your.dmg
        
        or
        
        ./dmg_pass.bash v.txt your.dmg N -- to skipp first N lines

# Donation

If this program helped you to restore the password, please donate some ETH to the address:

 0x281694Fabfdd9735e01bB59942B18c469b6e3df6
 
 Thank you
 

 
 
