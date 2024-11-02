// static/js/chartConfigs/8_topShipsKilledWordCloudConfig.js

export function initTopShipsKilledWordCloud() {
    document.addEventListener('DOMContentLoaded', () => {
        if (typeof Chart === 'undefined') {
            console.error('Chart.js is not loaded.');
            return;
        }

        const wordCloudController = Chart.registry.getController('wordCloud');
        if (!wordCloudController) {
            console.error('Word Cloud plugin is not loaded or not registered correctly.');
            return;
        }

        const wordCloudConfigs = [
            {
                timeframe: 'mtd',
                data: window.mtdTopShipsKilledData,
                canvasId: 'topShipsKilledWordCloud_mtd',
                title: 'Top Ships Killed (MTD)',
            },
            {
                timeframe: 'ytd',
                data: window.ytdTopShipsKilledData,
                canvasId: 'topShipsKilledWordCloud_ytd',
                title: 'Top Ships Killed (YTD)',
            },
            {
                timeframe: 'lastMonth',
                data: window.lastMTopShipsKilledData,
                canvasId: 'topShipsKilledWordCloud_lastM',
                title: 'Top Ships Killed (Last Month)',
            }
        ];

        wordCloudConfigs.forEach(config => {
            const { timeframe, data, canvasId, title } = config;

            if (!Array.isArray(data) || data.length === 0) {
                console.warn(`No data available for Word Cloud (${timeframe}).`);
                return;
            }

            const canvas = document.getElementById(canvasId);
            if (!canvas) {
                console.error(`Canvas with ID "${canvasId}" not found.`);
                return;
            }

            const ctx = canvas.getContext('2d');

            const mappedData = data.map(item => {
                if (!item || !item.Name || typeof item.KillCount !== 'number') {
                    console.warn(`Invalid data point in "${timeframe}":`, item);
                    return null;
                }
                return {
                    text: item.Name,
                    weight: item.KillCount
                };
            }).filter(item => item !== null);

            if (mappedData.length === 0) {
                console.warn(`No valid data points for Word Cloud (${timeframe}).`);
                return;
            }

            const maxWords = 10;
            const limitedData = mappedData.slice(0, maxWords);

            const maxKillCount = Math.max(...limitedData.map(d => d.weight));
            const scalingFactor = maxKillCount > 0 ? 60 / maxKillCount : 10;

            const scaledData = limitedData.map(d => ({
                text: d.text,
                weight: d.weight * scalingFactor
            }));

            const colorPalette = [
                '#1f77b4', '#ff7f0e', '#2ca02c', '#d62728', '#9467bd',
                '#8c564b', '#e377c2', '#7f7f7f', '#bcbd22', '#17becf'
            ];

            const getWordColor = (word, index) => {
                return colorPalette[index % colorPalette.length];
            };

            const labels = scaledData.map(d => d.text);
            const weights = scaledData.map(d => d.weight);
            const colors = scaledData.map((d, index) => getWordColor(d.text, index));
            const rotations = scaledData.map(() => (Math.random() > 0.5 ? 0 : 90));

            console.log(`Initializing Word Cloud for "${timeframe}" with data:`);
            console.log('Labels:', labels);
            console.log('Weights:', weights);
            console.log('Colors:', colors);
            console.log('Rotations:', rotations);

            try {
                new Chart(ctx, {
                    type: 'wordCloud',
                    data: {
                        labels: labels,
                        datasets: [{
                            data: weights,
                            color: colors,
                            rotation: rotations,
                        }]
                    },
                    options: {
                        responsive: true,
                        maintainAspectRatio: false,
                        plugins: {
                            legend: { display: false },
                            title: {
                                display: true,
                                text: title,
                                font: { size: 18, family: 'Arial', weight: 'bold' },
                                color: '#ffffff',
                                align: 'center',
                                padding: { top: 10, bottom: 20 }
                            },
                            tooltip: {
                                callbacks: {
                                    label: function (context) {
                                        const word = context.label || 'Unknown';
                                        const count = (context.raw / scalingFactor).toFixed(0) || 0;
                                        return `${word}: ${count} Killmails`;
                                    },
                                },
                                mode: 'nearest',
                                intersect: true,
                                backgroundColor: '#000000',
                                titleColor: '#ffffff',
                                bodyColor: '#ffffff',
                            },
                        },
                        scales: {
                            x: { display: false },
                            y: { display: false },
                        },
                        layout: { padding: 10 },
                        animation: { duration: 0 },
                        font: { weight: 'normal', family: 'Arial' },
                    },
                });
                console.log(`Word Cloud for "${timeframe}" successfully initialized.`);
            } catch (error) {
                console.error(`Error initializing Word Cloud (${timeframe}):`, error);
            }
        });
    })
}
