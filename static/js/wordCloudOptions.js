// static/js/chartConfigs/wordCloudOptions.js

export const wordCloudOptions = {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
        legend: { display: false },
        title: {
            display: true,
            text: 'Top Ships Killed Word Cloud',
            font: {
                size: 24,
                family: 'Montserrat, sans-serif',
                weight: 'bold',
            },
            color: '#000000',
            align: 'center',
        },
        tooltip: {
            callbacks: {
                label: function (context) {
                    const shipName = context.raw.text || 'Unknown';
                    const killCount = context.raw.weight || 0;
                    return `${shipName}: ${killCount} Killmails`;
                },
            },
            mode: 'nearest',
            intersect: true,
        },
        // Ensure datalabels are not interfering
        datalabels: false, // Explicitly disable datalabels for Word Cloud
    },
    scales: {
        x: { display: false },
        y: { display: false },
    },
    layout: {
        padding: {
            left: 10,
            right: 10,
            top: 10,
            bottom: 10,
        },
    },
};
