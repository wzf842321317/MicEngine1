window.onload=function(){
    if (getCookie("langType") == "en-US") {
        $("#language").append('中文简体');
    } else if (getCookie("langType") == "zh-CN") {
        $("#language").append('English');
    }
}

function cookieSwitch() {
    if (getCookie("langType") == "en-US") {
        setCookie("langType","zh-CN")
    } else if (getCookie("langType") == "zh-CN") {
        setCookie("langType","en-US")
    }
    document.location.reload();
}

function setCookie(c_name, value) {
    var exdate = new Date();
    document.cookie = c_name + "=" + escape(value);
}

function getCookie(name) {
    var prefix = name + "=";
    var start = document.cookie.indexOf(prefix);
    if (start == -1) {
        return null;
    }
    var end = document.cookie.indexOf(";", start + prefix.length);
    if (end == -1) {
        end = document.cookie.length;
    }
    var value = document.cookie.substring(start + prefix.length, end);
    return unescape(value);
}
