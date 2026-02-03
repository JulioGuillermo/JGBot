export function fetchContent(url, headers = false) {
    let request = HttpRequest().SetURL(url)

    if (headers) {
        request = request
            .AddHeader("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:147.0) Gecko/20100101 Firefox/147.0")
            .AddHeader("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
            .AddHeader("Accept-Language", "es-ES,en-US;q=0.9,en;q=0.8")
            .AddHeader("Accept-Encoding", "gzip, deflate, br, zstd")
    }

    const response = request.Get();

    if (response.StatusCode != 200 && response.StatusCode != 201) {
        throw `Request return status code ${response.StatusCode} with status message ${response.Status}`;
    }

    const html = response.BodyString();

    return html;
}