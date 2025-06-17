document.addEventListener('DOMContentLoaded', () => {
    const freeBtn = document.getElementById('free-btn');
    const proBtn = document.getElementById('pro-btn');
    const resultsBox = document.getElementById('results');

    const makeRequest = (tenant) => {
        resultsBox.textContent = '...';
        resultsBox.className = ''; // Reset classes

        fetch(`/query?tenant=${tenant}`)
            .then(response => {
                if (!response.ok) {
                    // This will handle errors like 429 Too Many Requests
                    throw new Error(`${response.status} - ${response.statusText}`);
                }
                return response.json();
            })
            .then(data => {
                // On success, format the JSON and display it
                resultsBox.textContent = JSON.stringify(data, null, 2);
                resultsBox.classList.add('success');
            })
            .catch(error => {
                // On failure, display the error message
                resultsBox.textContent = `// ERROR: ${error.message}`;
                resultsBox.classList.add('error');
            });
    };

    freeBtn.addEventListener('click', () => makeRequest('free-tier'));
    proBtn.addEventListener('click', () => makeRequest('pro-tier'));
});