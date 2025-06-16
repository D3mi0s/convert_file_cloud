const FILES_URL = 'http://localhost:8081/files';

async function loadFiles() {
    try {
        const response = await fetch(FILES_URL, {
            headers: {
                'Authorization': `Bearer ${localStorage.getItem('token')}`
            }
        });
        
        const files = await response.json();
        if (!response.ok) throw new Error(files.error || 'Ошибка загрузки');
        
        renderFiles(files);
    } catch (error) {
        console.error('Ошибка:', error);
        document.getElementById('fileItems').innerHTML = `
            <div class="message error">${error.message}</div>
        `;
    }
}

function renderFiles(files) {
    const container = document.getElementById('fileItems');
    container.innerHTML = files.map(file => `
        <div class="file-item">
            <div class="file-info">
                <div class="file-name">${file.OriginalName}</div>
                <div class="file-meta">
                    <span>${formatBytes(file.Size)}</span>
                    <span>${new Date(file.CreatedAt).toLocaleDateString()}</span>
                </div>
            </div>
            <div class="file-actions">
                <span class="status-${file.Status}">
                    ${getStatusText(file.Status)}
                </span>
                ${file.Status === 'completed' ? `
                    <button class="download-btn" onclick="downloadFile('${file.ConvertedName}')">
                        <i class="fas fa-download"></i> Скачать
                    </button>` : ''}
            </div>
        </div>
    `).join('');
}

function getStatusText(status) {
    return {
        pending: 'Ожидание',
        processing: 'В процессе',
        completed: 'Готово'
    }[status] || status;
}

function downloadFile(filename) {
    window.open(`http://localhost:8083/download/${filename}`);
}


loadFiles();
setInterval(loadFiles, 10000);