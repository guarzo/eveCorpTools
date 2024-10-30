// utils.js
export function truncateLabel(label, length) {
    if (!label || typeof label !== 'string') {
        return '';
    }
    return label.length > length ? label.substring(0, length) + '...' : label;
}

export function getColor(index) {
    const colors = [
        '#FF6384', '#36A2EB', '#FFCE56', '#4BC0C0',
        '#9966FF', '#FF9F40', '#E7E9ED', '#76D7C4',
        '#C0392B', '#8E44AD', '#2ECC71', '#1ABC9C',
        '#3498DB', '#F1C40F', '#E67E22', '#95A5A6',
    ];
    return colors[index % colors.length];
}

// Common Options (if needed globally)
export const commonOptions = {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
        legend: { display: true },
        tooltip: {
            callbacks: {
                title: function (context) {
                    const index = context[0].dataIndex;
                    return context[0].chart.data.fullLabels
                        ? context[0].chart.data.fullLabels[index]
                        : context[0].chart.data.labels[index];
                },
            },
        },
    },
    scales: {
        x: {
            ticks: { color: '#ffffff' },
            grid: { color: '#444' },
        },
        y: {
            beginAtZero: true,
            ticks: { color: '#ffffff' },
            grid: { color: '#444' },
        },
    },
};