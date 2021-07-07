/*
var https = require('follow-redirects').https;
var fs = require('fs');

var options = {
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
};

var req = https.request(options, function (res) {
    var chunks = [];

    res.on("data", function (chunk) {
        chunks.push(chunk);
    });

    res.on("end", function (chunk) {
        var body = Buffer.concat(chunks);
        fs.writeFile(__dirname + "\\img1.jpg", body, ()=>{});
    });

    res.on("error", function (error) {
        console.error(error);
    });
});

req.end();*/
