class SearchProvider {
    constructor(baseUrl) {
        this.baseUrl = baseUrl
    }

    search(query) { // returns a Promise with the posts list
        return fetch(`${this.baseUrl}/search?query=${query}`)
            .then(response => response.json())
            .then((data)=>data.posts)
    }
}