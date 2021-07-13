const nodemailer = require('nodemailer');
const readline = require('readline-sync');
const fs = require('fs');



module.exports = class mailer {
    #password = '';

    constructor() {
        if(fs.existsSync("./token.txt")){
            console.log("token exists");
            try{
                this.#password = fs.readFileSync('./token.txt', 'utf8');
            }
            catch(e){
                throw e;
            }

        }
        else{
            this.#password = readline.question("Enter password for Email\n");
        }
    }

    send_to_Email (filePath, fileName){
        let tranporter = nodemailer.createTransport({
            host: 'smtp.mail.ru',
            port: 465,
            secure: true,
            auth:{
                user: 'matvey2001xxl1976@mail.ru',
                pass: this.#password
            }
        })

        let result = tranporter.sendMail({
            from: 'matvey2001xxl1976@mail.ru matvey2001xxl1976@mail.ru',
            to: 'whitewolf_185@pbsync.com whitewolf_185@pbsync.com',
            attachments: {filename: fileName, path: filePath}
        })

        return result;
    }
}

async function send_to_Email (filePath, fileName){

}