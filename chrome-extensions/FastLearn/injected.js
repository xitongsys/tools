var AUTO_STUDY = true;

function getM3U8Duration(m3u8Text) {
    let total = 0;

    // 匹配所有 #EXTINF:xxx 行
    const regex = /#EXTINF:([\d\.]+)/g;
    let match;
    while ((match = regex.exec(m3u8Text)) !== null) {
        total += parseFloat(match[1]);
    }

    return total; // 单位：秒
}

async function getM3U8DurationFromURL(url) {
    console.log("==========" + url);
    var resp = await fetch(url);
    if(!resp.ok) throw new Error(resp);
    var data = await resp.text();
    return getM3U8Duration(data);
}


console.log("Injected script is running!");

const origFetch = window.fetch;
window.fetch = async function(...args) {
    //console.log("FETCH:", args);
    return origFetch.apply(this, args);
};

const origOpen = XMLHttpRequest.prototype.open;
XMLHttpRequest.prototype.open = function(method, url) {
    //console.log("XHR:", method, url);
    return origOpen.apply(this, arguments);
};

// hook send 才能拿 body
const origSend = XMLHttpRequest.prototype.send;
XMLHttpRequest.prototype.send = async function(body) {

    if (typeof body === "string" && body.includes("log_update") && AUTO_STUDY) {

        var headers = {
            "accept": "*/*",
            "accept-language": "en-US,en;q=0.9",
            "content-type": "application/x-www-form-urlencoded; charset=UTF-8",
            "sec-ch-ua": "\"Chromium\";v=\"142\", \"Google Chrome\";v=\"142\", \"Not_A Brand\";v=\"99\"",
            "sec-ch-ua-mobile": "?0",
            "sec-ch-ua-platform": "\"Linux\"",
            "sec-fetch-dest": "empty",
            "sec-fetch-mode": "cors",
            "sec-fetch-site": "same-origin",
            "x-requested-with": "XMLHttpRequest"
        };

        var total_secs = 0;
        const iframe = document.querySelector("iframe");
        if (iframe && iframe.contentDocument) {
            const video = iframe.contentDocument.querySelector("video");
            if (video) {
                console.log("src:", video.src);
                console.log("currentSrc:", video.currentSrc);
                total_secs = await getM3U8DurationFromURL(video.src);                
            }
        } 

        console.log("total_secs=" + String(total_secs));

        if(total_secs > 0) {

            const paras = new URLSearchParams(body);
            var logs = paras.get("log_list");
            var logs = JSON.parse(logs);
            var log = logs[0];
            
            console.log(log);
            console.log(log.studytime);

            for(dt=0; dt<total_secs-10; dt+=30) {
                log.studytime = String(dt)
                const log_encoded = encodeURIComponent(JSON.stringify([log]));

                var r = await fetch("index.php?m=Api&agency=3&log_update",{
                    method: "POST",
                    headers: headers,
                    body: "act=log_update&log_list=" + log_encoded,
                });

                console.log(r)
            }

            alert("已完成学习，请刷新页面");

        }

        AUTO_STUDY = false;

    }

    return origSend.apply(this, arguments);
};