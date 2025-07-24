chrome.runtime.sendMessage({ action: "getWebsocketHost" }, function ({ host }) {
	document.getElementById('host').value = host
})

document.getElementById('websocket-form').addEventListener('submit', function (e) {
	e.preventDefault();
	const host = document.getElementById('host').value;
	chrome.runtime.sendMessage({ action: "setWebsocketHost", host: host }) 
});

