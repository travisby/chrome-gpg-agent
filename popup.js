setInterval(function () { console.log(chrome.runtime.lastError); }, 500);

document.getElementById('sign').onsubmit = function (e) {
    e.preventDefault();

    chrome.runtime.sendNativeMessage(
        'io.linux_fu.chrome_gpg_agent',
        {method: 'sign', message: document.getElementById('signText').value},
        function (res) {
            console.log(res.message);
        }
    );
};

document.getElementById('verify').onsubmit = function (e) {
    e.preventDefault();

    chrome.runtime.sendNativeMessage(
        'io.linux_fu.chrome_gpg_agent',
        {method: 'verify', message: document.getElementById('verifyText').value},
        function (res) {
            console.log(res.message);
        }
    );
};

document.getElementById('encrypt').onsubmit = function (e) {
    e.preventDefault();

    chrome.runtime.sendNativeMessage(
        'io.linux_fu.chrome_gpg_agent',
        {
            method: 'encrypt',
            message: document.getElementById('encText').value,
            recipient: document.getElementById('encRecip').value
        },
        function (res) {
            console.log(res.message);
        }
    );

};

document.getElementById('decrypt').onsubmit = function (e) {
    e.preventDefault();

    chrome.runtime.sendNativeMessage(
        'io.linux_fu.chrome_gpg_agent',
        {method: 'decrypt', message: document.getElementById('decText').value},
        function (res) {
            console.log(res.message);
        }
    );
};
