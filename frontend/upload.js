const UPLOAD_URL = 'http://localhost:8081/upload';
let selectedFile = null;

function handleFileSelect(e) {
    selectedFile = e.target.files[0];
    const statusDiv = document.getElementById('uploadStatus');
    statusDiv.textContent = `Выбран файл: ${selectedFile.name} (${formatBytes(selectedFile.size)})`;
    statusDiv.className = 'message success';
}

async function uploadFile() {
    if (!selectedFile) {
        showUploadStatus('Пожалуйста, выберите файл', 'error');
        return;
    }

    const formData = new FormData();
    formData.append('file', selectedFile);

    try {
        const response = await fetch(UPLOAD_URL, {
            method: 'POST',
            headers: {
                'Authorization': `Bearer ${localStorage.getItem('token')}`
            },
            body: formData
        });

        const data = await response.json();
        if (response.ok) {
            showUploadStatus('Файл загружен! Конвертация начата.', 'success');
            setTimeout(loadFiles, 2000);
        } else {
            showUploadStatus(data.error || 'Ошибка загрузки', 'error');
        }
    } catch (error) {
        showUploadStatus('Ошибка сети', 'error');
    }
}

function showUploadStatus(text, type) {
    const statusDiv = document.getElementById('uploadStatus');
    statusDiv.textContent = text;
    statusDiv.className = `message ${type}`;
}

function formatBytes(bytes) {
    const units = ['B', 'KB', 'MB', 'GB'];
    let size = bytes;
    let unitIndex = 0;
    
    while (size >= 1024 && unitIndex < units.length - 1) {
        size /= 1024;
        unitIndex++;
    }
    
    return `${size.toFixed(1)} ${units[unitIndex]}`;
}