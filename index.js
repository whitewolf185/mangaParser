const util = require('util');
const https = require('follow-redirects').https;
const fs = require('fs');
const writeFile = util.promisify(fs.writeFile)
const parallel = require('async/parallel');
const exec = util.promisify(require('child_process').exec);
const nodemailer = require('nodemailer');
const fetch = require('node-fetch');

//whitewolf_185@pbsync.com
// ----временная информация----
const manga_name = 'abyss';
const ch_list = [51,33,30,19,19,19,23,18,19,19,19,18,19,18,21,19]

//-----------------------------
let sources = [
    {
        name: "catmanga",
        path: "https://images.catmanga.org/series/abyss/chapters/",
        options: {
            method: 'GET',
            redirect: 'follow'
        }
    },
    {
        name: "mangakakalot",
        options: {
            'method': 'GET',
            'hostname': 's6.mkklcdnv6tempv3.com',
            'path': '/mangakakalot/k1/kaguyasama_wa_kokurasetai_tensaitachi_no_renai_zunousen/chapter_197_the_shirogane_family_wants_to_move/1.jpg',
            'headers': {
                'Accept': 'image/webp,*!/!*',
                'Accept-Encoding': 'gzip, deflate, br',
                'Accept-Language': 'ru-RU,ru;q=0.8,en-US;q=0.5,en;q=0.3',
                'Cache-Control': 'max-age=0',
                'Connection': 'keep-alive',
                'Host': 's6.mkklcdnv6tempv3.com',
                'Referer': 'https://manganelo.com/',
                'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:87.0) Gecko/20100101 Firefox/87.0'
            },
            'maxRedirects': 20
        }
    }
]


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

let dwld_ch = async function(source, ch, img_count){
    return async function (callback){
        await fs.mkdir(__dirname + '/chapter'+ch+'/', (e) => {
            if(e.code === 'EEXIST'){
                console.log('Directory has already exist');
            }
            else if (e) {
                throw e;
            }
        })


        let filesPath = [];
        let b = [];
        let path = "";


        for (let i = 1; i < img_count + 1; i++) {
            if (i < 10){
                path = "00"+i+".png";
            }
            else{
                path = "0"+i+".png"
            }

            let a = (optionPath) => {
                return async function (callback) {
                    source.path += optionPath;
                    const res = await fetch(source.path, source.options);
                    if(res.ok) {
                        const fileStream = fs.createWriteStream(__dirname + __dirname + '/chapter' + ch
                            + '/img' + i + '.png');

                        callback(null, i);
                    }
                }
            }
            filesPath.push({
                id: i,
                path: __dirname + "/chapter1/img" + i
            });//TODO надо расставить в нужных местах нормальные имена файлов


            b.push(a(path));

        }

    }
}

let ch_chooser = async function (chapters_count){
    await fs.mkdir(__dirname + '/chapters/', (e) => {
        if(e.code === 'EEXIST'){
            console.log('Directory has already exist');
        }
        else if (e) {
            throw e;
        }
    })

    let tasks_pic = [];

    for (let i = 1; i <= chapters_count; i++) {
        let source = sources[0];
        source.path += i + '/';
        tasks_pic.push(dwld_ch(source,i,ch_list[i-1]));
    }

    parallel(tasks_pic, (err, result) => {
        console.log("ch_chooser ", result);
    })
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


