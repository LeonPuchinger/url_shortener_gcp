const API_URL = "https://api-server-3ucl2t344q-ue.a.run.app"

//setup submit button enabled/disabled state depending on if there is input or not
function init() {
    const url = $("#url-input").on("input", function () {
        if ($("#url-input").val() == "") {
            $("#submit-button").prop("disabled", true)
        } else {
            $("#submit-button").prop("disabled", false)
        }
    })
}

function hostname() {
    // only needed if client is hosted under same domain as api/server backend
    //return window.location.hostname
    return API_URL
}

// display returned shortened and original url
function display_shortened(base, key, url) {
    var element = `<li class="list-group-item list-group-item-info">
        <a href="${base}/${key}" class="alert-link" style="display: block">${base}/${key}</a>
        <a href="${url}" class="link-secondary">${url}</a>
    </li>
    `
    $("#url-list").prepend(element)
}

// display error message and original url (if any)
function display_error(url, error) {
    var element = `<li class="list-group-item list-group-item-danger">
        <p>${error}</p>
        <a href="${url}" class="link-secondary">${url}</a>
    </li>
    `
    $("#url-list").prepend(element)
}

// show loading animation
function enable_loading() {
    $("#submit-spinner").css("display", "inline-block")
    $("#submit-text").text("Loading...")
    $("#submit-button").prop("disabled", true)
}

// hide loading animation
function disable_loading() {
    $("#submit-spinner").css("display", "none")
    $("#submit-text").text("Shorten!")
    $("#submit-button").prop("disabled", false)
}

// call api with input url, start/stop loading animation, display response/error
async function shorten() {
    enable_loading()
    const url = $("#url-input").val()
    const encoded = encodeURIComponent(url)
    try {
        const result = await $.ajax({
            type: "POST",
            url: `${API_URL}/add?url=${encoded}`,
            dataType: "json",
        })
        display_shortened(hostname(), result["key"], result["url"])
    } catch (error) {
        var msg
        if (error == undefined || error.responseJSON == undefined ||
            error.responseJSON["error"] == undefined || error.responseJSON["error"] == "") {
            msg = "Error: unable to shorten url"
        } else {
            msg = `Error: ${error.responseJSON["error"]}`
        }
        display_error(url, msg)
    } finally {
        disable_loading()
    }
}
