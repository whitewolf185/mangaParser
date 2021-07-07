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

let dwld_ch = function(source, ch, img_count){
    return async function (callback){
        await fs.mkdir(__dirname + '/chapter'+ch+'/', (e) => {
            if(e){
                if(e.code === 'EEXIST'){
                    console.log('Directory has already exist');
                }
                else {
                    throw e;
                }
            }
        })


        let b = [];
        let path = "";


        for (let i = 1; i < img_count + 1; i++) {
            if (i < 10){
                path = source.path+"00"+i+".png";
            }
            else{
                path = source.path+"0"+i+".png";
            }

            let a = (optionPath) => {
                return async function (callback) {
                    let local_path = optionPath;
                    const res = await fetch(local_path, source.options);
                    if(res.ok) {
                        const fileStream = fs.createWriteStream(__dirname + __dirname + '/chapter' + ch
                            + '/img' + i + '.png');
                        await new Promise((resolve, reject) => {
                            res.body.pipe(fileStream)
                            res.body.on("error", reject);
                            fileStream.on("finish", resolve);
                        })

                        callback(null, i);
                    }
                }
            }


            b.push(a(path));

        }

        await parallel(b, (err, result) => {
            console.log("I've completed tasks ", result);
            if(err) throw err;
        })

        callback(null, "chapter " + ch);

    }
}

let ch_chooser = async function (chapters_count){
    await fs.mkdir(__dirname + '/chapters/', (e) => {
        if(e){
            if(e.code === 'EEXIST'){
                console.log('Directory has already exist');
            }
            else {
                throw e;
            }
        }
    })

    new Promise(resolve => {
        let tasks_pic = [];

        for (let i = 1; i <= chapters_count; i++) {
            let source = sources[0];
            source.path += i + '/';
            tasks_pic.push(dwld_ch(source,i,ch_list[i-1]));
        }

        parallel(tasks_pic, (err, result) => {
            console.log("ch_chooser ", result);

            resolve();
        })
    })
        .then(
            () => {
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
                // cleanOut();
            }
        )


}



ch_chooser(16).then(() => {
    console.log("Done");
})



