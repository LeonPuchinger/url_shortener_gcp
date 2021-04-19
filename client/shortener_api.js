const API_URL = "http://localhost:8080"; //just for testing

function hostname() {
    return window.location.hostname
}

function display_shortened(base, key, url) {
    var element = `<li class="list-group-item list-group-item-info">
        <a href="http://${base}/${key}" class="alert-link" style="display: block">${base}/${key}</a>
        <a href="${url}" class="link-secondary">${url}</a>
    </li>
    `
    $("#url-list").prepend(element)
}

function display_error(url, error) {
    var element = `<li class="list-group-item list-group-item-danger">
        <p>${error}</p>
        <a href="${url}" class="link-secondary">${url}</a>
    </li>
    `
    $("#url-list").prepend(element)
}

async function shorten() {
    const url = $("#url-input").val();
    const encoded = encodeURIComponent(url);
    try {
        const result = await $.ajax({
            type: "POST",
            url: `${API_URL}/add?url=${encoded}`,
            dataType: "json",
        });
        display_shortened(hostname(), result["key"], result["url"])
    } catch (error) {
        var msg
        if (error == undefined || error.responseJSON["error"] == undefined || error.responseJSON["error"] == "") {
            msg = "Error: unable to shorten url"
        } else {
            msg = `Error: ${error.responseJSON["error"]}`
        }
        display_error(url, msg)
    }
}
