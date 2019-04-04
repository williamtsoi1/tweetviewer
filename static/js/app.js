window.onload = function () {

    console.log("Protocol: " + location.protocol);
    var wsURL = "ws://" + document.location.host + "/ws"
    if (location.protocol == 'https:') {
        wsURL = "wss://" + document.location.host + "/ws"
    }
    console.log("WS URL: " + wsURL);

    var log = document.getElementById("tweets");

    function appendLog(item) {
        var doScroll = log.scrollTop > log.scrollHeight - log.clientHeight - 1;
        log.appendChild(item);
        if (doScroll) {
            log.scrollTop = log.scrollHeight - log.clientHeight;
        }

    }

    $('*').each(function () {
        if ($(this).not(':visible')) {
            $(this).remove();
        }
    });

    if (log) {

        sock = new WebSocket(wsURL);

        sock.onopen = function () {
            console.log("connected to " + wsURL);
        };

        sock.onclose = function (e) {
            console.log("connection closed (" + e.code + ")");
        };

        sock.onmessage = function (e) {
            console.log(e);
            var t = JSON.parse(e.data);
            console.log(t);
            var item = document.createElement("div");
            item.className = "item";
            item.innerHTML = "<img src='" + t.user.profile_image_url + "'/><div class='item-text'><b>" + t.user.screen_name + ":</b><br /><i>" + t.text + "</i></div>";
            appendLog(item);
        };

    } // if log

};