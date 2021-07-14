const nodemailer = require('nodemailer');
const readline = require('readline-sync');
const fs = require('fs');



module.exports = class mailer {
    #password = '';
    #user = 'matvey2001xxl1976@mail.ru';
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
            this.#password = readline.question("Enter password for Email\n", {
                hideEchoBack: true // The typed text on screen is hidden by `*` (default).
            });
        }
    }

    send_to_Email (filePath, fileName){
        return new Promise((resolve, reject) => {
            try {
                let transporter = nodemailer.createTransport({
                    host: 'smtp.mail.ru',
                    port: 465,
                    secure: true,
                    auth: {
                        user: this.#user,
                        pass: this.#password
                    }
                })

                transporter.sendMail({
                    from: this.#user + ' ' + this.#user,
                    to: 'whitewolf_185@pbsync.com whitewolf_185@pbsync.com',
                    attachments: {filename: fileName, path: filePath}
                }).then(() => {
                    resolve();
                })

            } catch (e) {
                reject(e);
            }
        });

    }
}