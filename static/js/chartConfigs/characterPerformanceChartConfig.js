// chartConfigs/characterPerformanceChartConfig.js
import { truncateLabel, commonOptions } from '../utils.js';

const characterPerformanceChartConfig = {
    id: 'characterPerformanceChart',
    instance: null,
    dataKeys: {
        mtd: 'mtdCharacterPerformanceData',
        ytd: 'ytdCharacterPerformanceData',
        lastMonth: 'lastMCharacterPerformanceData',
    },
    type: 'bar',
    options: {
        ...commonOptions,
        plugins: {
            ...commonOptions.plugins,
            legend: { display: true },
            tooltip: {
                ...commonOptions.plugins.tooltip,
                mode: 'index',
                intersect: false,
            },
        },
        scales: {
            ...commonOptions.scales,
            y: {
                ...commonOptions.scales.y,
                type: 'linear',
                position: 'left',
            },
            y1: {
                type: 'linear',
                position: 'right',
                beginAtZero: true,
                ticks: { color: '#ffffff' },
                grid: { drawOnChartArea: false },
            },
            x: {
                ...commonOptions.scales.x,
                ticks: {
                    ...commonOptions.scales.x.ticks,
                    maxRotation: 45,
                    minRotation: 45,
                    autoSkip: false,
                },
                grid: { display: false },
            },
        },
    },
    processData: function (data) {
        // ... processing logic remains the same ...
        const labels = data.map(item => item.CharacterName || 'Unknown');
        const killCountData = data.map(item => item.KillCount || 0);
        const soloKillsData = data.map(item => item.SoloKills || 0);
        const pointsData = data.map(item => item.Points || 0);

        const fullLabels = [...labels];
        const truncatedLabels = labels.map(label => truncateLabel(label, 10));

        const datasets = [
            {
                label: 'Kill Count',
                data: killCountData,
                backgroundColor: 'rgba(255, 77, 77, 0.7)',
                yAxisID: 'y',
                type: 'bar',
            },
            {
                label: 'Solo Kills',
                data: soloKillsData,
                backgroundColor: 'rgba(54, 162, 235, 0.7)',
                yAxisID: 'y',
                type: 'bar',
            },
            {
                label: 'Points',
                data: pointsData,
                borderColor: 'rgba(255, 206, 86, 1)',
                backgroundColor: 'rgba(255, 206, 86, 0.5)',
                yAxisID: 'y1',
                type: 'line',
            },
        ];

        return { labels: truncatedLabels, datasets, fullLabels };
    },
};

export default characterPerformanceChartConfig;
