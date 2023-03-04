cd $1
img2pdf `ls *.png | sort -t _ -k 2 -g` -o ../"$2" && echo "chapter 1 complete"

#pdfunite `ls | sort -t _ -k 2 -g` Juilet.pdf && rm chapter*