const API_URL = "http://localhost:8080"; //just for testing

async function shorten() {
    const url = $("#urlInput").val();
    const encoded = encodeURIComponent(url);
    const result = await $.ajax({
        type: "POST",
        url: `${API_URL}/add?url=${encoded}`,
        dataType: "text",
    });
    console.log(result);
}
