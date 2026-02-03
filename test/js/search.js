export function webSearch(query) {
    // const url = `https://searx.be/search?q=${encodeURIComponent(query)}`;
    const url = `https://www.mojeek.com/search?q=${encodeURIComponent(query)}`;

    const response = HttpRequest()
        .SetURL(url)
        // User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:147.0) Gecko/20100101 Firefox/147.0
        // Accept: text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8
        // Accept-Language: es-ES,en-US;q=0.9,en;q=0.8
        // Accept-Encoding: gzip, deflate, br, zstd
        .AddHeader("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:147.0) Gecko/20100101 Firefox/147.0")
        .AddHeader("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
        .AddHeader("Accept-Language", "es-ES,en-US;q=0.9,en;q=0.8")
        .AddHeader("Accept-Encoding", "gzip, deflate, br, zstd")
        .Get();
    const html = response.BodyString();

    return html;
}