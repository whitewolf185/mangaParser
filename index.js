const util = require('util');
const fs = require('fs');
const parallel = require('async/parallel');
const exec = util.promisify(require('child_process').exec);
const fetch = require('node-fetch');
const rimraf = require('rimraf');
let mailer = require('./mailer');


//whitewolf_185@pbsync.com
// ----временная информация----
const manga_name = 'abyss';
const ch_list = [51,33,30,19,19,19,23,18,19,19,19,18,19,18,21,19];
//TODO сделать автоматичкое получение количества страниц

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
    let new_ch = ch;

    return function (callback){
        try {
            fs.mkdirSync(__dirname + '/chapter'+new_ch+'/');
        }
        catch (e) {
            if(e.code === "EEXIST"){
                console.log("Directory ok");
            }
        }



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
                return async function () {
                    const res = await fetch(optionPath, source.options);
                    console.log(res.statusText);
                    if(res.statusText === "OK") {
                        const fileStream = fs.createWriteStream( __dirname + '/chapter' + new_ch
                            + '/img' + i + '.png');
                        await new Promise((resolve, reject) => {
                            res.body.pipe(fileStream)
                            res.body.on("error", reject);
                            fileStream.on("finish", resolve);
                        })
                        console.log("117 ", i);
                    }
                    return i;
                }
            }
            b.push(a(path));

        }

        parallel(b, (err, result) => {
            console.log("I've completed tasks ", result);
            callback(null, "chapter " + new_ch);
            if(err) throw err;
        })

    }
}

let ch_chooser = function (chapters_count, downloadOptions){
    if(downloadOptions.send_Email === true){
        new mailer();
    }

    try {
        fs.mkdirSync(__dirname + '/chapters/');
    }
    catch (e) {
        if(e.code === "EEXIST"){
            console.log("Directory ok");
        }
    }


    new Promise(resolve => {
        let tasks_pic = [];

        for (let i = 1; i <= chapters_count; i++) {
            let source = {...sources[0]};
            source.path += i + '/';
            tasks_pic.push(dwld_ch(source,i,ch_list[i-1]));
        }

        parallel(tasks_pic, (err, result) => {
            console.log("ch_chooser ", result);
            if(err) throw err;
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
                            if(downloadOptions.send_Email === true) {
                                mailer.send_to_Email(__dirname + "/chapter_1.pdf", "Abyss_ch1.pdf")
                                    .then(result =>{
                                        // cleanOut();
                                    })
                                    .catch(e => {
                                        console.log(e);
                                    })
                            }
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



ch_chooser(16, {
    send_Email: true
});



