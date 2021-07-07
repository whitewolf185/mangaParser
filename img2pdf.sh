cd $1
img2pdf `ls *.png | sort -t _ -k 2 -g` -o ../"$2" && echo "chapter 1 complete"