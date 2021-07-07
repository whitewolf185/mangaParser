const util = require('util');
const https = require('follow-redirects').https;
const fs = require('fs');
const writeFile = util.promisify(fs.writeFile)
const URL = '/series/abyss/chapters/1/';
const parallel = require('async/parallel');
const exec = util.promisify(require('child_process').exec);
const nodemailer = require('nodemailer');
//whitewolf_185@pbsync.com

var options = {
    'method': 'GET',
    'hostname': 'images.catmanga.org',
    'path': '/series/abyss/chapters/1/',
    'headers': {
    },
    'maxRedirects': 20
};

async function send_to_Email (filePath, fileName){
    let tranporter = nodemailer.createTransport({
        host: 'smtp.mail.ru',
        port: 465,
        secure: true,
        auth:{
            user: '',
            pass: ''
        }
    })

    let result = tranporter.sendMail({
        from: 'matvey2001xxl1976@mail.ru matvey2001xxl1976@mail.ru',
        to: 'whitewolf_185@pbsync.com whitewolf_185@pbsync.com',
        attachments: {filename: fileName, path: filePath}
    })

    console.log(result);
    return result;
}



function cleanOut () {
    setTimeout(()=>{fs.readdir(__dirname + "/chapter1", (error,files) => {
        if(error) throw error;
        console.info("deleting....\n");

        for (const file of files) {
            let filePath = __dirname + "/chapter1/" + file
            console.log(filePath);
            fs.unlink(filePath, err => {
                if (err) throw err;
            });
        }
    })},1000)
}

let check_dir = async function (path) {
    let result;
    await fs.stat(path, (err) => {
       if(!err){
           result = 'Directory is existing';
       }

       else if (err.code === 'ENOENT'){
           fs.mkdir(path, (e) => {
               if(e){
                   result = e;
               }
           });
       }
    });
    return result;
}

let dwld_imgs = new Promise(((resolve, reject) => {
    check_dir(__dirname + '/chapter1/').then(a => {
        let filesPath = [];
        let b = [];
        let path = "";
        for (let i = 1; i < 52; i++) {
            if (i < 10){
                path = URL+"00"+i+".png";
            }
            else{
                path = URL+"0"+i+".png"
            }

            let a = (optionPath) => {
                return function (callback){
                    options.path = optionPath;
                    let req = https.request(options, function (res) {
                        let chunks = [];

                        if(res.statusCode > 300){
                            return;
                        }

                        res.on("data", function (chunk) {
                            chunks.push(chunk);
                        });

                        res.on("end", function (chunk) {
                            let body = Buffer.concat(chunks);
                            writeFile(__dirname + "/chapter1/img_" + i + ".png", body)
                                .then(result => {
                                    console.log(optionPath + " ready")
                                    callback(null,i);
                                })
                        });

                        res.on("error", function (error) {
                            reject(error);
                        });
                    });
                    req.end();
                }
            }
            filesPath.push({
                id: i,
                path: __dirname + "/chapter1/img" + i
            });


            b.push(a(path));

        }
        //console.log(filesPath);
        parallel(b,(err,result) => {
            console.log(result);
            if (err) throw err;
            setTimeout(() =>{
                resolve(filesPath)
            }, 1000);
        })
    })

}));

dwld_imgs
    .then(
        filesPath => {
            filesPath.sort((lhs, rhs) => {
                return lhs.id - rhs.id;
            });
            console.log("im generating pdf file...");

            exec('./img2pdf.sh chapter1 chapter_1.pdf')
                .then(
                    output => {
                        console.log(output.stdout);
                        /*send_to_Email(__dirname + "/chapter_1.pdf", "Abyss_ch1.pdf")
                            .then(result =>{
                                // cleanOut();
                            })*/
                    }
                )
                .catch(
                    stderr => {
                        console.log(stderr);
                    }
                )

        }
    )

    .catch(
        err => {
            console.log(err);

            //clean out all files, when error
            cleanOut();
        }
    )


