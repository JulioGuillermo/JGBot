// Available documents
const docs = [
    { file: 'CUSTOM_SKILL.md', title: 'Custom Skills', desc: 'Create powerful JavaScript-based skills' },
    { file: 'AVAILABLE_SKILLS.md', title: 'Available Skills', desc: 'Built-in skills and tools' },
    { file: 'SESSION.md', title: 'Session Configuration', desc: 'Per-session behavior and tools' },
    { file: 'SCHEDULED_TASKS.md', title: 'Scheduled Tasks', desc: 'Cron jobs and timers' },
    { file: 'CONF.md', title: 'Configuration', desc: 'Full configuration reference' }
];

// Get doc file from URL hash
function getDocFile() {
    const hash = window.location.hash.slice(1); // Remove #
    if (hash) {
        // Handle both just filename and full path
        return hash.includes('.md') ? hash : hash + '.md';
    }
    return 'CUSTOM_SKILL.md'; // Default
}

// Find doc info by filename
function getDocInfo(filename) {
    return docs.find(d => d.file === filename) || { 
        file: filename, 
        title: filename.replace('.md', '').replace(/_/g, ' '), 
        desc: '' 
    };
}

// Get adjacent docs for navigation
function getAdjacentDocs(currentFile) {
    const currentIndex = docs.findIndex(d => d.file === currentFile);
    return {
        prev: currentIndex > 0 ? docs[currentIndex - 1] : null,
        next: currentIndex < docs.length - 1 ? docs[currentIndex + 1] : null
    };
}

// Update page with document
function loadDocument(filename) {
    const docInfo = getDocInfo(filename);
    const adjacent = getAdjacentDocs(filename);

    // Update title
    document.getElementById('pageTitle').textContent = docInfo.title;
    document.getElementById('docTitle').textContent = docInfo.title;
    document.title = docInfo.title + ' - JGBot Documentation';

    // Update select dropdown
    document.getElementById('docSelect').value = filename;

    // Update URL hash (without reloading)
    history.replaceState(null, null, '#' + filename);

    // Show loading
    const contentEl = document.getElementById('docContent');
    contentEl.innerHTML = `
        <div class="flex items-center justify-center py-20">
            <div class="flex flex-col items-center gap-4">
                <div class="w-10 h-10 border-4 border-cyan-500 border-t-transparent rounded-full animate-spin"></div>
                <p class="text-gray-400">Loading ${docInfo.title}...</p>
            </div>
        </div>
    `;

    // Fetch and render document
    fetch('doc/' + filename)
        .then(response => {
            if (!response.ok) {
                throw new Error('Document not found: ' + filename);
            }
            return response.text();
        })
        .then(markdown => {
            // Configure marked
            marked.setOptions({
                breaks: true,
                gfm: true
            });

            // Convert markdown to HTML
            const html = marked.parse(markdown);
            contentEl.innerHTML = '<div class="markdown-body">' + html + '</div>';

            // Update navigation
            const navEl = document.getElementById('docNav');
            let navHTML = '';
            
            if (adjacent.prev) {
                navHTML += `
                    <a href="#${adjacent.prev.file}" class="flex items-center gap-2 text-gray-400 hover:text-cyan-400 transition-colors">
                        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"></path>
                        </svg>
                        <div class="text-left">
                            <span class="text-xs text-gray-500 block">Previous</span>
                            <span class="text-sm">${adjacent.prev.title}</span>
                        </div>
                    </a>
                `;
            } else {
                navHTML += '<div></div>';
            }
            
            if (adjacent.next) {
                navHTML += `
                    <a href="#${adjacent.next.file}" class="flex items-center gap-2 text-gray-400 hover:text-cyan-400 transition-colors text-right">
                        <div>
                            <span class="text-xs text-gray-500 block">Next</span>
                            <span class="text-sm">${adjacent.next.title}</span>
                        </div>
                        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"></path>
                        </svg>
                    </a>
                `;
            }
            
            navEl.innerHTML = navHTML;
        })
        .catch(error => {
            contentEl.innerHTML = `
                <div class="text-center py-20">
                    <div class="w-16 h-16 mx-auto mb-4 bg-red-500/20 rounded-full flex items-center justify-center">
                        <svg class="w-8 h-8 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"></path>
                        </svg>
                    </div>
                    <h2 class="text-xl font-semibold text-red-400 mb-2">Document Not Found</h2>
                    <p class="text-gray-400 mb-4">${error.message}</p>
                    <a href="index.html" class="inline-flex items-center gap-2 text-cyan-400 hover:underline">
                        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18"></path>
                        </svg>
                        Back to Home
                    </a>
                </div>
            `;
        });
}

// Handle dropdown change
document.getElementById('docSelect').addEventListener('change', function() {
    if (this.value) {
        loadDocument(this.value);
    }
});

// Handle hash changes
window.addEventListener('hashchange', function() {
    const file = getDocFile();
    loadDocument(file);
});

// Initial load
document.addEventListener('DOMContentLoaded', function() {
    const file = getDocFile();
    loadDocument(file);
});
