window.onload = function () {
	var conn;
	var msg = document.getElementById("msg");
	var log = document.getElementById("log");

	function appendLog(item) {
		var doScroll = log.scrollTop > log.scrollHeight - log.clientHeight - 1;
		log.appendChild(item);
		if ( doScroll ) {
			log.scrollTop = log.scrollHeight - log.clientHeight;
		}
	}

	document.getElementById("form").onsubmit = function () {
		if ( !conn ) {
			return false;
		}
		if ( !msg.value ) {
			return false;
		}
		conn.send(msg.value);
		msg.value = "";
		return false;
	};

	if (window["WebSocket"]) {
		var u = new URL("../ws", document.location);
		u.protocol = "wss:";
		if ( document.location.protocol == "http:" ) {
			u.protocol = "ws:";
		}
		conn = new WebSocket( u );
		conn.onclose = function (evt) {
			var item = document.createElement("div");
			item.innerHTML = "<b>Connection closed.</b>";
			appendLog(item);
		};
		conn.onmessage = function (evt) {
			var messages = evt.data.split('\n');
			for ( var i = 0; i < messages.length; i++ ) {
				var item = document.createElement("div");
				item.innerText = messages[i];
				appendLog(item);
			}
		};
	} else {
		var item = document.createElement("div");
		item.innerHTML = "<b>Your browser does not support WebSockets.</b>";
		appendLog(item);
	}
};
