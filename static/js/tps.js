// Function to truncate labels
function truncateLabel(label, length) {
    return label.length > length ? label.substring(0, length) + '...' : label;
}

// Store full labels for tooltips
window.mtdKillCountData.fullLabels = [...window.mtdKillCountData.labels];
window.lmKillCountData.fullLabels = [...window.lmKillCountData.labels];
window.ytdKillCountData.fullLabels = [...window.ytdKillCountData.labels];

// Truncate labels for x-axis
window.mtdKillCountData.labels = window.mtdKillCountData.labels.map(label => truncateLabel(label, 10));
window.lmKillCountData.labels = window.lmKillCountData.labels.map(label => truncateLabel(label, 10));
window.ytdKillCountData.labels = window.ytdKillCountData.labels.map(label => truncateLabel(label, 10));

// Set background colors
window.mtdKillCountData.datasets[0].backgroundColor = 'rgba(255, 77, 77, 0.7)';
window.lmKillCountData.datasets[0].backgroundColor = 'rgba(255, 77, 77, 0.7)';
window.ytdKillCountData.datasets[0].backgroundColor = 'rgba(255, 77, 77, 0.7)';

// Chart options
const killCountChartOptions = {
    responsive: true,
    maintainAspectRatio: false,
    animation: {
        duration: 1000,
        easing: 'easeOutBounce',
    },
    plugins: {
        legend: { display: false },
        tooltip: {
            backgroundColor: '#23272a',
            titleColor: '#ffffff',
            bodyColor: '#ffffff',
            borderColor: '#7289da',
            borderWidth: 1,
            callbacks: {
                title: function(context) {
                    const index = context[0].dataIndex;
                    const fullLabel = context[0].chart.data.fullLabels[index];
                    return fullLabel;
                },
                label: function(context) {
                    let label = context.dataset.label || '';
                    if (label) {
                        label += ': ';
                    }
                    if (context.parsed.y !== null) {
                        label += context.parsed.y;
                    }
                    return label;
                },
            },
        },
    },
    scales: {
        x: {
            ticks: {
                color: '#ffffff',
                font: {
                    size: 10,
                },
                maxRotation: 45,
                minRotation: 45,
                autoSkip: false,
                // Removed 'callback' function
            },
            grid: { display: false },
        },
        y: {
            ticks: {
                color: '#ffffff',
                font: {
                    size: 12,
                },
            },
            grid: { color: '#444' },
            beginAtZero: true,
        },
    },
};

// Create the charts
if (document.getElementById('mtdKillCountChart')) {
    new Chart(document.getElementById('mtdKillCountChart').getContext('2d'), {
        type: 'bar',
        data: window.mtdKillCountData,
        options: killCountChartOptions,
    });
}

if (document.getElementById('lmKillCountChart')) {
    new Chart(document.getElementById('lmKillCountChart').getContext('2d'), {
        type: 'bar',
        data: window.lmKillCountData,
        options: killCountChartOptions,
    });
}

if (document.getElementById('ytdKillCountChart')) {
    new Chart(document.getElementById('ytdKillCountChart').getContext('2d'), {
        type: 'bar',
        data: window.ytdKillCountData,
        options: killCountChartOptions,
    });
}
